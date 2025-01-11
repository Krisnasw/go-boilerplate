package config

import (
	"strings"

	"github.com/spf13/viper"
)

// GetConfig loads the configuration based on the specified environment and additional config files.
func GetConfig(env string, confFiles map[string]string) (*viper.Viper, error) {
	conf := viper.New()

	// Default to "development" if env is not specified
	if env == "" {
		env = "development"
	}

	conf.SetDefault("environment", env)

	// Application Defaults
	conf.SetDefault("app.name", "Boilerplate")
	conf.SetDefault("app.env", "local")
	conf.SetDefault("app.debug", true)
	conf.SetDefault("app.timezone", "UTC")
	conf.SetDefault("app.url", "http://localhost")
	conf.SetDefault("app.locale", "en")
	conf.SetDefault("app.fallback_locale", "en")
	conf.SetDefault("app.faker_locale", "en_US")
	conf.SetDefault("app.maintenance_driver", "file")
	conf.SetDefault("app.bcrypt_rounds", 12)
	conf.SetDefault("app.port", 8080)
	conf.SetDefault("app.grpcPort", 5050)

	// Logging Defaults
	conf.SetDefault("log.channel", "stack")
	conf.SetDefault("log.stack", "single")
	conf.SetDefault("log.deprecations_channel", "null")
	conf.SetDefault("log.level", "debug")

	// Database Defaults
	conf.SetDefault("database.connection", "mysql")
	conf.SetDefault("database.host", "127.0.0.1")
	conf.SetDefault("database.port", 3306)
	conf.SetDefault("database.name", "example")
	conf.SetDefault("database.username", "root")
	conf.SetDefault("database.password", "")
	conf.SetDefault("database.logEnabled", true)
	conf.SetDefault("database.logLevel", 3)
	conf.SetDefault("database.logThreshold", 200)

	// Session Defaults
	conf.SetDefault("session.driver", "database")
	conf.SetDefault("session.lifetime", 120)
	conf.SetDefault("session.encrypt", false)
	conf.SetDefault("session.path", "/")
	conf.SetDefault("session.domain", "null")

	// Caching Defaults
	conf.SetDefault("cache.store", "database")
	conf.SetDefault("cache.prefix", "")

	// Redis Defaults
	conf.SetDefault("redis.client", "predis")
	conf.SetDefault("redis.host", "127.0.0.1")
	conf.SetDefault("redis.password", "null")
	conf.SetDefault("redis.port", 6379)

	// Queue Defaults
	conf.SetDefault("queue.connection", "sync")

	// Mail Defaults
	conf.SetDefault("mail.driver", "smtp")
	conf.SetDefault("mail.host", "smtp.example.com")
	conf.SetDefault("mail.port", 587)
	conf.SetDefault("mail.encryption", "tls")
	conf.SetDefault("mail.username", "your-email@example.com")
	conf.SetDefault("mail.password", "your-password")
	conf.SetDefault("mail.from_address", "noreply@example.com")
	conf.SetDefault("mail.from_name", "")

	// WooCommerce Defaults
	conf.SetDefault("woocommerce.url", "http://example.com")
	conf.SetDefault("woocommerce.store_url", "http://example.com/wp-json")
	conf.SetDefault("woocommerce.consumer_key", "")
	conf.SetDefault("woocommerce.consumer_secret", "")
	conf.SetDefault("woocommerce.wcfmp_username", "")
	conf.SetDefault("woocommerce.wcfmp_password", "")

	// AWS Defaults
	conf.SetDefault("aws.access_key_id", "")
	conf.SetDefault("aws.secret_access_key", "")
	conf.SetDefault("aws.region", "ap-southeast-1")
	conf.SetDefault("aws.bucket", "")
	conf.SetDefault("aws.photo_url", "")

	// Google Drive Defaults
	conf.SetDefault("google_drive.client_id", "")
	conf.SetDefault("google_drive.client_secret", "")
	conf.SetDefault("google_drive.refresh_token", "")
	conf.SetDefault("google_drive.folder", "")

	// OneBrick Defaults
	conf.SetDefault("onebrick.client_id", "")
	conf.SetDefault("onebrick.client_secret", "")
	conf.SetDefault("onebrick.base_url", "")
	conf.SetDefault("onebrick.platform_fee", 0)
	conf.SetDefault("onebrick.vat_amount", 0)

	// Backoffice Defaults
	conf.SetDefault("backoffice_url", "http://localhost:3000")

	// Firebase Defaults
	conf.SetDefault("firebase.project_id", "")

	// RabbitMQ Defaults
	conf.SetDefault("rabbitmq.host", "")
	conf.SetDefault("rabbitmq.port", 5672)
	conf.SetDefault("rabbitmq.user", "")
	conf.SetDefault("rabbitmq.password", "")
	conf.SetDefault("rabbitmq.vhost", "/")

	// OFD Sync Defaults
	conf.SetDefault("ofd_sync_delay", 10)

	// Scout Defaults
	conf.SetDefault("scout.driver", "")
	conf.SetDefault("scout.prefix", "")
	conf.SetDefault("scout.queue", "")
	conf.SetDefault("scout.identify", "")

	conf.SetDefault("health_check.route.group", "health")
	conf.SetDefault("health_check.route.live", "/live")
	conf.SetDefault("health_check.route.ready", "/ready")

	// Machinery Defaults
	conf.SetDefault("machinery.broker_dsn", "")
	conf.SetDefault("machinery.broker.retries", 2)
	conf.SetDefault("machinery.broker.retry_delay", 10)
	conf.SetDefault("machinery.broker.timeout", 100)
	conf.SetDefault("machinery.broker.max_conn", 10)
	conf.SetDefault("machinery.broker.vol_threshold", 20)
	conf.SetDefault("machinery.broker.sleep_window", 5000)
	conf.SetDefault("machinery.broker.err_per_threshold", 50)
	conf.SetDefault("machinery.default_queue", "")
	conf.SetDefault("machinery.result_backend_dsn", "")
	conf.SetDefault("machinery.exchange", "")
	conf.SetDefault("machinery.exchange_type", "")
	conf.SetDefault("machinery.binding_key", "")
	conf.SetDefault("machinery.consumer.enable", 0)
	conf.SetDefault("machinery.consumer.tag", "")
	conf.SetDefault("machinery.consumer.concurrent_tasks", 10)
	conf.SetDefault("machinery.consumer.prefetch_count", 1)

	// Environment Variable Handling
	conf.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	conf.AutomaticEnv()

	// Load YAML Config Files
	conf.SetConfigType("yaml")
	conf.SetConfigName(env)
	conf.AddConfigPath("./config/" + env)
	err := conf.ReadInConfig() // Find and read the config file
	if err != nil {
		return nil, err
	}

	// Merge Additional Config Files
	for confFile := range confFiles {
		conf.SetConfigName(confFile)
		if err := conf.MergeInConfig(); err != nil {
			return nil, err
		}
	}

	return conf, nil
}
