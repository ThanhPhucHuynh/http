package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/ThanhPhucHuynh/http/hystrix"
)

type linearBackoff struct {
	backoffInterval int
}

func (lb *linearBackoff) Next(retry int) time.Duration {
	if retry <= 0 {
		return 0 * time.Millisecond
	}
	return time.Duration(retry*lb.backoffInterval) * time.Millisecond
}

type demoRoundTripper struct{}

func (t *demoRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	r := req.Clone(req.Context())
	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("Accept", "application/json")
	r.Header.Set("Cookie", "TOKEN=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE2OTExMjU3MzUsImV4cCI6MTY5MTE2ODkzNX0.MgPXK3K6j1wPiHV_8HU9q3NBgHbSKyJZgS6suSGe1lc")

	return http.DefaultTransport.RoundTrip(r)
}
func main() {

	// Create a new hystrix-wrapped HTTP client with the command name, along with other required options
	client := hystrix.NewClient(
		hystrix.WithHTTPTimeout(10*time.Millisecond),
		hystrix.WithCommandName("google_get_request"),
		hystrix.WithHystrixTimeout(1000*time.Millisecond),
		hystrix.WithMaxConcurrentRequests(30),
		hystrix.WithErrorPercentThreshold(20),
		hystrix.WithFallbackFunc(func(err error) error {
			// do something in the fallback function
			fmt.Println("???", err)
			return nil
		}),
		hystrix.WithHTTPClient(&http.Client{
			Transport: &demoRoundTripper{},
		}),
	)
	// req, _ := http.NewRequest(http.MethodGet, "http://localhost:9090/swagger/auth/index.html", nil)
	h := http.Header{}

	res, err := client.Get("http://localhost:9090/swagger/auth/index.html", h)
	if err != nil {
		panic(err)
	}
	body, err := ioutil.ReadAll(res.Body)
	fmt.Println(string(body))
}
