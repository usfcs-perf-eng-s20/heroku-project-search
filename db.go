package main

import (
	"database/sql"
	"fmt"
	"os"
)

var (
	host     = os.Getenv("DATABASE_HOST")
	port     = 5432
	user     = os.Getenv("DATABASE_USER")
	password = os.Getenv("DATABASE_PSWD")
	dbname   = os.Getenv("DATABASE_NAME")
	sslmode  = "require"
)

func getDbConn() (db *sql.DB, err error) {
	// connect to database
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode)
	db, err = sql.Open("postgres", psqlInfo)
	return
}
