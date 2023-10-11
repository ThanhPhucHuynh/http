package hystrix

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	avhttp "github.com/ThanhPhucHuynh/http"
	"github.com/ThanhPhucHuynh/http/plugins"
	hystrixplugins "github.com/afex/hystrix-go/plugins"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type myHTTPClient struct {
	client http.Client
}

func (c *myHTTPClient) Do(request *http.Request) (*http.Response, error) {
	request.Header.Set("foo", "bar")
	return c.client.Do(request)
}

type errReader int

func (errReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("test error")
}

type MockStatsDCollector struct {
	ReturnError bool
}

func (m *MockStatsDCollector) InitializeStatsdCollector(config *hystrixplugins.StatsdCollectorConfig) (*hystrixplugins.StatsdCollectorClient, error) {
	if m.ReturnError {
		return nil, errors.New("mock error")
	}
	return nil, nil
}

func TestHystrixHTTPClientDoSuccess(t *testing.T) {
	client := NewClient(
		WithHTTPTimeout(50*time.Millisecond),
		WithCommandName("some_command_name"),
		WithHystrixTimeout(10*time.Millisecond),
		WithMaxConcurrentRequests(100),
		WithErrorPercentThreshold(10),
		WithSleepWindow(100),
		WithRequestVolumeThreshold(20),
	)

	dummyHandler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.Equal(t, "en", r.Header.Get("Accept-Language"))

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{ "response": "ok" }`))
	}

	server := httptest.NewServer(http.HandlerFunc(dummyHandler))
	defer server.Close()

	req, err := http.NewRequest(http.MethodGet, server.URL, nil)
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept-Language", "en")

	response, err := client.Do(req)
	require.NoError(t, err, "should not have failed to make a GET request")

	assert.Equal(t, http.StatusOK, response.StatusCode)
	body, err := io.ReadAll(response.Body)
	require.NoError(t, err)
	assert.Equal(t, "{ \"response\": \"ok\" }", string(body))
}

func TestHystrixHTTPClientGetSuccess(t *testing.T) {
	client := NewClient(
		WithHTTPTimeout(10*time.Millisecond),
		WithCommandName("some_command_name"),
		WithHystrixTimeout(10*time.Millisecond),
		WithMaxConcurrentRequests(100),
		WithErrorPercentThreshold(10),
		WithSleepWindow(100),
		WithRequestVolumeThreshold(10),
	)

	dummyHandler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.Equal(t, "en", r.Header.Get("Accept-Language"))

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{ "response": "ok" }`))
	}

	server := httptest.NewServer(http.HandlerFunc(dummyHandler))
	defer server.Close()

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")
	headers.Set("Accept-Language", "en")

	response, err := client.Get(server.URL, headers)
	require.NoError(t, err, "should not have failed to make a GET request")

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, "{ \"response\": \"ok\" }", respBody(t, response))
}

func TestHystrixHTTPClientPostSuccess(t *testing.T) {
	client := NewClient(
		WithHTTPTimeout(10*time.Millisecond),
		WithCommandName("some_command_name"),
		WithHystrixTimeout(10*time.Millisecond),
		WithMaxConcurrentRequests(100),
		WithErrorPercentThreshold(10),
		WithSleepWindow(100),
		WithRequestVolumeThreshold(10),
	)

	requestBodyString := `{ "name": "avhttp" }`

	dummyHandler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.Equal(t, "en", r.Header.Get("Accept-Language"))

		rBody, err := io.ReadAll(r.Body)
		require.NoError(t, err, "should not have failed to extract request body")

		assert.Equal(t, requestBodyString, string(rBody))

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{ "response": "ok" }`))
	}

	server := httptest.NewServer(http.HandlerFunc(dummyHandler))
	defer server.Close()

	requestBody := bytes.NewReader([]byte(requestBodyString))

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")
	headers.Set("Accept-Language", "en")

	response, err := client.Post(server.URL, requestBody, headers)
	require.NoError(t, err, "should not have failed to make a POST request")

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, "{ \"response\": \"ok\" }", respBody(t, response))
}

func TestHystrixHTTPClientDeleteSuccess(t *testing.T) {
	client := NewClient(
		WithHTTPTimeout(10*time.Millisecond),
		WithCommandName("some_command_name"),
		WithHystrixTimeout(10*time.Millisecond),
		WithMaxConcurrentRequests(100),
		WithErrorPercentThreshold(10),
		WithSleepWindow(100),
		WithRequestVolumeThreshold(10),
	)

	dummyHandler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.Equal(t, "en", r.Header.Get("Accept-Language"))

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{ "response": "ok" }`))
	}

	server := httptest.NewServer(http.HandlerFunc(dummyHandler))
	defer server.Close()

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")
	headers.Set("Accept-Language", "en")

	response, err := client.Delete(server.URL, headers)
	require.NoError(t, err, "should not have failed to make a DELETE request")

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, "{ \"response\": \"ok\" }", respBody(t, response))
}

