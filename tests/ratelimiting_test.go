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

	for i := 0; i < int(config.MaxRatePerIP)+100; i++ {
		wg.Add(1)
		func() {
			res, err := client.Do(request)

			if err != nil {
				t.FailNow()
			}

			defer res.Body.Close()

			wg.Done()
		}()
	}

	wg.Wait()

	res, err := client.Do(request)

	if err != nil {
		t.FailNow()
	}

	assert.Equal(t, res.StatusCode, 429)
}
