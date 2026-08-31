package services

import (
	"fmt"
	"radio_stream/model"
	"radio_stream/utils"
	"strings"

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

func QueryMakerInsertOneColumn(table string, value string, data utils.Set) (string, []any) {
	args := make([]any, 0, len(data))
	placeholders := make([]string, 0, len(data))

	for key, _ := range data {
		placeholders = append(placeholders, "(?)")
		args = append(args, key)
	}

	return fmt.Sprintf("INSERT OR IGNORE INTO %s (%s) VALUES %s", table, value, strings.Join(placeholders, ",")), args
}
