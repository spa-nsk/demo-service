package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/sites", searchSites)

	log.Println("Listening on :8080...")
	http.ListenAndServe(":8080", mux)
}

func searchSites(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, http.StatusText(405), 405)
		return
	}

	id := r.URL.Query().Get("search")
	if id == "" {
		http.Error(w, http.StatusText(400), 400)
		return
	}

	if id == "demo" {
		fmt.Fprintf(w, "Демонстрация microservice\n")
	} else {
		fmt.Fprintf(w, "/sites?search=%s\n", id)
	}

}
