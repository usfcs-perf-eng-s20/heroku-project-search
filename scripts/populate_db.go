package main

import (
	"bufio"
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"

	_ "github.com/lib/pq"
)

const MaxDbSize = 9999

func getDbConn() (db *sql.DB, err error) {
	// connect to database
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode)
	fmt.Println(psqlInfo)
	db, err = sql.Open("postgres", psqlInfo)
	return
}

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

func main() {
	// Path is specified from the CLI
	csvFile, _ := os.Open(os.Args[1])
	reader := csv.NewReader(bufio.NewReader(csvFile))
	var entries []Detail
	db, err := getDbConn()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	sqlStatement := `SELECT id FROM dvds;`
	rows, err := db.Query(sqlStatement)
	if err != nil {
		panic(err)
	}

	var tmpId int64
	ids := make(map[int64]int)
	for rows.Next() {
		err = rows.Scan(&tmpId)
		ids[tmpId] = 1
	}

	fmt.Println(ids)
	i := 0
	for i < MaxDbSize- len(ids) {
		line, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}

		id, err := strconv.Atoi(line[13])
		if err != nil {
			continue
		}

		// If id already exists
		if ids[int64(id)] != 0 {
			continue
		}

		// Some entries have missing information about year, skip those
		_, err = strconv.Atoi(line[8])
		if err != nil {
			continue
		}

		entries = append(entries, Detail{
			Title: line[0],
			Studio:  line[1],
			Price: line[6],
			Rating: line[7],
			Year: line[8],
			Genre: line[9],
			Upc: line[11],
			ID: int64(id),
		})
		i++
	}
	fmt.Println(len(entries))

	for i, entry := range entries {
		sqlStatement := `INSERT INTO dvds (title, studio, price, rating, year, genre, upc, id)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
		_, err = db.Exec(sqlStatement, entry.Title, entry.Studio, entry.Price, entry.Rating, entry.Year,
			entry.Genre, entry.Upc, entry.ID)
		if err != nil {
			panic(err)
		}
		fmt.Println("Inserted ", i)
	}
}