package db

import (
	"fmt"
	"github.com/ladmakhi81/learnup/pkg/env"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"time"
)

type Database struct {
	Core   *gorm.DB
	config *env.EnvConfig
}

func NewDatabase(config *env.EnvConfig) *Database {
	return &Database{
		config: config,
	}
}

func (db *Database) Connect() error {
	connection := db.getConnection()
	pgDialector := postgres.Open(connection)
	coreDb, err := gorm.Open(pgDialector, &gorm.Config{})
	if err != nil {
		return err
	}

	// Connection Managing
	sqlDb, sqlDbErr := coreDb.DB()
	if sqlDbErr != nil {
		return sqlDbErr
	}
	sqlDb.SetMaxOpenConns(50)
	sqlDb.SetMaxIdleConns(10)
	sqlDb.SetConnMaxLifetime(time.Hour * 1)
	sqlDb.SetConnMaxLifetime(time.Minute * 30)

	// TODO: add models into this AuthMigrate
	//coreDb.AutoMigrate()
	db.Core = coreDb
	return nil
}

func (db *Database) getConnection() string {
	dbHost := db.config.MainDB.Host
	dbUser := db.config.MainDB.Username
	dbPassword := db.config.MainDB.Password
	dbName := db.config.MainDB.Name
	dbPort := db.config.MainDB.Port

	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		dbHost, dbUser, dbPassword, dbName, dbPort,
	)
}
