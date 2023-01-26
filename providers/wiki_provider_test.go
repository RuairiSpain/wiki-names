package wiki_provider

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	w "wiki-names/clients/wiki_client"
	d "wiki-names/domains/wiki_domain"
)

var getContentMockFunc func(url string) (*http.Response, error)

type getClientMock struct{}

// We are mocking the client method "Get"
func (cm *getClientMock) Get(request string) (*http.Response, error) {
	return getContentMockFunc(request)
}

// When the everything is good
func TestGetContentNoError(t *testing.T) {
	// The error we will get is from the "response" so we make the second parameter of the function is nil
	getContentMockFunc = func(url string) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body: ioutil.NopCloser(strings.NewReader(`
			Content{
				Continue: {
					Rvcontinue : "json:rvcontinue",
					Continue   : "json:continue",
				},
				Warnings:  {
					Main : {
						Warnings : "json:warnings",
					},
					Revisions : {
						Warnings : "json:warnings",
					},
				},
				Query:  {
					Normalized: {[]{
						Fromencoded : true,
						From  :"json:from",
						To    :"json:to",
					}},
					Pages: {[]{
						Pageid    :    "json:pageid",
						Ns        :    1,
						Title     : "json:title",
						Revisions [ {
							Contentformat : "json:contentformat",
							Contentmodel  : "json:contentmodel",
							Content       : "{{Short description|Bob Smith is American}}",
						} ],
					}],
				},
			}`)),
		}, nil
	}
	w.Client = &getClientMock{} // without this line, the real api is fired

	response, err := WikiProvider.GetContent(d.WikiRequest{"Bob_Smith", "en"})
	assert.NotNil(t, response)
	assert.Nil(t, err)

	assert.EqualValues(t, "json:rvcontinue", response.Continue.Rvcontinue)
	assert.EqualValues(t, "json:continue", response.Continue.Continue)
	assert.EqualValues(t, "json:warnings", response.Warnings.Main.Warnings)
	assert.EqualValues(t, "json:warnings", response.Revisions.Main)

	assert.EqualValues(t, 1, len(response.Query.Normalized[0]))
	assert.EqualValues(t, true, response.Query.Normalized[0].Fromencoded)
	assert.EqualValues(t, "json:from", response.Query.Normalized[0].From)
	assert.EqualValues(t, "json:to", response.Query.Normalized[0].To)

	assert.EqualValues(t, 1, len(response.Query.Pages))
	assert.EqualValues(t, "json:pageid", response.Query.Pages[0].Pageid)
	assert.EqualValues(t, 1, response.Query.Pages[0].Ns)
	assert.EqualValues(t, "json:title", response.Query.Pages[0].Title)

	assert.EqualValues(t, 1, len(response.Query.Pages))
	assert.EqualValues(t, "json:contentformat", response.Query.Pages[0].Revisions.Contentformat)
	assert.EqualValues(t, "json:contentmodel", response.Query.Pages[0].Revisions.Contentmodel)
	assert.EqualValues(t, "Bob Smith is American", response.Query.Pages[0].Revisions.Content)
}

