package main

import (
	"fmt"
	"log"
	"net/http"
	"wxcloudrun-golang/db"
	"wxcloudrun-golang/service"
)

func main() {
	if err := db.Init(); err != nil {
		panic(fmt.Sprintf("mysql init failed with %+v", err))
	}

	http.HandleFunc("/", service.IndexHandler)
	http.HandleFunc("/api/count", service.CounterHandler)
	http.HandleFunc("/api/chat", service.Chat)
	http.HandleFunc("/api/updateScore", service.UploadScore)
	http.HandleFunc("/api/score", service.GetScore)
	http.HandleFunc("/api/getUserId", service.GetUserId)

	log.Fatal(http.ListenAndServe(":80", nil))
}
