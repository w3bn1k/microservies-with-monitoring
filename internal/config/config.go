package config

import (
	"log"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Service    ServiceConfig    `mapstructure:"service"`
	Kafka      KafkaConfig      `mapstructure:"kafka"`
	Redis      RedisConfig      `mapstructure:"redis"`
	Postgres   PostgresConfig   `mapstructure:"postgres"`
	Monitoring MonitoringConfig `mapstructure:"monitoring"`
}

type ServiceConfig struct {
	Name     string `mapstructure:"name"`
	Port     int    `mapstructure:"port"`
	GRPCPort int    `mapstructure:"grpc_port"`
}

type KafkaConfig struct {
	Brokers []string `mapstructure:"brokers"`
	Topic   string   `mapstructure:"topic"`
	GroupID string   `mapstructure:"group_id"`
}

type RedisConfig struct {
	Addr     string        `mapstructure:"addr"`
	Password string        `mapstructure:"password"`
	DB       int           `mapstructure:"db"`
	Timeout  time.Duration `mapstructure:"timeout"`
}

type PostgresConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Database string `mapstructure:"database"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	SSLMode  string `mapstructure:"ssl_mode"`
}

type MonitoringConfig struct {
	PrometheusPort int    `mapstructure:"prometheus_port"`
	JaegerEndpoint string `mapstructure:"jaeger_endpoint"`
}

func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs")
	viper.AddConfigPath(".")

	viper.SetDefault("service.port", 8080)
	viper.SetDefault("service.grpc_port", 9090)
	viper.SetDefault("kafka.brokers", []string{"localhost:9092"})
	viper.SetDefault("kafka.topic", "user-events")
	viper.SetDefault("kafka.group_id", "consumer-group")
	viper.SetDefault("redis.addr", "localhost:6379")
	viper.SetDefault("redis.db", 0)
	viper.SetDefault("redis.timeout", "5s")
	viper.SetDefault("postgres.host", "localhost")
	viper.SetDefault("postgres.port", 5432)
	viper.SetDefault("postgres.ssl_mode", "disable")
	viper.SetDefault("monitoring.prometheus_port", 9090)
	viper.SetDefault("monitoring.jaeger_endpoint", "http://localhost:14268/api/traces")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Warning: Could not read config file: %v", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

func LoadFromFile(filename string) (*Config, error) {
	viper.SetConfigFile(filename)

	viper.SetDefault("service.port", 8080)
	viper.SetDefault("service.grpc_port", 7070)
	viper.SetDefault("kafka.brokers", []string{"localhost:9092"})
	viper.SetDefault("kafka.topic", "user-events")
	viper.SetDefault("kafka.group_id", "consumer-group")
	viper.SetDefault("redis.addr", "localhost:6379")
	viper.SetDefault("redis.db", 0)
	viper.SetDefault("redis.timeout", "5s")
	viper.SetDefault("postgres.host", "localhost")
	viper.SetDefault("postgres.port", 5432)
	viper.SetDefault("postgres.ssl_mode", "disable")
	viper.SetDefault("monitoring.prometheus_port", 9090)
	viper.SetDefault("monitoring.jaeger_endpoint", "http://localhost:14268/api/traces")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
