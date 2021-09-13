package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync/atomic"
	"time"

	"golang.org/x/net/context"
)

func searchSites(w http.ResponseWriter, r *http.Request) {
	timeOutRequest := time.Millisecond * time.Duration(atomic.LoadUint64(&TimeOutWork))
	start := time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), timeOutRequest)
	defer cancel()
	defer func() {
		end := time.Now()
		fmt.Println("Время выполнения запроса", end.Sub(start))
	}()
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, http.StatusText(405), 405)
		return
	}

	search := r.URL.Query().Get("search")
	if search == "" {
		http.Error(w, http.StatusText(400), 400)
		return
	}

	var defaultTtransport http.RoundTripper = &http.Transport{Proxy: nil}
	client := &http.Client{Transport: defaultTtransport}

	resp, err := client.Get(baseYandexURL + search)

	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	//func parseYandexResponse(response []byte) (res responseStruct)
	res := parseYandexResponse(body)

	if res.Error != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	s := make(map[string]ResponseData)

	for _, item := range res.Items {
		select {
		case <-ctx.Done():
			fmt.Println("Истекло время выполнения запроса (", timeOutRequest, ").")
			json.NewEncoder(w).Encode(s)
			return
		default:
			count, timeResponse := checkAvailability(item.Url)
			s[item.Host] = ResponseData{count, timeResponse}
			fmt.Println(item.Host, count, timeResponse)
		}
	}
	json.NewEncoder(w).Encode(s)
}
