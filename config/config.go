package config

import (
	"fmt"
	"os"
	"strconv"
)

const (
	//!
	HTTP_HOST     = "HTTP_HOST"
	HTTP_PORT     = "HTTP_PORT"
	HTTP_BASE_URL = "BASE_URL"
	//!
	LOGS_FILE = "LOG_FILE"
	//!
	REDIS_URL = "REDIS_URL"
	//!
	MYSQL_HOST     = "MYSQL_HOST"
	MYSQL_PORT     = "MYSQL_PORT"
	MYSQL_USER     = "MYSQL_USER"
	MYSQL_PASSWORD = "MYSQL_PASSWORD"
	MYSQL_DB       = "MYSQL_DB"
)

type Config struct {
	MySQL  MySQLConfig
	Redis  RedisConfig
	Http   HttpConfig
	Logger LoggerConfig
	Kafka  KafkaConfig
}

type MySQLConfig struct {
	MysqlHost     string
	MysqlPort     string
	MysqlUser     string
	MysqlPassword string
	MysqlDBName   string
}

type RedisConfig struct {
	Url string
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
		MySQL: MySQLConfig{
			MysqlHost:     getEnv("MYSQL_HOST", "localhost"),
			MysqlPort:     getEnv("MYSQL_PORT", "3306"),
			MysqlUser:     getEnv("MYSQL_USER", "root"),
			MysqlPassword: getEnv("MYSQL_PASSWORD", ""),
			MysqlDBName:   getEnv("MYSQL_DB", "admetric"),
		},
		Redis: RedisConfig{
			Url: getEnv("REDIS_URL", "localhost:6379"),
		},
		Http: HttpConfig{
			Host: getEnv("HTTP_HOST", "localhost"),
			Port: getEnv("HTTP_PORT", ":8080"),
		},
		Logger: LoggerConfig{
			LogFile: getEnv("LOG_FILE", "admetric.log"),
		},
		Kafka: KafkaConfig{
			Brokers: []string{getEnv("KAFKA_BROKER", "localhost:9092")},
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
func getEnvInt(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	intValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return intValue
}

func (c *Config) Parse() {
	parseError := map[string]string{
		//!
		HTTP_HOST:     "",
		HTTP_PORT:     "",
		HTTP_BASE_URL: "",
		//!
		LOGS_FILE: "",
		//!
		REDIS_URL: "",
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
	parseError[LOGS_FILE] = c.Logger.LogFile
	//! Redi configs
	parseError[REDIS_URL] = c.Redis.Url
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
