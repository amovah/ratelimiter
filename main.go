package main

import (
  // "bytes"
	"fmt"
  // "io"
  "encoding/json"
  "net/http"
  "strings"
  // "bufio"
  "io/ioutil"
)

type Config struct {
	RequestPerDuration uint
	Duration uint
	TargetServer string
}

func proxy(res http.ResponseWriter, req *http.Request) {
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
}

var config Config

func main() {
	// read config file
	data, err := ioutil.ReadFile("./config.json")
	if err != nil {
		panic("Config file is not found")
	}

	json.Unmarshal(data, &config)

  http.HandleFunc("/", proxy)

  http.ListenAndServe(":8080", nil)
}