func TestHystrixHTTPClientPutSuccess(t *testing.T) {
	client := NewClient(
		WithHTTPTimeout(10*time.Millisecond),
		WithCommandName("some_command_name"),
		WithHystrixTimeout(10*time.Millisecond),
		WithMaxConcurrentRequests(100),
		WithErrorPercentThreshold(10),
		WithSleepWindow(100),
		WithRequestVolumeThreshold(10),
	)

	requestBodyString := `{ "name": "avhttp" }`

	dummyHandler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPut, r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.Equal(t, "en", r.Header.Get("Accept-Language"))

		rBody, err := io.ReadAll(r.Body)
		require.NoError(t, err, "should not have failed to extract request body")

		assert.Equal(t, requestBodyString, string(rBody))

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{ "response": "ok" }`))
	}

	server := httptest.NewServer(http.HandlerFunc(dummyHandler))
	defer server.Close()

	requestBody := bytes.NewReader([]byte(requestBodyString))

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")
	headers.Set("Accept-Language", "en")

	response, err := client.Put(server.URL, requestBody, headers)
	require.NoError(t, err, "should not have failed to make a PUT request")

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, "{ \"response\": \"ok\" }", respBody(t, response))
}

func TestHystrixHTTPClientPatchSuccess(t *testing.T) {
	client := NewClient(
		WithHTTPTimeout(10*time.Millisecond),
		WithCommandName("some_command_name"),
		WithHystrixTimeout(10*time.Millisecond),
		WithMaxConcurrentRequests(100),
		WithErrorPercentThreshold(10),
		WithSleepWindow(100),
		WithRequestVolumeThreshold(10),
	)

	requestBodyString := `{ "name": "avhttp" }`

	dummyHandler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPatch, r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.Equal(t, "en", r.Header.Get("Accept-Language"))

		rBody, err := io.ReadAll(r.Body)
		require.NoError(t, err, "should not have failed to extract request body")

		assert.Equal(t, requestBodyString, string(rBody))

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{ "response": "ok" }`))
	}

	server := httptest.NewServer(http.HandlerFunc(dummyHandler))
	defer server.Close()

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")
	headers.Set("Accept-Language", "en")

	requestBody := bytes.NewReader([]byte(requestBodyString))

	response, err := client.Patch(server.URL, requestBody, headers)
	require.NoError(t, err, "should not have failed to make a PATCH request")

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, "{ \"response\": \"ok\" }", respBody(t, response))
}

func TestHystrixHTTPClientRetriesGetOnFailure(t *testing.T) {
	backoffInterval := 1 * time.Millisecond
	maximumJitterInterval := 1 * time.Millisecond

	client := NewClient(
		WithHTTPTimeout(10*time.Millisecond),
		WithCommandName("some_command_name"),
		WithHystrixTimeout(10*time.Millisecond),
		WithMaxConcurrentRequests(100),
		WithErrorPercentThreshold(10),
		WithSleepWindow(100),
		WithRequestVolumeThreshold(10),
		WithRetryCount(3),
		WithRetrier(avhttp.NewRetrier(avhttp.NewConstantBackoff(backoffInterval, maximumJitterInterval))),
	)

	response, err := client.Get("url_doesnt_exist", http.Header{})

	assert.Contains(t, err.Error(), "unsupported protocol scheme")
	assert.Nil(t, response)
}

