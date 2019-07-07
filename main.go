package main

import (
	"fmt"
  "encoding/json"
  "net/http"
  "strings"
  "io/ioutil"
	"time"
)

type Config struct {
	RequestPerDuration uint
	Duration uint
	TargetServer string
}

func proxy(res http.ResponseWriter, req *http.Request) {
	if requestCount <= config.RequestPerDuration {
		requestCount = requestCount + 1

		go func() {
			time.Sleep(time.Duration(config.Duration) * time.Millisecond)
			requestCount = requestCount - 1
		}()

		req.URL.Path = config.TargetServer + req.URL.Path

		if req.Method == "GET" {
	    response, err := http.Get(req.URL.Path)

	    if err != nil {
	      fmt.Fprint(res, "ERROR")
	    } else {
	      body, err := ioutil.ReadAll(response.Body)
	      if err != nil {
	        fmt.Fprint(res, "ERROR")
	      } else {
	        for key, value := range response.Header {
	          res.Header().Set(key, strings.Join(value, ", "))
	        }

	        res.Write(body)
	      }
	    }
	  }
	} else {
		http.Error(res, "You reach your limit", 429)
	}
}

var config Config
var requestCount uint

func main() {
	// read config file
	data, err := ioutil.ReadFile("./config.json")
	if err != nil {
		panic("Config file is not found")
	}

	json.Unmarshal(data, &config)

	requestCount = 0

  http.HandleFunc("/", proxy)

  http.ListenAndServe(":8080", nil)
}
