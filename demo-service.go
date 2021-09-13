package main

import (
	"log"
	"net/http"
)

func main() {
	loadConfig()
	mux := http.NewServeMux()
	mux.HandleFunc("/sites", searchSites)
	mux.HandleFunc("/sitesclient", clientSearchSites)

	log.Println("Слушаем порт :8080...")
	http.ListenAndServe(":8080", mux)
}
