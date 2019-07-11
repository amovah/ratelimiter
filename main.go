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

func ProxyResponse(origin http.Response, target http.ResponseWriter) {
	body, err := ioutil.ReadAll(origin.Body)
	if err != nil {
		fmt.Fprint(target, "ERROR")
	} else {
		for key, value := range origin.Header {
			target.Header().Set(key, strings.Join(value, ", "))
		}

		target.Write(body)
	}
}

func proxy(res http.ResponseWriter, req *http.Request) {
	if record[req.RemoteAddr] <= uint64(config.RequestPerDuration) {
		record[req.RemoteAddr] = record[req.RemoteAddr] + 1

		go func() {
			time.Sleep(time.Duration(config.Duration) * time.Millisecond)
			record[req.RemoteAddr] = record[req.RemoteAddr] - 1
		}()

		req.URL.Path = config.TargetServer + req.URL.Path

		if req.Method == "GET" {
	    response, err := http.Get(req.URL.Path)

	    if err != nil {
				res.WriteHeader(400)
				fmt.Fprint(res, err)
				return
	    }

			ProxyResponse(*response, res)
	  }

		if req.Method == "POST" {
			response, err := http.Post(req.URL.Path, req.Header.Get("Content-Type"), req.Body)

	    if err != nil {
				res.WriteHeader(400)
				fmt.Fprint(res, err)
				return
	    }

			ProxyResponse(*response, res)
		} else {
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
	} else {
		http.Error(res, "You reach your limit", 429)
	}
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

  http.HandleFunc("/", proxy)

  http.ListenAndServe(":8080", nil)
}
