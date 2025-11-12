package databaseconn

import (
	"fmt"
	"log"

	"e-voting-service/data/configuration"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type LoaderError struct {
	Key      string
	KeyValue string
	Table    string
	message  string
}

func (e *LoaderError) Error() string {
	return fmt.Sprintf("Tabelle %s konnte nicht geladen werden: key %s; value: %s - %s", e.Table, e.Key, e.KeyValue, e.message)
}

var globalDB *gorm.DB = nil

func GetDB() (*gorm.DB, error) {
	if globalDB == nil {
		var config configuration.Config = configuration.GlobalConfig
		db, err := connectToMariaDB(config)
		if err != nil {
			return nil, err
		}
		globalDB = db
	}
	return globalDB, nil
}

func connectToMariaDB(config configuration.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8mb4&parseTime=True&loc=Local",
		config.Database.User,
		config.Database.Pwd,
		config.Database.Address,
		config.Database.Port,
		config.Database.Database)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Printf("Error connection to  Database %v", err)
		return nil, err
	}

	return db, nil
}

type WahlTest struct {
	WahlID int64 `gorm:"column:id"`
}

func TestConnection(config configuration.Config) error {
	db, err := connectToMariaDB(config)
	if err != nil {
		return err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	err = sqlDB.Ping()
	if err != nil {
		return err
	}

	// Testing a Query
	rows, err := sqlDB.Query("select * from election")
	if err != nil {
		return err
	}

	rows.Close()

	sqlDB.Close()

	return nil
}
