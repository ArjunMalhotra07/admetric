package config

import (
	"fmt"
	"os"
)

const (
	//!
	HTTP_HOST     = "HTTP_HOST"
	HTTP_PORT     = "HTTP_PORT"
	HTTP_BASE_URL = "HTTP_BASE_URL"
	//!
	LOGS_FILE = "LOG_FILE"
	//!
	MYSQL_HOST     = "MYSQL_HOST"
	MYSQL_PORT     = "MYSQL_PORT"
	MYSQL_USER     = "MYSQL_USER"
	MYSQL_PASSWORD = "MYSQL_PASSWORD"
	MYSQL_DB       = "MYSQL_DB"
)

type Config struct {
	Http   Http
	MySQL  MySQL
	Logger Logger
}

type Http struct {
	BaseUrl string
	Host    string
	Port    string
}

type MySQL struct {
	MysqlHost     string
	MysqlPort     string
	MysqlUser     string
	MysqlPassword string
	MysqlDBName   string
}

type Logger struct {
	DisableCaller     bool
	DisableStacktrace bool
	Encoding          string
	Level             string
	LogFile           string
}

func NewConfig() *Config {
	http := Http{}
	logger := Logger{}
	mysql := MySQL{}
	c := &Config{
		Http:   http,
		Logger: logger,
		MySQL:  mysql,
	}
	parseError := map[string]string{
		//!
		HTTP_HOST:     "",
		HTTP_PORT:     "",
		HTTP_BASE_URL: "",
		//!
		LOGS_FILE: "",
		//!
		MYSQL_HOST:     "",
		MYSQL_PORT:     "",
		MYSQL_USER:     "",
		MYSQL_PASSWORD: "",
		MYSQL_DB:       "",
	}
	//! http configs
	httpHost := os.Getenv(HTTP_HOST)
	if httpHost != "" {
		c.Http.Host = httpHost
		parseError[HTTP_HOST] = httpHost
	}
	httpPort := os.Getenv(HTTP_PORT)
	if httpPort != "" {
		c.Http.Port = httpPort
		parseError[HTTP_PORT] = httpPort
	}
	httpBaseURL := os.Getenv(HTTP_BASE_URL)
	if httpBaseURL != "" {
		c.Http.BaseUrl = httpBaseURL
		parseError[HTTP_BASE_URL] = httpBaseURL
	}
	//! Logger configs
	logsFile := os.Getenv(LOGS_FILE)
	if logsFile != "" {
		c.Logger.LogFile = logsFile
		parseError[LOGS_FILE] = logsFile
	}
	//! mysql configs
	mysqlHost := os.Getenv(MYSQL_HOST)
	if mysqlHost != "" {
		c.MySQL.MysqlHost = mysqlHost
		parseError[MYSQL_HOST] = mysqlHost
	}
	mysqlPort := os.Getenv(MYSQL_PORT)
	if mysqlPort != "" {
		c.MySQL.MysqlPort = mysqlPort
		parseError[MYSQL_PORT] = mysqlPort

	}
	mysqlUser := os.Getenv(MYSQL_USER)
	if mysqlUser != "" {
		c.MySQL.MysqlUser = mysqlUser
		parseError[MYSQL_USER] = mysqlUser
	}
	mysqlPassword := os.Getenv(MYSQL_PASSWORD)
	if mysqlPassword != "" {
		c.MySQL.MysqlPassword = mysqlPassword
		parseError[MYSQL_PASSWORD] = mysqlPassword
	}
	mysqlDBName := os.Getenv(MYSQL_DB)
	if mysqlDBName != "" {
		c.MySQL.MysqlDBName = mysqlDBName
		parseError[MYSQL_DB] = mysqlDBName
	}
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
	return c
}
