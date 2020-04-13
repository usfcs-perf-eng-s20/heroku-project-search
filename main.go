package main

import (
	"database/sql"
	"github.com/lib/pq"
	newrelic "github.com/newrelic/go-agent"
	"log"

	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
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

var loaderIOPath = "loaderio-b5db249a1a78a9873b364017c18a4edb.txt"
var loaderIOUrlPath = "/loaderio-b5db249a1a78a9873b364017c18a4edb.txt"
var hostName = "https://perfeng-search.herokuapp.com/"
var serviceName = "search"
var APP_NAME = "perfeng_search"
var NEWRELIC_KEY = "310ab9b1832faccbf5d072e2a828f4535fedNRAL"

func updateFlag(value string, statusOk bool, varName string) {
	if statusOk {
		_, err := strconv.ParseBool(value)
		if err == nil {
			os.Setenv(varName, value)
		} else {
			log.Println("Invalid config value for ", varName)
		}
	}
}

func config(c *gin.Context) {
	start := time.Now()
	analytics, analyticsOk := c.GetQuery("analytics")
	faves, favesOk := c.GetQuery("favs")
	login, loginOk := c.GetQuery("login")

	updateFlag(analytics, analyticsOk, "STORE_ANALYTICS")
	updateFlag(faves, favesOk, "CALL_FAVS")
	updateFlag(login, loginOk, "CALL_LOGIN")

	logMessage := getConfigLogMessage(hostName, serviceName, "GET", "/config", analytics, faves, login, 0, time.Since(start).Nanoseconds()/1000000, "config flag updated")
	log.Println(logMessage)

	c.HTML(http.StatusOK, "config.tmpl.html", gin.H{
		"analytics_status": os.Getenv("STORE_ANALYTICS"),
		"faves_status":     os.Getenv("CALL_FAVS"),
		"login_status":     os.Getenv("CALL_LOGIN"),
	})
}

func search(c *gin.Context) {
	start := time.Now()
	// get search keywords
	keyword := c.Query("keyword")
	var event *analyticsEvent
	queryKey := getCacheKey("/search", keyword)

	mutex.RLock()
	resp, hit := cache.Get(queryKey)
	mutex.RUnlock()

	if hit {
		c.JSON(200, gin.H{
			"success": "true",
			"results": resp,
		})
		event = getEvent("/search", time.Since(start).Nanoseconds()/1000, "200", true, start)
		go postEvent(event)

		logMessage := getLogMessage(hostName, serviceName, "GET", "/search", keyword, 0, time.Since(start).Nanoseconds()/1000000, 1, "fetch result from cache")
		log.Println(logMessage)
		return
	}

	var dvds []Detail

	db, err := getDbConn()
	if err != nil {
		logMessage := getLogMessage(hostName, serviceName, "GET", "/search", keyword, 1, time.Since(start).Nanoseconds()/1000000, 0, "db connection error")
		log.Println(logMessage)
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

	sqlStatement := `SELECT * FROM dvds WHERE LOWER(title) LIKE $1 || '%' LIMIT 50;`
	rows, err := db.Query(sqlStatement, strings.ToLower(keyword))

	if err != nil {
		logMessage := getLogMessage(hostName, serviceName, "GET", "/search", keyword, 1, time.Since(start).Nanoseconds()/1000000, 0, "error querying db")
		log.Println(logMessage)
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
		queryErr := rows.Scan(&dvd.ID, &dvd.Title, &dvd.Studio, &dvd.Price, &dvd.Rating, &dvd.Year,
			&dvd.Genre, &dvd.Upc)
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
			logMessage := getLogMessage(hostName, serviceName, "GET", "/search", keyword, 1, time.Since(start).Nanoseconds()/1000000, 1, "error querying db")
			log.Println(logMessage)
			log.Println(err)
			c.JSON(500, gin.H{
				"success": "false",
				"message": "internal server error",
			})
			event = getEvent("/search", time.Since(start).Nanoseconds()/1000, "500", false, start)
			go postEvent(event)
			return
		}
	}

	logMessage := getLogMessage(hostName, serviceName, "GET", "/search", keyword, 0, time.Since(start).Nanoseconds()/1000000, 0, "fetch result from db")
	log.Println(logMessage)
	c.JSON(200, gin.H{
		"success": "true",
		"results": dvds,
	})
	event = getEvent("/search", time.Since(start).Nanoseconds()/1000, "200", true, start)
	go postEvent(event)

	mutex.Lock()
	cache.Add(queryKey, dvds)
	mutex.Unlock()

}

func getMoviesByIDs(c *gin.Context) {
	start := time.Now()
	var event *analyticsEvent
	ids := c.Query("ids")
	queryKey := getCacheKey("/getMoviesByIDs", ids)

	mutex.RLock()
	resp, hit := cache.Get(queryKey)
	mutex.RUnlock()

	if hit {
		c.JSON(200, gin.H{
			"success": "true",
			"results": resp,
		})
		event = getEvent("/getMoviesByIds", time.Since(start).Nanoseconds()/1000, "200", true, start)
		go postEvent(event)

		logMessage := getLogMessage(hostName, serviceName, "GET", "/getMoviesByIDs", ids, 0, time.Since(start).Nanoseconds()/1000000, 1, "fetch result from cache")
		log.Println(logMessage)
		return
	}
	var idList []int
	var dvds []Detail

	idsString := strings.Split(ids, ",")

	for i := 0; i < len(idsString); i++ {
		id, err := strconv.Atoi(strings.TrimSpace(idsString[i]))
		if err != nil {
			logMessage := getLogMessage(hostName, serviceName, "GET", "/getMoviesByIds", ids, 1, time.Since(start).Nanoseconds()/1000000, 0, "bad request")
			log.Println(logMessage)
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
		logMessage := getLogMessage(hostName, serviceName, "GET", "/getMoviesByIds", ids, 1, time.Since(start).Nanoseconds()/1000000, 0, "db connection error")
		log.Println(logMessage)
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
		logMessage := getLogMessage(hostName, serviceName, "GET", "/getMoviesByIds", ids, 1, time.Since(start).Nanoseconds()/1000000, 0, "error querying db")
		log.Println(logMessage)
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
		queryErr := rows.Scan(&dvd.ID, &dvd.Title, &dvd.Studio, &dvd.Price, &dvd.Rating, &dvd.Year,
			&dvd.Genre, &dvd.Upc)
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
			logMessage := getLogMessage(hostName, serviceName, "GET", "/getMoviesByIds", ids, 1, time.Since(start).Nanoseconds()/1000000, 0, "error querying db")
			log.Println(logMessage)
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

	logMessage := getLogMessage(hostName, serviceName, "GET", "/getMoviesByIds", ids, 0, time.Since(start).Nanoseconds()/1000000, 0, "fetch result from db")
	log.Println(logMessage)
	c.JSON(200, gin.H{
		"success": "true",
		"results": dvds,
	})
	event = getEvent("/getMoviesByIds", time.Since(start).Nanoseconds()/1000, "200", true, start)
	go postEvent(event)

	mutex.Lock()
	cache.Add(queryKey, dvds)
	mutex.Unlock()
}

func getMovieByID(c *gin.Context) {
	start := time.Now()
	idStr := c.Query("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	var event *analyticsEvent
	queryKey := getCacheKey("/getMovieByID", strconv.Itoa(int(id)))

	mutex.RLock()
	resp, hit := cache.Get(queryKey)
	mutex.RUnlock()

	if hit {
		c.JSON(200, gin.H{
			"success": "true",
			"results": resp,
		})
		event = getEvent("/getMovieById", time.Since(start).Nanoseconds()/1000, "200", true, start)
		go postEvent(event)

		logMessage := getLogMessage(hostName, serviceName, "GET", "/getMovieById", idStr, 0, time.Since(start).Nanoseconds()/1000000, 1, "fetch result from cache")
		log.Println(logMessage)
		return
	}

	db, err := getDbConn()
	if err != nil {
		logMessage := getLogMessage(hostName, serviceName, "GET", "/getMovieById", idStr, 1, time.Since(start).Nanoseconds()/1000000, 0, "db connection error")
		log.Println(logMessage)
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

	queryErr := row.Scan(&dvd.ID, &dvd.Title, &dvd.Studio, &dvd.Price, &dvd.Rating, &dvd.Year,
		&dvd.Genre, &dvd.Upc)

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
		logMessage := getLogMessage(hostName, serviceName, "GET", "/getMovieById", idStr, 1, time.Since(start).Nanoseconds()/1000000, 0, "error querying db")
		log.Println(logMessage)
		log.Println(queryErr)
		c.JSON(500, gin.H{
			"success": "false",
			"message": "internal server error",
		})
		event = getEvent("/getMovieById", time.Since(start).Nanoseconds()/1000, "500", false, start)
		go postEvent(event)
		return
	}

	logMessage := getLogMessage(hostName, serviceName, "GET", "/getMovieById", idStr, 0, time.Since(start).Nanoseconds()/1000000, 0, "fetch result from db")
	log.Println(logMessage)
	c.JSON(200, dvd)
	event = getEvent("/getMovieById", time.Since(start).Nanoseconds()/1000, "200", true, start)
	go postEvent(event)

	mutex.Lock()
	cache.Add(queryKey, dvd)
	mutex.Unlock()
}

func newRelicMiddleware(appName string, license string) gin.HandlerFunc {

	if appName == "" || license == "" {
		return func(c *gin.Context) {}
	}

	config := newrelic.NewConfig(appName, license)
	app, err := newrelic.NewApplication(config)

	if err != nil {
		panic(err)
	}

	return func(c *gin.Context) {
		txn := app.StartTransaction(c.Request.URL.Path, c.Writer, c.Request)
		defer txn.End()
		c.Next()
	}
}

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(cors.Default())
	router.Use(newRelicMiddleware(APP_NAME, NEWRELIC_KEY))
	router.LoadHTMLGlob("templates/*.tmpl.html")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl.html", nil)
	})

	router.GET("/getMovieById", getMovieByID)
	router.GET("/search", search)
	router.GET("/getMoviesByIds", getMoviesByIDs)
	router.GET("/config", config)
	router.StaticFile(loaderIOUrlPath, loaderIOPath)

	router.Run(":" + port)
}
