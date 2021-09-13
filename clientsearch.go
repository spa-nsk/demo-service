package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
)

func clientSearchSites(w http.ResponseWriter, r *http.Request) {

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

	resp, err := client.Get("http://127.0.0.1:8080/sites?search=" + search)

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
	//декодировать ответ и сформировать страничку ответа в структурированном виде
	var s map[string]ResponseData
	err = json.Unmarshal(body, &s)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	data := ClientData{Title: search,
		Data: s}
	tmpl, _ := template.ParseFiles("/opt/demo-service/view/search.html")
	err = tmpl.Execute(w, &data)
	if err != nil {
		fmt.Println("Ошибка парсинга шаблона", err)
		http.Error(w, http.StatusText(500), 500)
	}
}
