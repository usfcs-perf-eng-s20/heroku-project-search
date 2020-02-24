package main

import (
	"database/sql"
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

	db, err := getDbConn()
	if err != nil {
		log.Println(err)
		c.JSON(500, "internal server error")
		return
	}
	defer db.Close()

	sqlStatement := `SELECT * FROM dvds WHERE LOWER(title) LIKE '%' || $1 || '%' ;`
	rows, err := db.Query(sqlStatement, keyword)

	if err != nil {
		log.Println(err)
		c.JSON(500, "internal server error")
		return
	}

	if !rows.Next() {
		c.JSON(200, "[]")
		return
	}

	for rows.Next() {
		var dvd Detail
		queryErr := rows.Scan(&dvd.Title, &dvd.Studio, &dvd.Price, &dvd.Rating, &dvd.Year,
			&dvd.Genre, &dvd.Upc, &dvd.ID)
		switch queryErr {
		case sql.ErrNoRows:
			c.JSON(200, "[]")
			return
		case nil:
			dvds = append(dvds, dvd)
		default:
			log.Println(queryErr)
			c.JSON(500, "internal server error")
			return
		}
	}

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

	db, err := getDbConn()
	if err != nil {
		log.Println(err)
		c.JSON(500, "internal server error")
		return
	}
	defer db.Close()

	sqlStatement := `SELECT * FROM dvds WHERE id = any($1);`
	rows, err := db.Query(sqlStatement, pq.Array(idList))

	if err != nil {
		log.Println(err)
		c.JSON(500, "internal server error")
		return
	}

	if !rows.Next() {
		c.JSON(200, "[]")
		return
	}

	for rows.Next() {
		var dvd Detail
		queryErr := rows.Scan(&dvd.Title, &dvd.Studio, &dvd.Price, &dvd.Rating, &dvd.Year,
			&dvd.Genre, &dvd.Upc, &dvd.ID)
		switch queryErr {
		case sql.ErrNoRows:
			c.JSON(200, "{}")
			return
		case nil:
			dvds = append(dvds, dvd)
		default:
			log.Println(queryErr)
			c.JSON(500, "internal server error")
			return
		}
	}

	c.JSON(200, dvds)
}

func getMovieByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Query("id"), 10, 64)

	db, err := getDbConn()
	if err != nil {
		log.Println(err)
		c.JSON(500, "internal server error")
		return
	}
	defer db.Close()

	sqlStatement := `SELECT * FROM dvds WHERE id=$1;`
	var dvd Detail
	row := db.QueryRow(sqlStatement, id)

	queryErr := row.Scan(&dvd.Title, &dvd.Studio, &dvd.Price, &dvd.Rating, &dvd.Year,
		&dvd.Genre, &dvd.Upc, &dvd.ID)

	switch queryErr {
	case sql.ErrNoRows:
		c.JSON(200, "{}")
		return
	case nil:
	default:
		log.Println(queryErr)
		c.JSON(500, "internal server error")
		return
	}

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
