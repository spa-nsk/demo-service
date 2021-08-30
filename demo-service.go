package main

import (
	"fmt"
	"log"
	"net/http"
	//	"github.com/gomodule/redigo/redis"
)

func main() {
	/*
		pool = &redis.Pool{
			MaxIdle:     10,
			IdleTimeout: 240 * time.Second,
			Dial: func() (redis.Conn, error) {
				return redis.Dial("tcp", "localhost:6379")
			},
		}
	*/
	mux := http.NewServeMux()
	mux.HandleFunc("/sites", searchSites)

	log.Println("Listening on :8080...")
	http.ListenAndServe(":8080", mux)
}

/*
func searchSites(w http.ResponseWriter, r *http.Request) {
	// Unless the request is using the POST method, return a 405
	// Method Not Allowed response.
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, http.StatusText(405), 405)
		return
	}

	// Retrieve the id from the POST request body. If there is no
	// parameter named "id" in the request body then PostFormValue()
	// will return an empty string. We check for this, returning a 400
	// Bad Request response if it's missing.
	search := r.PostFormValue("search")
	if search == "" {
		http.Error(w, http.StatusText(400), 400)
		return
	}

	if search == "demo" {
		fmt.Fprintf(w, "Демонстрация микросервиса\n")
	}

	// Validate that the id is a valid integer by trying to convert it,
	// returning a 400 Bad Request response if the conversion fails.
		if _, err := strconv.Atoi(search); err != nil {
			http.Error(w, http.StatusText(400), 400)
			return
		}
	// Call the IncrementLikes() function passing in the user-provided
	// id. If there's no album found with that id, return a 404 Not
	// Found response. In the event of any other errors, return a 500
	// Internal Server Error response.
		err := IncrementLikes(id)
		if err == ErrNoAlbum {
			http.NotFound(w, r)
			return
		} else if err != nil {
			http.Error(w, http.StatusText(500), 500)
			return
		}
	// Redirect the client to the GET /album route, so they can see the
	// impact their like has had.
	//	http.Redirect(w, r, "/album?id="+id, 303)
}
*/

func searchSites(w http.ResponseWriter, r *http.Request) {
	// Unless the request is using the GET method, return a 405 'Method
	// Not Allowed' response.
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, http.StatusText(405), 405)
		return
	}

	// Retrieve the id from the request URL query string. If there is
	// no id key in the query string then Get() will return an empty
	// string. We check for this, returning a 400 Bad Request response
	// if it's missing.
	id := r.URL.Query().Get("search")
	if id == "" {
		http.Error(w, http.StatusText(400), 400)
		return
	}

	if id == "demo" {
		fmt.Fprintf(w, "Деменстрация microservice\n")
	} else {
		fmt.Fprintf(w, "/sites?search=%s\n", id)
	}

	// Validate that the id is a valid integer by trying to convert it,
	// returning a 400 Bad Request response if the conversion fails.
	/*
		if _, err := strconv.Atoi(id); err != nil {
			http.Error(w, http.StatusText(400), 400)
			return
		}
	*/
	// Call the FindAlbum() function passing in the user-provided id.
	// If there's no matching album found, return a 404 Not Found
	// response. In the event of any other errors, return a 500
	// Internal Server Error response.
	/*
		bk, err := FindAlbum(id)
		if err == ErrNoAlbum {
			http.NotFound(w, r)
			return
		} else if err != nil {
			http.Error(w, http.StatusText(500), 500)
			return
		}

		// Write the album details as plain text to the client.
		fmt.Fprintf(w, "%s by %s: £%.2f [%d likes] \n", bk.Title, bk.Artist, bk.Price, bk.Likes)
	*/

}

/*
func listPopular(w http.ResponseWriter, r *http.Request) {
	// Unless the request is using the GET method, return a 405 'Method Not
	// Allowed' response.
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, http.StatusText(405), 405)
		return
	}

	// Call the FindTopThree() function, returning a return a 500 Internal
	// Server Error response if there's any error.
	albums, err := FindTopThree()
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	// Loop through the 3 albums, writing the details as a plain text list
	// to the client.
	for i, ab := range albums {
		fmt.Fprintf(w, "%d) %s by %s: £%.2f [%d likes] \n", i+1, ab.Title, ab.Artist, ab.Price, ab.Likes)
	}
}
*/
