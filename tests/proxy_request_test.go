package tests

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProxyGet(t *testing.T) {
	go Server()

	config := LoadConfig()
	client := http.Client{}
	request, _ := http.NewRequest("GET", config.ProxyServerPath+"/get", nil)

	res, err := client.Do(request)
	if err != nil {
		t.Fatal(err)
	}

	defer res.Body.Close()

	assert.Equal(t, 200, res.StatusCode, "status code is not proxied")

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "200 successful get", string(body), "body is not proxied")
}

func TestProxyHeader(t *testing.T) {
	go Server()

	config := LoadConfig()
	client := http.Client{}
	request, err := http.NewRequest("GET", config.ProxyServerPath+"/header", nil)
	if err != nil {
		t.Fatal(err)
	}

	request.Header.Add("x-test", "sample text")

	res, err := client.Do(request)
	if err != nil {
		t.Fatal(err)
	}

	defer res.Body.Close()

	assert.Equal(t, 200, res.StatusCode, "status code is not proxied")
	assert.Equal(t, "sample text", res.Header.Get("x-test"), "header is not proxied")

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "header x-test", string(body), "body is not proxied")
}

func TestProxyStatusCode(t *testing.T) {
	go Server()

	config := LoadConfig()
	client := http.Client{}
	request, err := http.NewRequest("GET", config.ProxyServerPath+"/get-failure", nil)
	if err != nil {
		t.Fatal(err)
	}

	res, err := client.Do(request)
	if err != nil {
		t.Fatal(err)
	}

	defer res.Body.Close()

	assert.Equal(t, 400, res.StatusCode, "status code is not proxied")

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "400 very bad request", string(body), "body is not proxied")
}

func TestProxyBody(t *testing.T) {
	go Server()

	data, err := json.Marshal(PostedData{
		Username: "Ali",
		Password: "123",
		Email:    "ali_movahedi@aol.com",
	})
	if err != nil {
		t.Fatal(err)
	}

	config := LoadConfig()
	client := http.Client{}
	request, err := http.NewRequest("POST", config.ProxyServerPath+"/post", bytes.NewReader(data))
	if err != nil {
		t.Fatal(err)
	}

	res, err := client.Do(request)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	assert.Equal(t, 200, res.StatusCode, "status code is not proxied")

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}

	parsed := PostedData{}
	err = json.Unmarshal(body, &parsed)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "Alimodified", parsed.Username, "body is not proxied")
	assert.Equal(t, "123", parsed.Password, "body is not proxied")
	assert.Equal(t, "ali_movahedi@aol.com", parsed.Email, "body is not proxied")
}

func TestProxyPut(t *testing.T) {
	go Server()

	config := LoadConfig()
	client := http.Client{}
	request, _ := http.NewRequest("PUT", config.ProxyServerPath+"/put", nil)

	res, err := client.Do(request)
	if err != nil {
		t.Fatal(err)
	}

	defer res.Body.Close()

	assert.Equal(t, 403, res.StatusCode, "status code is not proxied")

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "403 Forbidden", string(body), "body is not proxied")
}
