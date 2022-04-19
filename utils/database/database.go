package database

import (
	"database/sql"
	"golang-practise-project/configs"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var db *gorm.DB
var sqlDB *sql.DB

func InitDatabase() error {
	// open the database
	dsn := "root:lila_db@tcp(127.0.0.1:13306)/lila_db?charset=utf8mb4&parseTime=True&loc=Local"

	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		return err
	}

	// creating the SQL DB Handler
	sqlDB, err = db.DB()
	if err != nil {
		return err
	}

	// try to ping
	err = sqlDB.Ping()
	if err != nil {
		return err
	}

	// now set the configurations
	sqlDB.SetMaxOpenConns(configs.MaxOpenConnections)
	sqlDB.SetMaxIdleConns(configs.MaxIdleConnections)
	sqlDB.SetConnMaxIdleTime(time.Duration(configs.ConnectionMaxLifetimeInSeconds) * time.Second)
	sqlDB.SetConnMaxLifetime(time.Duration(configs.ConnectionMaxIdleTimeInSeconds) * time.Second)

	return nil
}

// returns the gorm handler
func Get() *gorm.DB {
	return db
}

// to close the connection
func Close() error {
	return sqlDB.Close()
}
