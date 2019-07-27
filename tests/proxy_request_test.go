package tests

import (
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
