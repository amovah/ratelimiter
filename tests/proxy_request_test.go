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
