package wiki_client

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

// RoundTripFunc .
type RoundTripFunc func(req *http.Request) *http.Response

// RoundTrip .
func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

// NewTestClient returns *http.Client with Transport replaced to avoid making real calls
func NewTestClient(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: RoundTripFunc(fn),
	}
}

func TestClientStructGet(t *testing.T) {
	client := NewTestClient(func(req *http.Request) *http.Response {
		// Test request parameters
		assert.EqualValues(t, req.URL.String(), "/search")
		return &http.Response{
			StatusCode: 200,
			// Send response to be tested
			Body: ioutil.NopCloser(bytes.NewBufferString(`OK`)),
			// Must be set to non-nil value or it panics
			Header: make(http.Header),
		}
	})

	Url := "/search"
	response := &http.Response{
		StatusCode: 200,
		// Send response to be tested
		Body: ioutil.NopCloser(bytes.NewBufferString(`OK`)),
		// Must be set to non-nil value or it panics
		Header: make(http.Header),
	}

	actualResponse, actualError := client.Get(Url)

	assert.Equal(t, response, actualResponse)
	assert.Equal(t, nil, actualError)
}
