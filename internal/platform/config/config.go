package config

import (
	"fmt"
	"os"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
	App      AppConfig
	Postgres PostgresConfig
	Kafka    KafkaConfig
}

func NewConfig() (*Config, error) {
	if err := loadDotenv(); err != nil {
		return nil, err
	}

	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	return cfg, nil
}

func loadDotenv() error {
	if _, err := os.Stat(".env"); err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("stat .env: %w", err)
	}

	if err := godotenv.Load(".env"); err != nil {
		return fmt.Errorf("load .env: %w", err)
	}

	return nil
}

type AppConfig struct {
	Env      string `env:"APP_ENV" envDefault:"local"`
	LogLevel string `env:"APP_LOG_LEVEL" envDefault:"info"`
	Port     string `env:"APP_HTTP_PORT,required"`
	GRPCPort string `env:"APP_GRPC_PORT" envDefault:"50051"`
}

type PostgresConfig struct {
	Host     string `env:"PG_HOST,required"`
	Port     string `env:"PG_PORT,required"`
	DB       string `env:"PG_DB,required"`
	User     string `env:"PG_USER,required"`
	Password string `env:"PG_PASSWORD,required"`
	SSLMode  string `env:"PG_SSLMODE,required"`
}

func (pc *PostgresConfig) DSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		pc.Host, pc.Port, pc.User, pc.Password, pc.DB, pc.SSLMode,
	)
}

type KafkaConfig struct {
	Brokers                string `env:"KAFKA_BROKERS,required"`
	IncidentProcessorGroup string `env:"KAFKA_INCIDENT_PROCESSOR_GROUP" envDefault:"pulseops-incident-processor"`
}

func (kc KafkaConfig) BrokerList() []string {
	parts := strings.Split(kc.Brokers, ",")
	brokers := make([]string, 0, len(parts))
	for _, part := range parts {
		broker := strings.TrimSpace(part)
		if broker != "" {
			brokers = append(brokers, broker)
		}
	}

	return brokers
}
