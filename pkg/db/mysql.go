package db

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/ArjunMalhotra/config"
	"github.com/ArjunMalhotra/internal/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

const (
	ADS_PATH = "./assets/ads.json"
)

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

func (db *MysqlDB) Migrate() error {
	if err := db.DB.AutoMigrate(&model.Ad{}, &model.Click{}); err != nil {
		return err
	}
	return nil
}

func (db *MysqlDB) Seed() error {
	data, err := os.ReadFile(ADS_PATH)
	if err != nil {
		return fmt.Errorf("Failed to load ads from json file : %w", err)
	}
	var ads []model.Ad
	if err := json.Unmarshal(data, &ads); err != nil {
		return fmt.Errorf("failed to parse JSON: %w", err)
	}
	tx := db.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	//! seed
	if err := tx.Create(&ads).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to seed data -> ", err)
	}
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction -> ", err)
	}
	return nil
}
