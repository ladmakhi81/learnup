package db

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"time"
)

type Database struct {
	Core *gorm.DB
}

func NewDatabase() *Database {
	return &Database{}
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
	//TODO: replace this hard coded values with config values
	dbHost := "learnup_main_database_service"
	dbUser := "postgres"
	dbPassword := "postgres"
	dbName := "learnup_db"
	dbPort := 5432

	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		dbHost, dbUser, dbPassword, dbName, dbPort,
	)
}
