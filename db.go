package main

import (
	"database/sql"
	"fmt"
)

var (
	host     = "ec2-184-72-236-57.compute-1.amazonaws.com"
	port     = 5432
	user     = "wrwcqifhvfkkjw"
	password = "08f2594b47185df91e8bd0513405b8f5f4831089dc59c27387cd89e465b96015"
	dbname   = "d3hpcrkvokd5i"
	sslmode  = "require"
)

//var (
//	host     = os.Getenv("DATABASE_HOST")
//	port     = 5432
//	user     = os.Getenv("DATABASE_USER")
//	password = os.Getenv("DATABASE_PSWD")
//	dbname   = os.Getenv("DATABASE_NAME")
//	sslmode  = "require"
//)

func getDbConn() (db *sql.DB, err error) {
	// connect to database
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode)
	db, err = sql.Open("postgres", psqlInfo)
	return
}
