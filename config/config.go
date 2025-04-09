package config

import (
	"fmt"
	"os"
)

const (
	//!
	HTTP_HOST     = "HTTP_HOST"
	HTTP_PORT     = "HTTP_PORT"
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
	c := Config{
		Http: HttpConfig{
			Host: getEnv(HTTP_HOST),
			Port: getEnv(HTTP_PORT),
		},
		Logger: LoggerConfig{
			LogFile: getEnv(LOG_FILE),
		},
		Kafka: KafkaConfig{
			Brokers: []string{getEnv(KAFKA_BROKER)},
		},
		MySQL: MySQLConfig{
			MysqlHost:     getEnv(MYSQL_HOST),
			MysqlPort:     getEnv(MYSQL_PORT),
			MysqlUser:     getEnv(MYSQL_USER),
			MysqlPassword: getEnv(MYSQL_PASSWORD),
			MysqlDBName:   getEnv(MYSQL_DB),
		},
	}
	fmt.Println(c)
	return &c
}

func getEnv(key string) string {
	value := os.Getenv(key)
	return value
}

func (c *Config) Parse() {
	parseError := map[string]string{
		//!
		HTTP_HOST:     "",
		HTTP_PORT:     "",
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
	//! Logger configs
	parseError[LOG_FILE] = c.Logger.LogFile
	//! kafka
	parseError[KAFKA_BROKER] = c.Kafka.Brokers[0]
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