func TestGetContentNoShortDecription(t *testing.T) {
	// The error we will get is from the "response" so we make the second parameter of the function is nil
	getContentMockFunc = func(url string) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body: ioutil.NopCloser(strings.NewReader(`
			Content{
				Continue: {
					Rvcontinue : "json:rvcontinue",
					Continue   : "json:continue",
				},
				Warnings:  {
					Main : {
						Warnings : "json:warnings",
					},
					Revisions : {
						Warnings : "json:warnings",
					},
				},
				Query:  {
					Normalized: {[]{
						Fromencoded : true,
						From  :"json:from",
						To    :"json:to",
					}},
					Pages: {[]{
						Pageid    :    "json:pageid",
						Ns        :    1,
						Title     : "json:title",
						Revisions [ {
							Contentformat : "json:contentformat",
							Contentmodel  : "json:contentmodel",
							Content       : "{{descripcion corto|Bob Smith es Americano}}",
						} ],
					}],
				},
			}`)),
		}, nil
	}
	w.Client = &getClientMock{} // without this line, the real api is fired

	response, err := WikiProvider.GetContent(d.WikiRequest{"Bob_Smith", "en"})
	assert.NotNil(t, response)
	assert.Nil(t, err)
	assert.EqualValues(t, http.StatusPartialContent, response.StatusCode)
	assert.EqualValues(t, "{{descripcion corto|Bob Smith es Americano}}", response.Query.Pages[0].Revisions.Content)

	// Same as previous test
	assert.EqualValues(t, "json:rvcontinue", response.Continue.Rvcontinue)
	assert.EqualValues(t, "json:continue", response.Continue.Continue)
	assert.EqualValues(t, "json:warnings", response.Warnings.Main.Warnings)
	assert.EqualValues(t, "json:warnings", response.Revisions.Main)

	assert.EqualValues(t, 1, len(response.Query.Normalized[0]))
	assert.EqualValues(t, true, response.Query.Normalized[0].Fromencoded)
	assert.EqualValues(t, "json:from", response.Query.Normalized[0].From)
	assert.EqualValues(t, "json:to", response.Query.Normalized[0].To)

	assert.EqualValues(t, 1, len(response.Query.Pages))
	assert.EqualValues(t, "json:pageid", response.Query.Pages[0].Pageid)
	assert.EqualValues(t, 1, response.Query.Pages[0].Ns)
	assert.EqualValues(t, "json:title", response.Query.Pages[0].Title)

	assert.EqualValues(t, 1, len(response.Query.Pages))
	assert.EqualValues(t, "json:contentformat", response.Query.Pages[0].Revisions.Contentformat)
	assert.EqualValues(t, "json:contentmodel", response.Query.Pages[0].Revisions.Contentmodel)
}

func TestGetContentMissingPageRevisions(t *testing.T) {
	// The error we will get is from the "response" so we make the second parameter of the function is nil
	getContentMockFunc = func(url string) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body: ioutil.NopCloser(strings.NewReader(`
			Content{
				Continue: {
					Rvcontinue : "json:rvcontinue",
					Continue   : "json:continue",
				},
				Warnings:  {
					Main : {
						Warnings : "json:warnings",
					},
					Revisions : {
						Warnings : "json:warnings",
					},
				},
				Query:  {
					Normalized: {[]{
						Fromencoded : true,
						From  :"json:from",
						To    :"json:to",
					}},
					Pages: [],
				},
			}`)),
		}, nil
	}
	w.Client = &getClientMock{} // without this line, the real api is fired

	response, err := WikiProvider.GetContent(d.WikiRequest{"Bob_Smith", "en"})
	assert.Nil(t, response)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusNotFound, err.Code)
	assert.EqualValues(t, "Missing page revisions in json response body", err.Message)
}

func TestGetContentInvalidJsonError(t *testing.T) {
	// The error we will get is from the "response" so we make the second parameter of the function is nil
	getContentMockFunc = func(url string) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(strings.NewReader(`INVALID JSON`)),
		}, nil
	}
	w.Client = &getClientMock{} // without this line, the real api is fired

	request := d.WikiRequest{"Bob_Smith", "en"}
	got, err := WikiProvider.GetContent(request)

	assert.Nil(t, got)
	assert.NotNil(t, err)

	assert.EqualValues(t, http.StatusInternalServerError, err.Code)
	assert.EqualValues(t, "error unmarshaling wiki fetch response", err.ErrorMessage)
}

func TestGetContentInvalidWikiPermissions(t *testing.T) {
	getContentMockFunc = func(url string) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusForbidden,
			Body:       ioutil.NopCloser(strings.NewReader(`{"code": 403, "error": "permission denied"}`)),
		}, nil
	}
	w.Client = &getClientMock{} // without this line, the real api is fired

	got, err := WikiProvider.GetContent(d.WikiRequest{"Bob_Smith", "en"})

	assert.NotNil(t, err)
	assert.Nil(t, got)
	assert.EqualValues(t, http.StatusBadRequest, err.Code)
	assert.EqualValues(t, "permission denied", err.ErrorMessage)
}

