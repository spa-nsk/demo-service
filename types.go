package main

import "time"

type ResponseData struct {
	ResponseCount uint64
	TimeResponse  time.Duration
}

type ClientData struct {
	Title string
	Data  map[string]ResponseData
}

type responseStruct struct {
	Error error
	Items []responseItem
}

type responseItem struct {
	Host string
	Url  string
}