func TestHystrixHTTPClientRetriesGetOnFailure5xx(t *testing.T) {
	count := 0
	backoffInterval := 1 * time.Millisecond
	maximumJitterInterval := 1 * time.Millisecond
	client := NewClient(
		WithHTTPTimeout(10*time.Millisecond),
		WithCommandName("some_command_name_5xx"),
		WithHystrixTimeout(10*time.Millisecond),
		WithMaxConcurrentRequests(100),
		WithErrorPercentThreshold(10),
		WithSleepWindow(100),
		WithRequestVolumeThreshold(10),
		WithRetryCount(3),
		WithRetrier(avhttp.NewRetrier(avhttp.NewConstantBackoff(backoffInterval, maximumJitterInterval))),
	)

	dummyHandler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{ "response": "something went wrong" }`))
		count = count + 1
	}

	server := httptest.NewServer(http.HandlerFunc(dummyHandler))
	defer server.Close()

	response, err := client.Get(server.URL, http.Header{})
	require.NoError(t, err)

	assert.Equal(t, 4, count)

	assert.Equal(t, http.StatusInternalServerError, response.StatusCode)
	assert.Equal(t, "{ \"response\": \"something went wrong\" }", respBody(t, response))
}

func BenchmarkHystrixHTTPClientRetriesGetOnFailure(b *testing.B) {
	backoffInterval := 1 * time.Millisecond
	maximumJitterInterval := 1 * time.Millisecond

	client := NewClient(
		WithHTTPTimeout(10*time.Millisecond),
		WithCommandName("some_command_name"),
		WithHystrixTimeout(10*time.Millisecond),
		WithMaxConcurrentRequests(100),
		WithErrorPercentThreshold(10),
		WithSleepWindow(100),
		WithRequestVolumeThreshold(10),
		WithRetryCount(3),
		WithRetrier(avhttp.NewRetrier(avhttp.NewConstantBackoff(backoffInterval, maximumJitterInterval))),
	)

	dummyHandler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{ "response": "something went wrong" }`))
	}

	server := httptest.NewServer(http.HandlerFunc(dummyHandler))
	defer server.Close()

	for i := 0; i < b.N; i++ {
		_, _ = client.Get(server.URL, http.Header{})
	}
}

func TestHystrixHTTPClientRetriesPostOnFailure(t *testing.T) {
	count := 0
	backoffInterval := 1 * time.Millisecond
	maximumJitterInterval := 1 * time.Millisecond

	client := NewClient(
		WithHTTPTimeout(50*time.Millisecond),
		WithCommandName("some_command_name"),
		WithHystrixTimeout(10*time.Millisecond),
		WithMaxConcurrentRequests(100),
		WithErrorPercentThreshold(10),
		WithSleepWindow(100),
		WithRequestVolumeThreshold(20),
		WithRetryCount(3),
		WithRetrier(avhttp.NewRetrier(avhttp.NewConstantBackoff(backoffInterval, maximumJitterInterval))),
	)

	dummyHandler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{ "response": "something went wrong" }`))
		count = count + 1
	}

	server := httptest.NewServer(http.HandlerFunc(dummyHandler))
	defer server.Close()

	response, err := client.Post(server.URL, strings.NewReader("a=1&b=2"), http.Header{})
	require.NoError(t, err)

	assert.Equal(t, 4, count)
	assert.Equal(t, http.StatusInternalServerError, response.StatusCode)
	assert.JSONEq(t, `{ "response": "something went wrong" }`, respBody(t, response))
}

func BenchmarkHystrixHTTPClientRetriesPostOnFailure(b *testing.B) {
	backoffInterval := 1 * time.Millisecond
	maximumJitterInterval := 1 * time.Millisecond

	client := NewClient(
		WithHTTPTimeout(10*time.Millisecond),
		WithCommandName("some_command_name"),
		WithHystrixTimeout(10*time.Millisecond),
		WithMaxConcurrentRequests(100),
		WithErrorPercentThreshold(10),
		WithSleepWindow(100),
		WithRequestVolumeThreshold(10),
		WithRetryCount(3),
		WithRetrier(avhttp.NewRetrier(avhttp.NewConstantBackoff(backoffInterval, maximumJitterInterval))),
	)

	dummyHandler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{ "response": "something went wrong" }`))
	}

	server := httptest.NewServer(http.HandlerFunc(dummyHandler))
	defer server.Close()

	for i := 0; i < b.N; i++ {
		_, _ = client.Post(server.URL, strings.NewReader("a=1&b=2"), http.Header{})
	}
}

