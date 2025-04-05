package db

import (
	"fmt"
	"time"

	"github.com/ArjunMalhotra/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// MysqlDB instant to pass to handlers
type MysqlDB struct {
	DB *gorm.DB
}

// NewMysqlDB will return a valid connection to Mysql DB Session
func NewMysqDB(cfg *config.Config) (*MysqlDB, error) {
	dns := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.MySQL.MysqlUser,
		cfg.MySQL.MysqlPassword,
		cfg.MySQL.MysqlHost,
		cfg.MySQL.MysqlPort,
		cfg.MySQL.MysqlDBName,
	)

	fmt.Println("DNS", dns)

	db, err := gorm.Open(mysql.Open(dns), &gorm.Config{
		PrepareStmt:                              true,
		DisableForeignKeyConstraintWhenMigrating: true,
		SkipDefaultTransaction:                   true, // Disable automatic transactions for read-only operations

		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "admetric_",
			SingularTable: true,
			NoLowerCase:   true,
		},
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// Set connection pool settings
	sqlDB.SetMaxIdleConns(10)                 // Maximum idle connections
	sqlDB.SetMaxOpenConns(1000)               // Maximum open connections
	sqlDB.SetConnMaxLifetime(5 * time.Minute) // Maximum connection lifetime

	dbc := &MysqlDB{
		DB: db,
	}

	return dbc, nil
}

// Migrate when you change your model, called from main only
func (db *MysqlDB) Migrate() error {
	return nil
}
