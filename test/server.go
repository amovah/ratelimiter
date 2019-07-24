package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
)

func getTest(res http.ResponseWriter, req *http.Request) {
	if req.Method == "GET" {
		res.Write([]byte("200 successful get"))
	}
}

func getTestFailure(res http.ResponseWriter, req *http.Request) {
	if req.Method == "GET" {
		res.WriteHeader(400)
		res.Write([]byte("400 bad request"))
	}
}

type PostedData struct {
	Username string
	Password string
	Email    string
}

func postTest(res http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" {
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			res.WriteHeader(400)
			res.Write([]byte("400 bad request"))
			return
		}

		var postedData PostedData
		json.Unmarshal(body, &postedData)
		postedData.Username = postedData.Username + "modified"

		jsonied, err := json.Marshal(postedData)
		if err != nil {
			res.WriteHeader(400)
			res.Write([]byte("400 bad request"))
			return
		}

		res.Write(jsonied)
	}
}

func putTest(res http.ResponseWriter, req *http.Request) {
	if req.Method == "PUT" {
		res.WriteHeader(403)
		res.Write([]byte("403 Forbidden"))
	}
}

func deleteTest(res http.ResponseWriter, req *http.Request) {
	if req.Method == "DELETE" {
		res.WriteHeader(404)
		res.Write([]byte("404 Model Not Found"))
	}
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	http.HandleFunc("/get", getTest)
	http.HandleFunc("/get-failure", getTestFailure)
	http.HandleFunc("/post", postTest)
	http.HandleFunc("/put", putTest)
	http.HandleFunc("/delete", deleteTest)

	http.ListenAndServe(":"+port, nil)
}
