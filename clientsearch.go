package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net"
	"net/http"
	"sync/atomic"
	"time"
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

	var sec = time.Millisecond * time.Duration(atomic.LoadUint64(&TimeOutRequest))
	var defaultTtransport http.RoundTripper = &http.Transport{
		Dial: (&net.Dialer{
			Timeout:   sec,
			KeepAlive: sec}).Dial,
		TLSHandshakeTimeout: sec}
	client := &http.Client{Transport: defaultTtransport}

	resp, err := client.Get(ClientSearchPoint + search)

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
