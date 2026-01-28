package config

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB() {
	var err error
	DB, err = sql.Open("sqlite3", "./pdf.db")
	if err != nil {
		log.Fatal(err)
	}

	createTable := `
	CREATE TABLE IF NOT EXISTS pdf_files (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		filename TEXT,
		original_name TEXT,
		filepath TEXT,
		size INTEGER,
		status TEXT,
		created_at DATETIME,
		updated_at TIMESTAMP,
		deleted_at DATETIME
	);
	`

	if _, err = DB.Exec(createTable); err != nil {
		log.Fatal(err)
	}
}
