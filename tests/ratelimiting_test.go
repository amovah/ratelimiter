package tests

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRateLimiting(t *testing.T) {
	go Server()

	config := LoadConfig()
	client := http.Client{}
	request, _ := http.NewRequest("GET", config.ProxyServerPath+"/get", nil)
	maxRate := int(config.MaxRatePerIP) + 1

	for i := 0; i < maxRate; i++ {
		res, err := client.Do(request)
		if err != nil {
			t.Fatal(err)
		}
		res.Body.Close()
	}

	res, err := client.Do(request)

	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 429, res.StatusCode)
}
