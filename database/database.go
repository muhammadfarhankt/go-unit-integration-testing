package database

import (
	"database/sql"
	"userApiTest/model"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Db *gorm.DB
var MockDB *sql.DB

func AutoMigrate(db *gorm.DB) {
	db.AutoMigrate(model.User{})
}

func InitDb() {
	DSN := "host=localhost user=postgres password=0000 dbname=userapitest port=5432"

	gormDb, err := gorm.Open(postgres.Open(DSN))
	if err != nil {
		panic(err)
	}
	
	AutoMigrate(gormDb)
	Db = gormDb
}

func DbSet() sqlmock.Sqlmock {

	mockDB, mock, err := sqlmock.New()
	if err != nil {
		panic(err)
	}

	dialecter := postgres.New(postgres.Config{
		Conn: mockDB,
	})

	gormDb, err := gorm.Open(dialecter, &gorm.Config{})
	if err != nil {
		panic(err)
	}

	Db = gormDb
	MockDB = mockDB

	return mock
}
