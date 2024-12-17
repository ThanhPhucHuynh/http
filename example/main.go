package examples

import (
	"fmt"
	"io"
	"net/http"
	"time"

	avhttp "github.com/thanhphuchuynh/http"
	"github.com/thanhphuchuynh/http/httpclient"
	"github.com/thanhphuchuynh/http/hystrix"

	"github.com/pkg/errors"
)

const (
	baseURL = "http://localhost:9090"
)

func httpClientUsage() error {
	timeout := 100 * time.Millisecond

	httpClient := httpclient.NewClient(
		httpclient.WithHTTPTimeout(timeout),
		httpclient.WithRetryCount(2),
		httpclient.WithRetrier(avhttp.NewRetrier(avhttp.NewConstantBackoff(10*time.Millisecond, 50*time.Millisecond))),
	)
	headers := http.Header{}
	headers.Set("Content-Type", "application/json")

	response, err := httpClient.Get(baseURL, headers)
	if err != nil {
		return errors.Wrap(err, "failed to make a request to server")
	}

	defer response.Body.Close()

	respBody, err := io.ReadAll(response.Body)
	if err != nil {
		return errors.Wrap(err, "failed to read response body")
	}

	fmt.Printf("Response: %s", string(respBody))
	return nil
}

func hystrixClientUsage() error {
	timeout := 100 * time.Millisecond
	hystrixClient := hystrix.NewClient(
		hystrix.WithHTTPTimeout(timeout),
		hystrix.WithCommandName("MyCommand"),
		hystrix.WithHystrixTimeout(1100*time.Millisecond),
		hystrix.WithMaxConcurrentRequests(100),
		hystrix.WithErrorPercentThreshold(25),
		hystrix.WithSleepWindow(10),
		hystrix.WithRequestVolumeThreshold(10),
		hystrix.WithStatsDCollector("localhost:8125", "myapp.hystrix"),
	)
	headers := http.Header{}
	response, err := hystrixClient.Get(baseURL, headers)
	if err != nil {
		return errors.Wrap(err, "failed to make a request to server")
	}

	defer response.Body.Close()

	respBody, err := io.ReadAll(response.Body)
	if err != nil {
		return errors.Wrap(err, "failed to read response body")
	}

	fmt.Printf("Response: %s", string(respBody))
	return nil
}

type myHTTPClient struct {
	client http.Client
}

func (c *myHTTPClient) Do(request *http.Request) (*http.Response, error) {
	request.SetBasicAuth("username", "passwd")
	return c.client.Do(request)
}

func customHTTPClientUsage() error {
	httpClient := httpclient.NewClient(
		httpclient.WithHTTPTimeout(0*time.Millisecond),
		httpclient.WithHTTPClient(&myHTTPClient{
			// replace with custom HTTP client
			client: http.Client{Timeout: 25 * time.Millisecond},
		}),
		httpclient.WithRetryCount(2),
		httpclient.WithRetrier(avhttp.NewRetrier(avhttp.NewConstantBackoff(10*time.Millisecond, 50*time.Millisecond))),
	)

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")

	response, err := httpClient.Get(baseURL, headers)
	if err != nil {
		return errors.Wrap(err, "failed to make a request to server")
	}

	defer response.Body.Close()

	respBody, err := io.ReadAll(response.Body)
	if err != nil {
		return errors.Wrap(err, "failed to read response body")
	}

	fmt.Printf("Response: %s", string(respBody))
	return nil
}

func customHystrixClientUsage() error {
	timeout := 0 * time.Millisecond

	hystrixClient := hystrix.NewClient(
		hystrix.WithHTTPTimeout(timeout),
		hystrix.WithCommandName("MyCommand"),
		hystrix.WithHystrixTimeout(1100*time.Millisecond),
		hystrix.WithMaxConcurrentRequests(100),
		hystrix.WithErrorPercentThreshold(25),
		hystrix.WithSleepWindow(10),
		hystrix.WithRequestVolumeThreshold(10),
		hystrix.WithHTTPClient(&myHTTPClient{
			// replace with custom HTTP client
			client: http.Client{Timeout: 25 * time.Millisecond},
		}),
	)

	headers := http.Header{}
	response, err := hystrixClient.Get(baseURL, headers)
	if err != nil {
		return errors.Wrap(err, "failed to make a request to server")
	}

	defer response.Body.Close()

	respBody, err := io.ReadAll(response.Body)
	if err != nil {
		return errors.Wrap(err, "failed to read response body")
	}

	fmt.Printf("Response: %s", string(respBody))
	return nil
}
