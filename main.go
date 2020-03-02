package main

import (
	"database/sql"
	"fmt"
	"github.com/lib/pq"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

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
	start := time.Now()
	// get search keywords
	keyword := c.Query("keyword")
	var event *analyticsEvent

	var dvds []Detail

	db, err := getDbConn()
	if err != nil {
		log.Println(err)
		c.JSON(500, gin.H{
			"success": "false",
			"message": "internal server error",
		})
		event = getEvent("/search", time.Since(start).Nanoseconds()/1000, "500", false, start)
		go postEvent(event)
		return
	}
	defer db.Close()

	sqlStatement := `SELECT * FROM dvds WHERE LOWER(title) LIKE '%' || $1 || '%' LIMIT 50;`
	fmt.Println(strings.ToLower(keyword))
	rows, err := db.Query(sqlStatement, strings.ToLower(keyword))

	if err != nil {
		log.Println(err)
		c.JSON(500, gin.H{
			"success": "false",
			"message": "internal server error",
		})
		event = getEvent("/search", time.Since(start).Nanoseconds()/1000, "500", false, start)
		go postEvent(event)
		return
	}

	for rows.Next() {
		var dvd Detail
		queryErr := rows.Scan(&dvd.Title, &dvd.Studio, &dvd.Price, &dvd.Rating, &dvd.Year,
			&dvd.Genre, &dvd.Upc, &dvd.ID)
		switch queryErr {
		case sql.ErrNoRows:
			c.JSON(200, gin.H{
				"success": "true",
				"results": "[]",
			})
			event = getEvent("/search", time.Since(start).Nanoseconds()/1000, "200", true, start)
			go postEvent(event)
			return
		case nil:
			dvds = append(dvds, dvd)
		default:
			log.Println(queryErr)
			c.JSON(500, gin.H{
				"success": "false",
				"message": "internal server error",
			})
			event = getEvent("/search", time.Since(start).Nanoseconds()/1000, "500", false, start)
			go postEvent(event)
			return
		}
	}

	c.JSON(200, gin.H{
		"success": "true",
		"results": dvds,
	})
	event = getEvent("/search", time.Since(start).Nanoseconds()/1000, "200", true, start)
	go postEvent(event)
}

func getMoviesByIDs(c *gin.Context) {
	start := time.Now()
	var idList []int
	var dvds []Detail
	var event *analyticsEvent

	idsString := strings.Split(c.Query("ids"), ",")

	for i := 0; i < len(idsString); i++ {
		id, err := strconv.Atoi(strings.TrimSpace(idsString[i]))
		if err != nil {
			log.Println(err)
			c.JSON(400, gin.H{
				"success": "false",
				"message": "bad request",
			})
			event = getEvent("/getMoviesByIds", time.Since(start).Nanoseconds()/1000, "400", false, start)
			go postEvent(event)
			return
		}
		idList = append(idList, id)
	}

	db, err := getDbConn()
	if err != nil {
		log.Println(err)
		c.JSON(500, gin.H{
			"success": "false",
			"message": "internal server error",
		})
		event = getEvent("/getMoviesByIds", time.Since(start).Nanoseconds()/1000, "500", false, start)
		go postEvent(event)
		return
	}
	defer db.Close()

	sqlStatement := `SELECT * FROM dvds WHERE id = any($1);`
	rows, err := db.Query(sqlStatement, pq.Array(idList))

	if err != nil {
		log.Println(err)
		c.JSON(500, gin.H{
			"success": "false",
			"message": "internal server error",
		})
		event = getEvent("/getMoviesByIds", time.Since(start).Nanoseconds()/1000, "500", false, start)
		go postEvent(event)
		return
	}

	for rows.Next() {
		var dvd Detail
		queryErr := rows.Scan(&dvd.Title, &dvd.Studio, &dvd.Price, &dvd.Rating, &dvd.Year,
			&dvd.Genre, &dvd.Upc, &dvd.ID)
		switch queryErr {
		case sql.ErrNoRows:
			c.JSON(200, gin.H{
				"success": "true",
				"results": "[]",
			})
			event = getEvent("/getMoviesByIds", time.Since(start).Nanoseconds()/1000, "200", true, start)
			go postEvent(event)
			return
		case nil:
			dvds = append(dvds, dvd)
		default:
			log.Println(queryErr)
			c.JSON(500, gin.H{
				"success": "false",
				"message": "internal server error",
			})
			event = getEvent("/getMoviesByIds", time.Since(start).Nanoseconds()/1000, "500", false, start)
			go postEvent(event)
			return
		}
	}

	c.JSON(200, gin.H{
		"success": "true",
		"results": dvds,
	})
	event = getEvent("/getMoviesByIds", time.Since(start).Nanoseconds()/1000, "200", true, start)
	go postEvent(event)
}

func getMovieByID(c *gin.Context) {
	start := time.Now()
	id, err := strconv.ParseInt(c.Query("id"), 10, 64)
	var event *analyticsEvent

	db, err := getDbConn()
	if err != nil {
		log.Println(err)
		c.JSON(500, gin.H{
			"success": "false",
			"message": "internal server error",
		})
		event = getEvent("/getMovieById", time.Since(start).Nanoseconds()/1000, "500", false, start)
		go postEvent(event)
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
		c.JSON(200, gin.H{
			"success": "true",
			"results": "[]",
		})
		event = getEvent("/getMovieById", time.Since(start).Nanoseconds()/1000, "200", true, start)
		go postEvent(event)
		return
	case nil:
	default:
		log.Println(queryErr)
		c.JSON(500, gin.H{
			"success": "false",
			"message": "internal server error",
		})
		event = getEvent("/getMovieById", time.Since(start).Nanoseconds()/1000, "500", false, start)
		go postEvent(event)
		return
	}

	c.JSON(200, dvd)
	event = getEvent("/getMovieById", time.Since(start).Nanoseconds()/1000, "200", true, start)
	go postEvent(event)
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
	router.GET("/getMoviesByIds", getMoviesByIDs)

	router.Run(":" + port)
}
