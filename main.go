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

type idList struct {
	IdList []int64 `json:"ids" binding:"required"`
}

func search(c *gin.Context) {
	// get search keywords
	keyword := c.Query("keyword")
	log.Printf("keyword is: %s\n", keyword)

	var dvds []Detail

	// connect to database
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Println(err)
		c.JSON(500, "internal server error")
		return
	}

	sqlStatement := `SELECT * FROM dvds WHERE LOWER(title) LIKE '%' || $1 || '%' ;`
	rows, err := db.Query(sqlStatement, keyword)

	if err != nil {
		log.Println(err)
		c.JSON(500, "internal server error")
		return
	}

	for rows.Next() {
		var dvd Detail
		queryErr := rows.Scan(&dvd.Title, &dvd.Studio, &dvd.Price, &dvd.Rating, &dvd.Year,
			&dvd.Genre, &dvd.Upc, &dvd.ID)
		switch queryErr {
		case sql.ErrNoRows:
			return
		case nil:
			dvds = append(dvds, dvd)
			log.Println(dvd)
		default:
			log.Println(queryErr)
			c.JSON(500, "internal server error")
			return
		}
	}

	defer db.Close()
	log.Println(dvds)
	c.JSON(200, dvds)
}

func getMoviesByIDs(c *gin.Context) {
	// get query ids
	ids := idList{}
	var idList []int64
	var dvds []Detail

	// This reads c.Request.Body and stores the result into the context.
	if err := c.ShouldBindBodyWith(&ids, binding.JSON); err == nil {
		idList = ids.IdList
		log.Println(idList)
	}

	// connect to database
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Println(err)
		c.JSON(500, "internal server error")
		return
	}

	sqlStatement := `SELECT * FROM dvds WHERE id = any($1);`
	rows, err := db.Query(sqlStatement, pq.Array(idList))

	if err != nil {
		log.Println(err)
		c.JSON(500, "internal server error")
		return
	}

	for rows.Next() {
		var dvd Detail
		queryErr := rows.Scan(&dvd.Title, &dvd.Studio, &dvd.Price, &dvd.Rating, &dvd.Year,
			&dvd.Genre, &dvd.Upc, &dvd.ID)
		switch queryErr {
		case sql.ErrNoRows:
			log.Println("No rows were returned!")
			return
		case nil:
			dvds = append(dvds, dvd)
		default:
			log.Println(queryErr)
			c.JSON(500, "internal server error")
			return
		}
	}

	defer db.Close()
	log.Println(dvds)
	c.JSON(200, dvds)
}

func getMovieByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Query("id"), 10, 64)

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Println(err)
		c.JSON(500, "internal server error")
		return
	}

	sqlStatement := `SELECT * FROM dvds WHERE id=$1;`
	var dvd Detail
	row := db.QueryRow(sqlStatement, id)

	queryErr := row.Scan(&dvd.Title, &dvd.Studio, &dvd.Price, &dvd.Rating, &dvd.Year,
		&dvd.Genre, &dvd.Upc, &dvd.ID)

	log.Println(dvd)

	switch queryErr {
	case sql.ErrNoRows:
		log.Println("No rows were returned!")
		return
	case nil:
		log.Println("success")
	default:
		log.Println(queryErr)
		c.JSON(500, "internal server error")
		return
	}

	defer db.Close()
	c.JSON(200, dvd)
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
