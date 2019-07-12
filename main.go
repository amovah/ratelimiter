package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type Config struct {
	RequestPerDuration uint
	Duration uint
	TargetServer string
}

func ProxyResponse(origin http.Response, target http.ResponseWriter) {
	body, err := ioutil.ReadAll(origin.Body)
	if err != nil {
		fmt.Fprint(target, err)
		return
	}

	for key, value := range origin.Header {
		target.Header().Set(key, strings.Join(value, ", "))
	}

	target.WriteHeader(origin.StatusCode)
	target.Write(body)
}

func ProxyRequest(res http.ResponseWriter, req *http.Request) {
	if record[req.RemoteAddr] > uint64(config.RequestPerDuration) {
		http.Error(res, "You reach your limit", 429)
		return
	}

	record[req.RemoteAddr] = record[req.RemoteAddr] + 1

	go func() {
		time.Sleep(time.Duration(config.Duration) * time.Millisecond)
		record[req.RemoteAddr] = record[req.RemoteAddr] - 1
	}()

	req.URL.Path = config.TargetServer + req.URL.Path

	createdReq, err := http.NewRequest(req.Method, req.URL.Path, req.Body)
	if err != nil {
		res.WriteHeader(400)
		fmt.Fprint(res, err)
		return
	}

	client := http.Client{}
	response, err := client.Do(createdReq)
	if err != nil {
		res.WriteHeader(400)
		fmt.Fprint(res, err)
		return
	}

	ProxyResponse(*response, res)
}

var config Config
var record map[string]uint64

func main() {
	// read config file
	data, err := ioutil.ReadFile("./config.json")
	if err != nil {
		panic("Config file is not found")
	}

	json.Unmarshal(data, &config)

	record = make(map[string]uint64)

	http.HandleFunc("/", ProxyRequest)

	http.ListenAndServe(":8080", nil)
}
