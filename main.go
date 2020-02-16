package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/heroku/x/hmetrics/onload"
	_ "github.com/lib/pq"
)

// local host
// const (
// 	host     = "localhost"
// 	port     = 5432
// 	user     = "zini"
// 	password = "zini"
// 	dbname   = "perfsearch"
// 	sslmode  = "disable"
// )

// heroku server
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
	price  string
	rating string
	year   string
	genre  string
	upc    string
	id     int64
}

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
		//log.Fatal("$PORT must be set")
	}

	router := gin.New()
	router.Use(gin.Logger())
	router.LoadHTMLGlob("templates/*.tmpl.html")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl.html", nil)
	})

	router.GET("/helloworld", func(c *gin.Context) {
		c.HTML(http.StatusOK, "helloworld.tmpl.html", nil)
	})

	router.GET("/foo", foo)

	router.Run(":" + port)
}

func foo(c *gin.Context) {
	id, err := strconv.ParseInt(c.Query("id"), 10, 64)

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}

	sqlStatement := `SELECT * FROM dvds WHERE id=$1;`
	var dvd Detail
	row := db.QueryRow(sqlStatement, id)

	queryErr := row.Scan(&dvd.Title, &dvd.Studio, &dvd.price, &dvd.rating, &dvd.year,
		&dvd.genre, &dvd.upc, &dvd.id)

	switch queryErr {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned!")
		return
	case nil:
		fmt.Println(dvd)
	default:
		panic(queryErr)
	}

	defer db.Close()
	fmt.Println("Successfully connected to DB")

	c.JSON(200, dvd)
}
