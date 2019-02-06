package db

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type Model struct {
	gorm.Model
}

func OpenSqlite3(args ...string) *gorm.DB {
	db, err := gorm.Open("sqlite3", args)
	if err != nil {
		panic("failed to connect sqlite3 database")
	}
	return db
}
func OpenMySql(args ...string) *gorm.DB {
	db, err := gorm.Open("mysql", args)
	if err != nil {
		panic("failed to connect mysql database")
	}
	return db
}
