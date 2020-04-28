package database

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var DSN = ""
var connection *sql.DB

func Initialize() {
	log.Print("Initializing Database...")
	var DB_USER = os.Getenv("RAYMOND_DB_USER")
	var DB_PASS = os.Getenv("RAYMOND_DB_PASS")
	var DB_HOST = os.Getenv("RAYMOND_DB_HOST")
	var DB_NAME = os.Getenv("RAYMOND_DB_NAME")

	error := false
	if DB_USER == "" {
		log.Printf("RAYMOND_DB_USER environment variable not set")
		error = true
	}
	if DB_PASS == "" {
		log.Printf("RAYMOND_DB_PASS environment variable not set")
		error = true
	}
	if DB_NAME == "" {
		log.Printf("RAYMOND_DB_NAME environment variable not set")
		error = true
	}
	if DB_HOST == "" {
		log.Printf("RAYMOND_DB_HOST environment variable not set")
		error = true
	}

	if error {
		log.Fatal("Missing Database login information. Not Starting.")
	}

	DSN = DB_USER + ":" + DB_PASS + "@tcp(" + DB_HOST + ":3306)/" + DB_NAME

	// Open DB connection
	db, err := sql.Open("mysql", DSN)
	if err != nil {
		log.Fatalf("Failed to initiate database connection. Not Starting. %s", err.Error())
	}

	// make sure our connection is available
	err = db.Ping()
	if err != nil {
		log.Fatalf("Failed to ping database connection. Not Starting. %s", err.Error())
	}

	SetConnection(db)

	log.Print("done")
}

func SetConnection(db *sql.DB) {
	connection = db
}

func Connection() (database *sql.DB) {
	return connection
}
