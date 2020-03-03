package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
	"strconv"
)

type analyticsEvent struct {
	Method string `json:"method"`
	Path string `json:"path"`
	TimeMilis int64 `json:"processingTimeInMiliseconds"`
	Response string `json:"responseCode"`
	Service string `json:"serviceName"`
	Success bool `json:"success"`
	Timestamp string `json:"timestamp"`
	Username string `json:"username"`
}

var analyticsHost = fmt.Sprint(os.Getenv("ANALYTICS_URL"), "/saveEdr")
var storeAnalytics = os.Getenv("STORE_ANALYTICS") == "1"

func getEvent(path string, timeMillis int64, response string, success bool, timestamp time.Time) *analyticsEvent {
	e := analyticsEvent{
		Method:    "GET",
		Path:      path,
		TimeMilis: timeMillis,
		Response:  response,
		Service:   "search",
		Success:   success,
		Timestamp: strconv.Itoa(int(timestamp.UTC().Unix())),
		Username:  "",
	}
	return &e
}

func postEvent(e *analyticsEvent) {
	if !storeAnalytics {
		return
	}
	jsonEvent, err := json.Marshal(e)
	if err != nil {
		log.Println("internal error:", err)
		return
	}

	req, err := http.NewRequest("POST", analyticsHost, bytes.NewBuffer(jsonEvent))
	if err != nil {
		log.Println("internal error:", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("error:", err)
		return
	}
	defer resp.Body.Close()

	log.Println("updated analytics with status:", resp.Status)
}