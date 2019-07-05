package main

import (
  // "bytes"
	"fmt"
  // "io"
  // "encoding/json"
  "net/http"
  "strings"
  // "bufio"
  "io/ioutil"
)

// type Shit struct {
//   Name string
// }

func proxy(res http.ResponseWriter, req *http.Request) {
  req.URL.Path = "http://localhost:8010" + req.URL.Path

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

func main() {
  http.HandleFunc("/", proxy)

  http.ListenAndServe(":8080", nil)
}
