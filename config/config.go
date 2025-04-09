package config

import (
	"fmt"
	"os"
)

const (
	//!
	HTTP_HOST     = "HTTP_HOST"
	HTTP_PORT     = "HTTP_PORT"
	HTTP_BASE_URL = "BASE_URL"
	//!
	LOG_FILE = "LOG_FILE"
	//!
	KAFKA_BROKER = "KAFKA_BROKER"
	//!
	MYSQL_HOST     = "MYSQL_HOST"
	MYSQL_PORT     = "MYSQL_PORT"
	MYSQL_USER     = "MYSQL_USER"
	MYSQL_PASSWORD = "MYSQL_PASSWORD"
	MYSQL_DB       = "MYSQL_DB"
)

type Config struct {
	Http   HttpConfig
	Logger LoggerConfig
	Kafka  KafkaConfig
	MySQL  MySQLConfig
}

type MySQLConfig struct {
	MysqlHost     string
	MysqlPort     string
	MysqlUser     string
	MysqlPassword string
	MysqlDBName   string
}

type HttpConfig struct {
	Host string
	Port string
}

type LoggerConfig struct {
	LogFile string
}

type KafkaConfig struct {
	Brokers []string
}

func NewConfig() *Config {
	return &Config{
		Http: HttpConfig{
			Host: getEnv(HTTP_HOST, "localhost"),
			Port: getEnv(HTTP_PORT, ":8080"),
		},
		Logger: LoggerConfig{
			LogFile: getEnv(LOG_FILE, "admetric.log"),
		},
		Kafka: KafkaConfig{
			Brokers: []string{getEnv(KAFKA_BROKER, "localhost:9092")},
		},
		MySQL: MySQLConfig{
			MysqlHost:     getEnv(MYSQL_HOST, "localhost"),
			MysqlPort:     getEnv(MYSQL_PORT, "3306"),
			MysqlUser:     getEnv(MYSQL_USER, "root"),
			MysqlPassword: getEnv(MYSQL_PASSWORD, ""),
			MysqlDBName:   getEnv(MYSQL_DB, "admetric"),
		},
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func (c *Config) Parse() {
	parseError := map[string]string{
		//!
		HTTP_HOST:     "",
		HTTP_PORT:     "",
		HTTP_BASE_URL: "",
		//!
		LOG_FILE: "",
		//!
		KAFKA_BROKER: "",
		//!
		MYSQL_HOST:     "",
		MYSQL_PORT:     "",
		MYSQL_USER:     "",
		MYSQL_PASSWORD: "",
		MYSQL_DB:       "",
	}
	//! http configs
	parseError[HTTP_HOST] = c.Http.Host
	parseError[HTTP_PORT] = c.Http.Port
	parseError[HTTP_BASE_URL] = c.Http.Host
	//! Logger configs
	parseError[LOG_FILE] = c.Logger.LogFile
	//! kafka
	parseError[KAFKA_BROKER]=c.Kafka.Brokers[0]
	//! mysql configs
	parseError[MYSQL_HOST] = c.MySQL.MysqlHost
	parseError[MYSQL_PORT] = c.MySQL.MysqlPort
	parseError[MYSQL_USER] = c.MySQL.MysqlUser
	parseError[MYSQL_PASSWORD] = c.MySQL.MysqlPassword
	parseError[MYSQL_DB] = c.MySQL.MysqlDBName
	//! check all env vars are set
	exitParse := false
	for k, v := range parseError {
		if v == "" {
			exitParse = true
			fmt.Printf("%s = %s\n", k, v)
		}
	}
	if exitParse {
		panic("Env vars not set see list")
	}
}
