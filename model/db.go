package model

import (
	"log/slog"
	"os"
	"time"

	"github.com/swartzfoundation/feedr/pkg/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB

// GetDB returns the database connection
func GetDB() *gorm.DB {
	if db == nil {
		slog.Error("db is not initialized")
		os.Exit(1)
	}
	return db
}

// ConnectDatabase connects to the database
func ConnectDatabase(dbConfig config.DBConfig) {
	logMode := logger.Silent
	if config.Config.DEBUG {
		logMode = logger.Warn
	}
	gormDB, err := gorm.Open(postgres.Open(dbConfig.DSN), &gorm.Config{
		Logger:                 logger.Default.LogMode(logMode),
		SkipDefaultTransaction: false,
	})
	if err != nil {
		slog.Error("db: connecting to database", "error", err.Error())
	}
	sqlDB, err := gormDB.DB()
	if err != nil {
		slog.Error("db: setting max idle connections", "error", err.Error())
		os.Exit(1)
	}
	sqlDB.SetMaxIdleConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)
	db = gormDB
}

func CloseDatabase() {
	d := GetDB()
	sqlDB, err := d.DB()
	if err != nil {
		slog.Error("db: getting database", "error", err.Error())
	}
	if err := sqlDB.Close(); err != nil {
		slog.Error("db: closing database", "error", err.Error())
	}
	slog.Warn("Database closed successfully")
}

func PingDatabase() error {
	sqlDB, err := db.DB()
	if err != nil {
		slog.Error("db: pinging database", "error", err.Error())
	}
	return sqlDB.Ping()
}

var tables = []interface{}{
	&User{},
	&Session{},
}

func Tables() []interface{} {
	return tables
}

// MigrateDatabase migrates the database
func MigrateDatabase() error {
	slog.Warn("db migrations started")
	if err := db.AutoMigrate(tables...); err != nil {
		slog.Error("db: migrating database", "error", err.Error())
		return err
	}

	slog.Warn("db migration complete")
	return nil
}
