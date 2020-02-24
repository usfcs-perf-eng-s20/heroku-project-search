package main

import "os"

var (
	host     = os.Getenv("DATABASE_HOST")
	port     = 5432
	user     = os.Getenv("DATABASE_USER")
	password = os.Getenv("DATABASE_PSWD")
	dbname   = os.Getenv("DATABASE_NAME")
	sslmode  = "require"
)
