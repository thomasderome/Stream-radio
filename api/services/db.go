package services

import (
	"fmt"
	"radio_stream/model"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

var DB *sqlx.DB

func InitDB() {
	var err error
	DB, err = sqlx.Connect("sqlite3", "db.sqlite")

	if err != nil {
		panic(err)
	}

	fmt.Println("DB connected")
	DB.MustExec(model.SCHEMA_DATABASE)
	return
}
