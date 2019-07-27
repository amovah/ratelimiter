package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

func proxyResponse(origin http.Response, target http.ResponseWriter) {
	defer origin.Body.Close()

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

func proxyRequest(res http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	ip := strings.Split(req.RemoteAddr, ":")[0]

	if record[ip] > maxRatePerIP || totalRequest > totalMaxRate {
		http.Error(res, "You reach your limit", 429)
		return
	}

	mutex.Lock()
	record[ip] = record[ip] + 1
	totalRequest = totalRequest + 1
	mutex.Unlock()

	go func() {
		time.Sleep(time.Second)
		mutex.Lock()
		record[ip] = record[ip] - 1
		totalRequest = totalRequest - 1
		mutex.Unlock()
	}()

	createdReq, err := http.NewRequest(req.Method, targetServer+req.URL.Path, req.Body)
	if err != nil {
		res.WriteHeader(400)
		fmt.Fprint(res, err)
		return
	}

	for key, value := range req.Header {
		createdReq.Header.Set(key, strings.Join(value, ", "))
	}

	client := http.Client{}
	response, err := client.Do(createdReq)
	if err != nil {
		res.WriteHeader(400)
		fmt.Fprint(res, err)
		return
	}

	proxyResponse(*response, res)
}

var record map[string]uint64
var totalRequest uint64
var maxRatePerIP uint64
var totalMaxRate uint64
var targetServer string
var mutex sync.Mutex

func main() {
	record = make(map[string]uint64)
	totalRequest = 0
	mutex = sync.Mutex{}
	var err error

	maxRatePerIP, err = strconv.ParseUint(os.Getenv("MAX_RATE_PER_IP"), 10, 64)
	if err != nil {
		maxRatePerIP = 100
	}

	totalMaxRate, err = strconv.ParseUint(os.Getenv("TOTAL_MAX_RATE"), 10, 64)
	if err != nil {
		totalMaxRate = 10000
	}

	targetServer = os.Getenv("TARGET_SERVER")
	if targetServer == "" {
		fmt.Println("TARGET_SERVER environment variable cannot be empty")
		return
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/", proxyRequest)

	http.ListenAndServe(":"+port, nil)
}
