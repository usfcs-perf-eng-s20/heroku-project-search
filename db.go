package main

const (
	host     = "ec2-184-72-236-57.compute-1.amazonaws.com"
	port     = 5432
	user     = "wrwcqifhvfkkjw"
	password = "08f2594b47185df91e8bd0513405b8f5f4831089dc59c27387cd89e465b96015"
	dbname   = "d3hpcrkvokd5i"
	sslmode  = "require"
)

type Detail struct {
	Title  string
	Studio string
	Price  string
	Rating string
	Year   string
	Genre  string
	Upc    string
	ID     int64
}
