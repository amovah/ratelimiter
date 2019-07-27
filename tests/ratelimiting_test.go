package tests

import (
	"net/http"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRateLimiting(t *testing.T) {
	go Server()

	var wg sync.WaitGroup
	config := LoadConfig()
	client := http.Client{}
	request, _ := http.NewRequest("GET", config.ProxyServerPath+"/get", nil)

	for i := 0; i < int(config.MaxRatePerIP)+1; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			res, err := client.Do(request)

			if err != nil {
				t.Fatal("something with proxy server went wrong!", err)
			}

			defer res.Body.Close()
		}()
	}

	wg.Wait()

	res, err := client.Do(request)

	if err != nil {
		t.Fatal("something with proxy server went wrong!", err)
	}

	assert.Equal(t, 429, res.StatusCode)
}