func TestHystrixHTTPClientReturnsFallbackFailureWithoutFallBackFunction(t *testing.T) {
	client := NewClient(
		WithHTTPTimeout(10*time.Millisecond),
		WithCommandName("some_command_name"),
		WithHystrixTimeout(10*time.Millisecond),
		WithMaxConcurrentRequests(100),
		WithErrorPercentThreshold(10),
		WithSleepWindow(100),
		WithRequestVolumeThreshold(10),
	)

	_, err := client.Get("http://foobar.example", http.Header{})
	assert.Equal(t, err.Error(), "hystrix: circuit open")
}

func TestHystrixHTTPClientReturnsFallbackFailureWithAFallBackFunctionWhichReturnAnError(t *testing.T) {
	client := NewClient(
		WithHTTPTimeout(10*time.Millisecond),
		WithCommandName("some_command_name"),
		WithHystrixTimeout(10*time.Millisecond),
		WithMaxConcurrentRequests(100),
		WithErrorPercentThreshold(10),
		WithSleepWindow(100),
		WithRequestVolumeThreshold(10),
		WithFallbackFunc(func(err error) error {
			// do something in the fallback function
			return err
		}),
	)

	_, err := client.Get("http://foobar.example", http.Header{})
	require.Error(t, err, "should have failed")

	assert.True(t, strings.Contains(err.Error(), "fallback failed"))
}

func TestFallBackFunctionIsCalledWithHystrixHTTPClient(t *testing.T) {
	called := false

	client := NewClient(
		WithHTTPTimeout(10*time.Millisecond),
		WithCommandName("some_command_name"),
		WithHystrixTimeout(10*time.Millisecond),
		WithMaxConcurrentRequests(100),
		WithErrorPercentThreshold(10),
		WithSleepWindow(100),
		WithRequestVolumeThreshold(10),
		WithFallbackFunc(func(err error) error {
			called = true
			return err
		}),
	)

	_, err := client.Get("http://foobar.example", http.Header{})
	require.Error(t, err, "should have failed")

	assert.True(t, called)
}

func TestHystrixHTTPClientReturnsFallbackFailureWithAFallBackFunctionWhichReturnsNil(t *testing.T) {
	client := NewClient(
		WithHTTPTimeout(10*time.Millisecond),
		WithCommandName("some_command_name"),
		WithHystrixTimeout(10*time.Millisecond),
		WithMaxConcurrentRequests(100),
		WithErrorPercentThreshold(10),
		WithSleepWindow(100),
		WithRequestVolumeThreshold(10),
		WithFallbackFunc(func(err error) error {
			// do something in the fallback function
			return nil
		}),
	)

	_, err := client.Get("http://foobar.example", http.Header{})
	assert.Nil(t, err)
}

