package repository

import (
	"database/sql"
	"fmt"
	"time"

	"testsocket/config"

	_ "github.com/go-sql-driver/mysql"
)

func NewDb() (*sql.DB, error) {
	config := config.NewConfig()

	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", config.MySQL["user"], config.MySQL["password"], config.MySQL["host"], config.MySQL["port"], config.MySQL["db"]))

	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}
	// See "Important settings" section. "user:password@/dbname"
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	return db, err
}
