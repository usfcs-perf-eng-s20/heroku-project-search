package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin/binding"
	"github.com/lib/pq"
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

type idList struct {
	IdList []int64 `json:"ids" binding:"required"`
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
	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	router := gin.New()
	router.Use(gin.Logger())
	router.LoadHTMLGlob("templates/*.tmpl.html")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl.html", nil)
	})

	router.GET("/getMovieById", getMovieByID)
	router.GET("/search", search)
	router.POST("/getMoviesByIds", getMoviesByIDs)

	router.Run(":" + port)
}

func search(c *gin.Context) {
	c.String(200, "Search endpoint")
}

func getMoviesByIDs(c *gin.Context) {
	// get query ids
	ids := idList{}
	idList := []int64{}
	dvds := []Detail{}
	// This reads c.Request.Body and stores the result into the context.
	if err := c.ShouldBindBodyWith(&ids, binding.JSON); err == nil {
		idList = ids.IdList
		fmt.Println(idList)
	}

	// connect to database
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Successfully connected to DB")

	sqlStatement := `SELECT * FROM dvds WHERE id = any($1);`
	rows, err := db.Query(sqlStatement, pq.Array(idList))

	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		var dvd Detail
		queryErr := rows.Scan(&dvd.Title, &dvd.Studio, &dvd.Price, &dvd.Rating, &dvd.Year,
			&dvd.Genre, &dvd.Upc, &dvd.ID)
		switch queryErr {
		case sql.ErrNoRows:
			fmt.Println("No rows were returned!")
			return
		case nil:
			dvds = append(dvds, dvd)
			fmt.Println(dvd)
		default:
			panic(queryErr)
		}
	}

	defer db.Close()
	fmt.Println(dvds)
	c.JSON(200, dvds)
}

func getMovieByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Query("id"), 10, 64)

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Successfully connected to DB")

	sqlStatement := `SELECT * FROM dvds WHERE id=$1;`
	var dvd Detail
	row := db.QueryRow(sqlStatement, id)

	queryErr := row.Scan(&dvd.Title, &dvd.Studio, &dvd.Price, &dvd.Rating, &dvd.Year,
		&dvd.Genre, &dvd.Upc, &dvd.ID)

	fmt.Println(dvd)

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
	c.JSON(200, dvd)
}
