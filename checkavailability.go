package main

import (
	"io/ioutil"
	"net"
	"net/http"
	"sync/atomic"
	"time"
)

func checkAvailability(url string) (uint64, time.Duration) {
	var i, index uint64
	countRequest := atomic.LoadUint64(&CountRequest)
	timeOutRequest := time.Millisecond * time.Duration(atomic.LoadUint64(&TimeOutRequest))
	timeResponse := time.Millisecond * 0

	ch := make(chan time.Duration)

	for i = 0; i < countRequest; i++ {
		go readUrl(url, timeOutRequest, ch)
	}

	for i = 0; i < countRequest; i++ {
		t := <-ch
		if t == time.Second*9999 {
			if index == 0 {
				index = i
			}
			continue
		}
		if t > timeResponse {
			timeResponse = t
		}
	}
	if index == 0 {
		return i, timeResponse
	}
	return index, timeResponse
}

func readUrl(url string, sec time.Duration, ch chan time.Duration) {
	var defaultTtransport http.RoundTripper = &http.Transport{
		Dial: (&net.Dialer{
			Timeout:   sec,
			KeepAlive: sec}).Dial,
		TLSHandshakeTimeout: sec}
	client := &http.Client{Transport: defaultTtransport}
	start := time.Now()
	resp, err := client.Get(url)

	if err != nil {
		ch <- time.Second * 9999
		//fmt.Println("ошибка client.Get(", url, ") ", err)
		return
	}
	if resp.StatusCode == 429 { //слишком много запросов
		ch <- time.Second * 9999
		//fmt.Println("Слишком много запросов", url)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	_ = body
	if err != nil {
		ch <- time.Second * 9999
		//fmt.Println("ошибка ioutil.ReadAll()", err)
		return
	}
	end := time.Now()
	ch <- end.Sub(start)
}
