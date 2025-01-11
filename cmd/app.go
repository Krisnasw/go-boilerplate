package main

import (
	"context"
	"github.com/cloudwego/hertz/pkg/common/adaptor"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/hertz-contrib/swagger"
	"go-grst-boilerplate/routers/gateways"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/heptiolabs/healthcheck"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	swaggerFiles "github.com/swaggo/files"
	_ "go-grst-boilerplate/cmd/docs"
	"go-grst-boilerplate/config"
	"go-grst-boilerplate/database/mysql"
	"go-grst-boilerplate/routers"
)

func main() {
	// Load configuration
	conf, err := config.GetConfig(os.Getenv("environment"), nil)
	if err != nil {
		zap.L().Fatal("Failed to load configuration: %v", zap.Error(err))
	}

	// Initialize logger
	zapLogger := initLogger(conf)
	defer zapLogger.Sync()

	// Initialize database
	db := initDatabase(conf, zapLogger)

	// Initialize health checks
	healthChecker := initHealthChecks(conf, db)

	// Start health check server
	go startHealthCheckServer(healthChecker, conf, zapLogger)

	// Start gRPC server
	grpcServer := initGRPCServer(conf, db)
	go startGRPCServer(grpcServer, conf, zapLogger)

	// Start gRPC-Gateway server with Hertz
	gatewayServer := initGatewayServer(conf)
	startGatewayServer(gatewayServer, zapLogger)
}

// initLogger initializes and returns a zap logger based on the configuration.
func initLogger(conf *viper.Viper) *zap.Logger {
	zapConfig := zap.NewProductionConfig()
	zapConfig.Level.UnmarshalText([]byte(conf.GetString("log.level")))
	zapConfig.Development = conf.GetString("environment") != "production"

	zapLogger, err := zapConfig.Build()
	if err != nil {
		zap.L().Fatal("Failed to initialize logger: %v", zap.Error(err))
	}
	return zapLogger
}

// initDatabase initializes and returns a database connection based on the configuration.
func initDatabase(conf *viper.Viper, zapLogger *zap.Logger) *gorm.DB {
	db, err := mysql.Connect(
		conf.GetString("database.host"),
		conf.GetInt("database.port"),
		conf.GetString("database.username"),
		conf.GetString("database.password"),
		conf.GetString("database.name"),
		mysql.SetPrintLog(
			conf.GetBool("database.logEnabled"),
			logger.LogLevel(conf.GetInt("database.logLevel")),
			time.Duration(conf.GetInt("database.logThreshold"))*time.Millisecond,
		),
	)
	if err != nil {
		zapLogger.Fatal("Failed to initialize database connection", zap.Error(err))
	}

	return db
}

// initHealthChecks initializes and returns a health check handler with database and goroutine checks.
func initHealthChecks(conf *viper.Viper, db *gorm.DB) healthcheck.Handler {
	healthChecker := healthcheck.NewMetricsHandler(prometheus.DefaultRegisterer, "health_check")
	healthChecker.AddLivenessCheck("Goroutine Threshold", healthcheck.GoroutineCountCheck(conf.GetInt("health_check.goroutine_threshold")))

	sqlDB, _ := db.DB()
	healthChecker.AddReadinessCheck(conf.GetString("database.driver"), healthcheck.DatabasePingCheck(sqlDB, 1*time.Second))

	return healthChecker
}

// startHealthCheckServer starts an HTTP server to expose health check endpoints.
func startHealthCheckServer(healthChecker healthcheck.Handler, conf *viper.Viper, zapLogger *zap.Logger) {
	healthCheckPort := conf.GetString("health_check.port")
	if healthCheckPort == "" {
		healthCheckPort = "8081" // Default port for health checks
	}

	http.HandleFunc(conf.GetString("health_check.route.group")+conf.GetString("health_check.route.live"), healthChecker.LiveEndpoint)
	http.HandleFunc(conf.GetString("health_check.route.group")+conf.GetString("health_check.route.ready"), healthChecker.ReadyEndpoint)

	zapLogger.Info("Health check server starting on Port : " + healthCheckPort)
	if err := http.ListenAndServe(":"+healthCheckPort, nil); err != nil {
		zapLogger.Fatal("Failed to start health check server", zap.Error(err))
	}
}

// initGRPCServer initializes and returns a gRPC server with registered microservices.
func initGRPCServer(conf *viper.Viper, db *gorm.DB) *grpc.Server {
	srv := grpc.NewServer()
	reflection.Register(srv)
	routers.InitMicroservices(db, srv)
	return srv
}

// startGRPCServer starts the gRPC server on the configured port.
func startGRPCServer(srv *grpc.Server, conf *viper.Viper, zapLogger *zap.Logger) {
	listenGrpcPort, err := net.Listen("tcp", ":"+conf.GetString("app.grpcPort"))
	if err != nil {
		zapLogger.Fatal("Failed to listen gRPC Port", zap.Error(err))
	}

	zapLogger.Info("gRPC Starting on Port : " + conf.GetString("app.grpcPort"))
	if err := srv.Serve(listenGrpcPort); err != nil {
		zapLogger.Fatal("Failed to serve gRPC server", zap.Error(err))
	}
}

// initGatewayServer initializes and returns a Hertz server for the gRPC-Gateway.
func initGatewayServer(conf *viper.Viper) *server.Hertz {
	conn, err := grpc.NewClient(
		"localhost:"+conf.GetString("app.grpcPort"),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		zap.L().Fatal("Failed to dial server: %v", zap.Error(err))
	}

	gatewayMux := runtime.NewServeMux()

	err = gateways.RegisterGatewayHandlers(context.Background(), gatewayMux, conn)
	if err != nil {
		return nil
	}

	h := server.New(server.WithHostPorts(":" + conf.GetString("app.port")))
	h.GET("/", func(c context.Context, ctx *app.RequestContext) {
		ctx.JSON(consts.StatusOK, utils.H{
			"error":   false,
			"message": "Welcome to the Go gRPC and Rest Microservice",
		})
	})

	// Swagger
	url := swagger.URL("http://localhost:8080/api/docs.json") // The url pointing to API definition
	h.GET("/api/*any", swagger.WrapHandler(swaggerFiles.Handler, url))

	h.Any("/*path", func(c context.Context, ctx *app.RequestContext) {
		// Adapt Hertz's RequestContext to http.ResponseWriter and http.Request
		httpReq, err := adaptor.GetCompatRequest(&ctx.Request)
		if err != nil {
			ctx.JSON(consts.StatusInternalServerError, utils.H{"error": "failed to adapt request"})
			return
		}

		httpWriter := adaptor.GetCompatResponseWriter(&ctx.Response)
		gatewayMux.ServeHTTP(httpWriter, httpReq)
	})

	return h
}

// startGatewayServer starts the Hertz server for the gRPC-Gateway.
func startGatewayServer(gwServer *server.Hertz, zapLogger *zap.Logger) {
	zapLogger.Info("Serving gRPC-Gateway for REST on Port :" + gwServer.GetOptions().Addr)
	if err := gwServer.Run(); err != nil {
		zapLogger.Fatal("Failed to serve gRPC-Gateway server", zap.Error(err))
	}
}