func TestCustomHystrixHTTPClientDoSuccess(t *testing.T) {
	client := NewClient(
		WithHTTPTimeout(10*time.Millisecond),
		WithCommandName("some_new_command_name"),
		WithHystrixTimeout(10*time.Millisecond),
		WithMaxConcurrentRequests(100),
		WithErrorPercentThreshold(10),
		WithSleepWindow(100),
		WithRequestVolumeThreshold(10),
		WithHTTPClient(&myHTTPClient{
			client: http.Client{Timeout: 25 * time.Millisecond}}),
	)

	dummyHandler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Header.Get("foo"), "bar")
		assert.NotEqual(t, r.Header.Get("foo"), "baz")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{ "response": "ok" }`))
	}

	server := httptest.NewServer(http.HandlerFunc(dummyHandler))
	defer server.Close()

	req, err := http.NewRequest(http.MethodGet, server.URL, nil)
	require.NoError(t, err)
	response, err := client.Do(req)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	body, err := io.ReadAll(response.Body)
	require.NoError(t, err)
	assert.Equal(t, "{ \"response\": \"ok\" }", string(body))
}

func respBody(t *testing.T, response *http.Response) string {
	if response.Body != nil {
		defer func() {
			_ = response.Body.Close()
		}()
	}

	respBody, err := io.ReadAll(response.Body)
	require.NoError(t, err, "should not have failed to read response body")

	return string(respBody)
}

func TestDurationToInt(t *testing.T) {
	t.Run("1sec should return 1 when unit is second", func(t *testing.T) {
		timeout := 1 * time.Second
		timeoutInSec := durationToInt(timeout, time.Second)

		assert.Equal(t, 1, timeoutInSec)
	})

	t.Run("30sec should return 30000 when unit is millisecond", func(t *testing.T) {
		timeout := 30 * time.Second
		timeoutInMs := durationToInt(timeout, time.Millisecond)

		assert.Equal(t, 30000, timeoutInMs)
	})

	t.Run("max should return int(maxUint >> 1) when unit is millisecond", func(t *testing.T) {
		timeout := time.Duration(1<<63 - 1)
		timeoutInMs := durationToInt(timeout, time.Nanosecond)

		assert.Equal(t, int(maxUint>>1), timeoutInMs)
	})
}

func TestHTTPClientGetReturnsErrorDoReader(t *testing.T) {
	client := NewClient(WithHTTPTimeout(10 * time.Millisecond))

	response, err := client.Do(httptest.NewRequest("POST", "http://127.0.0.1:40014", errReader(0)))
	assert.Contains(t, err.Error(), "test error")
	assert.Nil(t, response)
}

func TestHTTPClientGetReturnsErrorNewRequest(t *testing.T) {
	client := NewClient(WithHTTPTimeout(10 * time.Millisecond))

	_, err := client.Get("http://127.0.0.1:40014endpoint", http.Header{})
	assert.Contains(t, err.Error(), "GET - request creation failed")
}

func TestHTTPClientPostReturnsErrorNewRequest(t *testing.T) {
	client := NewClient(WithHTTPTimeout(10 * time.Millisecond))

	_, err := client.Post("http://127.0.0.1:40014endpoint", nil, http.Header{})
	assert.Contains(t, err.Error(), "POST - request creation failed")
}

func TestHTTPClientPatchReturnsErrorNewRequest(t *testing.T) {
	client := NewClient(WithHTTPTimeout(10 * time.Millisecond))

	_, err := client.Patch("http://127.0.0.1:40014endpoint", nil, http.Header{})
	assert.Contains(t, err.Error(), "PATCH - request creation failed")
}

func TestHTTPClientDeleteReturnsErrorNewRequest(t *testing.T) {
	client := NewClient(WithHTTPTimeout(10 * time.Millisecond))

	_, err := client.Delete("http://127.0.0.1:40014endpoint", http.Header{})
	assert.Contains(t, err.Error(), "DELETE - request creation failed")
}

func TestHTTPClientPutReturnsErrorNewRequest(t *testing.T) {
	client := NewClient(WithHTTPTimeout(10 * time.Millisecond))

	_, err := client.Put("http://127.0.0.1:40014endpoint", nil, http.Header{})
	assert.Contains(t, err.Error(), "PUT - request creation failed")
}

func TestHTTPClientAddPlugin(t *testing.T) {
	var buf bytes.Buffer

	client := NewClient(WithHTTPTimeout(10 * time.Millisecond))
	requestLogger := plugins.NewRequestLogger(nil, &buf)
	client.AddPlugin(requestLogger)

	_, err := client.Put("http://127.0.0.1:40014endpoint", nil, http.Header{})
	assert.Contains(t, err.Error(), "PUT - request creation failed")
}

func TestHystrixHTTPClientFailedWhenCreate(t *testing.T) {
	mockCollector := &MockStatsDCollector{ReturnError: true}
	// Replace the actual implementation with the mock
	oldInitializeStatsdCollector := hystrixplugins.InitializeStatsdCollector
	initializeStatsdCollector = func(config *hystrixplugins.StatsdCollectorConfig) (*hystrixplugins.StatsdCollectorClient, error) {
		return mockCollector.InitializeStatsdCollector(config)
	}
	defer func() {
		initializeStatsdCollector = oldInitializeStatsdCollector
	}()

	shouldPanic(t, func() {
		_ = NewClient(
			WithHTTPTimeout(50*time.Millisecond),
			WithCommandName("some_command_name"),
			WithHystrixTimeout(10*time.Millisecond),
			WithMaxConcurrentRequests(100),
			WithErrorPercentThreshold(10),
			WithSleepWindow(100),
			WithRequestVolumeThreshold(20),
			WithStatsDCollector("localhost", "x"),
		)
	})
}

func shouldPanic(t *testing.T, f func()) {
	t.Helper()
	defer func() { _ = recover() }()
	f()
	t.Errorf("should have panicked")
}
