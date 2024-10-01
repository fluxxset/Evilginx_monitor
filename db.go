package main

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

var dbPath = filepath.Join(os.Getenv("HOME"), ".evilginx_monitor", "record_tracker.db")

func initDB() {
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		db, err := sql.Open("sqlite3", dbPath)
		if err != nil {
			fmt.Println("Failed to create database:", err)
			return
		}
		defer db.Close()

		// Create table
		createTableSQL := `CREATE TABLE IF NOT EXISTS credentials (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            data TEXT NOT NULL
        );`
		_, err = db.Exec(createTableSQL)
		if err != nil {
			fmt.Println("Failed to create table:", err)
		}
	}
}
