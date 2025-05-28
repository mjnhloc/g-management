package infrastructure

import (
	"context"
	"database/sql"
	"os"

	"g-management/pkg/log"

	"github.com/joho/godotenv"
	mysqlDriver "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func GetGormConfig() *gorm.Config {
	return &gorm.Config{
		DisableAutomaticPing: true,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	}
}

// create a new instance of database
func NewDatabase() (db *gorm.DB, master *sql.DB, err error) {
	gormConfig := GetGormConfig()

	err = godotenv.Load()
	if err != nil {
		log.Error(context.Background(), "Error loading .env file", "error", err)
	}

	database, err := sql.Open("mysql", os.Getenv("DSN"))
	if err != nil {
		return nil, nil, err
	}

	database.SetMaxIdleConns(20)
	database.SetMaxOpenConns(200)

	driverMaster := mysqlDriver.New(mysqlDriver.Config{
		Conn:                      database,
		SkipInitializeWithVersion: true,
		DontSupportForShareClause: true,
	})

	db, err = gorm.Open(driverMaster, gormConfig)
	if err != nil {
		return nil, nil, err
	}

	err = PingCtx(context.Background(), db)
	if err != nil {
		return nil, nil, err
	}

	return db, database, nil
}

func PingCtx(ctx context.Context, db *gorm.DB) error {
	mysqlDB, err := db.DB()
	if err != nil {
		return err
	}

	return mysqlDB.PingContext(ctx)
}

func CloseDB(db *sql.DB) {
	if err := db.Close(); err != nil {
		log.Error(context.Background(), "Error closing the database", "error", err)
	}
}