func TestGetContentBadRequest(t *testing.T) {
	getContentMockFunc = func(url string) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusBadRequest,
			Body:       ioutil.NopCloser(strings.NewReader(`{"code": 400, "error": "The given location is invalid"}`)),
		}, nil
	}
	w.Client = &getClientMock{} // without this line, the real api is fired

	response, err := WikiProvider.GetContent(d.WikiRequest{"Bob_Smith", "en"})
	assert.NotNil(t, err)
	assert.Nil(t, response)
	assert.EqualValues(t, http.StatusBadRequest, err.Code)
	assert.EqualValues(t, "The given location is invalid", err.ErrorMessage)
}

func TestGetContentInvalidHttpBodyResponse(t *testing.T) {
	getContentMockFunc = func(url string) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusBadRequest,
			Body:       ioutil.NopCloser(strings.NewReader(`{"code": 400, "error": "The given location is invalid"}`)),
		}, nil
	}
	w.Client = &getClientMock{} // without this line, the real api is fired

	response, err := WikiProvider.GetContent(d.WikiRequest{"Bob_Smith", "en"})
	assert.NotNil(t, err)
	assert.Nil(t, response)
	assert.EqualValues(t, http.StatusBadRequest, err.Code)
	assert.EqualValues(t, "The given location is invalid", err.ErrorMessage)
}

type mockReadCloser struct {
	mock.Mock
}

func (m *mockReadCloser) Read(p []byte) (n int, err error) {
	args := m.Called(p)
	return args.Int(0), args.Error(1)
}

func (m *mockReadCloser) Close() error {
	args := m.Called()
	return args.Error(0)
}

// When no body is provided
func TestGetContentInvalidWikiResponseBody(t *testing.T) {
	mockReadCloser := mockReadCloser{}
	// if Read is called, it will return error
	mockReadCloser.On("Read", mock.AnythingOfType("[]uint8")).Return(0, fmt.Errorf("error reading"))
	// if Close is called, it will return error
	mockReadCloser.On("Close").Return(fmt.Errorf("error closing"))

	getContentMockFunc = func(url string) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       &mockReadCloser,
		}, nil
	}
	w.Client = &getClientMock{} // without this line, the real api is fired

	response, err := WikiProvider.GetContent(d.WikiRequest{"Bob_Smith", "en"})
	assert.NotNil(t, err)
	assert.Nil(t, response)
	assert.EqualValues(t, http.StatusBadRequest, err.Code)
	assert.EqualValues(t, "error reading", err.ErrorMessage)
}

func TestGetContentInvalidWikiResponseBody2(t *testing.T) {
	getContentMockFunc = func(url string) (*http.Response, error) {
		invalidCloser, _ := os.Open("-asf3")
		return &http.Response{
			StatusCode: http.StatusBadRequest,
			Body:       invalidCloser,
		}, nil
	}
	w.Client = &getClientMock{} // without this line, the real api is fired

	response, err := WikiProvider.GetContent(d.WikiRequest{"Bob_Smith", "en"})
	assert.Nil(t, response)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusBadRequest, err.Code)
	assert.EqualValues(t, "error reading", err.ErrorMessage)
}

// When the error response is invalid, here the code is supposed to be an integer, but a string was given.
// This can happen when the api owner changes some data types in the api
func TestGetWeatherInvalidErrorInterface(t *testing.T) {
	getContentMockFunc = func(url string) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(strings.NewReader(`{"XXX": "string code"}`)),
		}, nil
	}
	w.Client = &getClientMock{} // without this line, the real api is fired

	response, err := WikiProvider.GetContent(d.WikiRequest{"Bob_Smith", "en"})
	assert.Nil(t, response)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusInternalServerError, err.Code)
	assert.EqualValues(t, "invalid json response body", err.ErrorMessage)
}
