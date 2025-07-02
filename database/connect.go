package database

import (
	"fmt"
	"log"

	"gitag.ir/armogroup/armo/services/reality/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect() *gorm.DB {

	var (
		host     = config.AppConfig.DatabaseHost
		username = config.AppConfig.DatabaseUsername
		password = config.AppConfig.DatabasePassword
		dbName   = config.AppConfig.DatabaseName
		port     = config.AppConfig.DatabasePort
		sslMode  = config.AppConfig.DatabaseSslMode
	)

	DSN := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d"+
		" sslmode=%s TimeZone=Asia/Tehran",
		host, username, password, dbName, port, sslMode)

	db, err := gorm.Open(postgres.Open(DSN), &gorm.Config{})

	if err != nil {
		log.Fatal(err)
	}

	return db
}
