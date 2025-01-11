package mysql

import (
	"fmt"
	"log"
	"os"
	"time"

	driverMysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

type mysql struct {
	DBHost         string
	DBPort         int
	DBUserName     string
	DBPassword     string
	DBDatabaseName string
	DBTimezone     string

	printLog     bool
	logLevel     logger.LogLevel
	logThreshold time.Duration

	maxIdleConnection             int
	maxOpenConnection             int
	connectionMaxLifetimeInSecond int
	namingStrategy                schema.Namer
}

type mysqlOption func(*mysql)

func Connect(DBHost string, DBPort int, DBUserName string, DBPassword string, DBDatabaseName string, options ...mysqlOption) (*gorm.DB, error) {
	db := &mysql{
		DBHost:         DBHost,
		DBPort:         DBPort,
		DBUserName:     DBUserName,
		DBPassword:     DBPassword,
		DBDatabaseName: DBDatabaseName,
		DBTimezone:     "UTC",

		printLog:     false,
		logLevel:     logger.Silent,
		logThreshold: 200 * time.Millisecond,

		maxIdleConnection:             5,
		maxOpenConnection:             10,
		connectionMaxLifetimeInSecond: 60,
		namingStrategy:                nil,
	}

	for _, o := range options {
		o(db)
	}

	return connect(db)
}

func connect(param *mysql) (*gorm.DB, error) {
	// Construct MySQL DSN
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=%s",
		param.DBUserName,
		param.DBPassword,
		param.DBHost,
		param.DBPort,
		param.DBDatabaseName,
		param.DBTimezone,
	)

	// GORM Config
	cfg := &gorm.Config{}
	if param.printLog {
		cfg.Logger = logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), logger.Config{
			SlowThreshold: param.logThreshold,
			LogLevel:      param.logLevel,
			Colorful:      true,
		})
	}
	if param.namingStrategy != nil {
		cfg.NamingStrategy = param.namingStrategy
	}

	// Open Database Connection
	db, err := gorm.Open(driverMysql.Open(dsn), cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Set Connection Pool Settings
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database handle: %w", err)
	}
	sqlDB.SetMaxOpenConns(param.maxOpenConnection)
	sqlDB.SetConnMaxLifetime(time.Duration(param.connectionMaxLifetimeInSecond) * time.Second)
	sqlDB.SetMaxIdleConns(param.maxIdleConnection)

	return db, nil
}

func SetMaxIdleConns(conns int) mysqlOption {
	return func(c *mysql) {
		if conns > 0 {
			c.maxIdleConnection = conns
		}
	}
}

func SetMaxOpenConns(conns int) mysqlOption {
	return func(c *mysql) {
		if conns > 0 {
			c.maxOpenConnection = conns
		}
	}
}

func SetConnMaxLifetime(conns int) mysqlOption {
	return func(c *mysql) {
		if conns > 0 {
			c.connectionMaxLifetimeInSecond = conns
		}
	}
}

func SetNamingStrategy(namingStrategy schema.Namer) mysqlOption {
	return func(c *mysql) {
		c.namingStrategy = namingStrategy
	}
}

// SetPrintLog level: 1=silent, 2=Error, 3=Warn, 4=Info. latencyThreshold: suggestion 200ms.
func SetPrintLog(isEnable bool, level logger.LogLevel, latencyThreshold time.Duration) mysqlOption {
	return func(c *mysql) {
		if latencyThreshold > 0 {
			c.printLog = isEnable
			c.logLevel = level
			c.logThreshold = latencyThreshold
		}
	}
}

func SetTimezone(timezone string) mysqlOption {
	return func(c *mysql) {
		if timezone == "" {
			c.DBTimezone = "Etc/UTC"
		} else {
			c.DBTimezone = timezone
		}
	}
}

// SetTablePrefix sets the table prefix for all tables
func SetTablePrefix(prefix string) mysqlOption {
	return func(c *mysql) {
		c.namingStrategy = schema.NamingStrategy{
			TablePrefix: prefix,
		}
	}
}
