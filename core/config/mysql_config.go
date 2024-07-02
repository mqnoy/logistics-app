package config

import (
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func InitMySQLDatabase(cfg Configuration) *gorm.DB {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.LogLevel(cfg.Database.LogLevel),
			IgnoreRecordNotFoundError: true,
			ParameterizedQueries:      true,
			Colorful:                  false,
		},
	)

	db, err := gorm.Open(mysql.Open(cfg.Database.Mysql.DSN()), &gorm.Config{
		Logger: newLogger,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
			NoLowerCase:   true,
		},
	})

	if err != nil {
		log.Fatalf("Error creating database connection err: %v", err)
	}

	log.Printf("Database is connected to %s:%d", cfg.Database.Mysql.Host, cfg.Database.Mysql.Port)

	return db
}
