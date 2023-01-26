package wiki_provider

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"reflect"
	"strings"
	"testing"

	wiki_client "wiki-names/clients"
	wiki_domain "wiki-names/domains"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var getContentMockFunc func(url string) (*http.Response, error)

type getClientMock struct{}

// We are mocking the client method "Get"
func (cm *getClientMock) Get(request string) (*http.Response, error) {
	return getContentMockFunc(request)
}

func TestGetContentSummaryNoError(t *testing.T) {
	// The error we will get is from the "response" so we make the second parameter of the function is nil
	getContentMockFunc = func(url string) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(`{"continue":{"rvcontinue":"20221220124836|1128498933","continue":"||"},"warnings":{"main":{"warnings":"Subscribe to the mediawiki-api-announce mailing list at or notice of APIdeprecations and breaking changes. Use [[Special:ApiFeatureUsage]]to see usage of deprecated features by your application."},"revisions":{"warnings":"Because \"rvslots\" waformat is deprecated, and in the future the new format will always be used."}},"query":{"normalized":[{"fromencoded":false,"from":"Yoshua_Bengio","to":"Yoshua Bengio"}],"pages":[{"pageid":47749536,"ns":0,"title":"Yoshua Bengio","revisions":[{	"contentformat":"text/x-wiki",	"contentmodel":"wikitext","content":"{{Short description|Canadian computer scientist}}\n{{Use mdy dates|date=March 2019}}}"}]}]}}	`)),
		}, nil
	}
	wiki_client.Client = &getClientMock{} // without this line, the real api is fired

	response, err := WikiProvider.GetContentSummary(wiki_domain.RequestQuery{Name: "Bob_Smith", Locale: "en"})
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.EqualValues(t, response.ShortDescription, "Canadian computer scientist")
}

// When the everything is good
func TestGetContentNoError(t *testing.T) {
	// The error we will get is from the "response" so we make the second parameter of the function is nil
	getContentMockFunc = func(url string) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(`{"continue":{"rvcontinue":"20221220124836|1128498933","continue":"||"},"warnings":{"main":{"warnings":"Subscribe to the mediawiki-api-announce mailing list at or notice of APIdeprecations and breaking changes. Use [[Special:ApiFeatureUsage]]to see usage of deprecated features by your application."},"revisions":{"warnings":"Because \"rvslots\" waformat is deprecated, and in the future the new format will always be used."}},"query":{"normalized":[{"fromencoded":false,"from":"Yoshua_Bengio","to":"Yoshua Bengio"}],"pages":[{"pageid":47749536,"ns":0,"title":"Yoshua Bengio","revisions":[{"contentformat":"text/x-wiki","contentmodel":"wikitext","content":"{{Short description|Canadian computer scientist}}\n{{Use mdy dates|date=March 2019}}}"}]}]}}	`)),
		}, nil
	}
	wiki_client.Client = &getClientMock{} // without this line, the real api is fired

	response, err := WikiProvider.GetContent(wiki_domain.RequestQuery{Name: "Bob_Smith", Locale: "en"})
	assert.NotNil(t, response)
	assert.Nil(t, err)

	assert.EqualValues(t, "20221220124836|1128498933", response.Continue.Rvcontinue)
	assert.EqualValues(t, "||", response.Continue.Continue)
	assert.EqualValues(t, "Subscribe to the mediawiki-api-announce mailing list at or notice of APIdeprecations and breaking changes. Use [[Special:ApiFeatureUsage]]to see usage of deprecated features by your application.", response.Warnings.Main.Warnings)
	assert.EqualValues(t, "Because \"rvslots\" waformat is deprecated, and in the future the new format will always be used.", response.Warnings.Revisions.Warnings)

	assert.EqualValues(t, 1, len(response.Query.Normalized))
	assert.EqualValues(t, false, response.Query.Normalized[0].Fromencoded)
	assert.EqualValues(t, "Yoshua_Bengio", response.Query.Normalized[0].From)
	assert.EqualValues(t, "Yoshua Bengio", response.Query.Normalized[0].To)

	assert.EqualValues(t, 1, len(response.Query.Pages))
	assert.EqualValues(t, 47749536, response.Query.Pages[0].Pageid)
	assert.EqualValues(t, 0, response.Query.Pages[0].Ns)
	assert.EqualValues(t, "Yoshua Bengio", response.Query.Pages[0].Title)

	assert.EqualValues(t, 1, len(response.Query.Pages))
	assert.EqualValues(t, "text/x-wiki", response.Query.Pages[0].Revisions[0].Contentformat)
	assert.Contains(t, response.Query.Pages[0].Revisions[0].Content, "Canadian computer scientist")
}

func TestGetContentSummaryNoShortDecription(t *testing.T) {
	// The error we will get is from the "response" so we make the second parameter of the function is nil
	getContentMockFunc = func(url string) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(`{"continue":{"rvcontinue":"20221220124836|1128498933","continue":"||"},"warnings":{"main":{"warnings":"Subscribe to the mediawiki-api-announce mailing list at or notice of APIdeprecations and breaking changes. Use [[Special:ApiFeatureUsage]]to see usage of deprecated features by your application."},"revisions":{"warnings":"Because \"rvslots\" waformat is deprecated, and in the future the new format will always be used."}},"query":{"normalized":[{"fromencoded":false,"from":"Yoshua_Bengio","to":"Yoshua Bengio"}],"pages":[{"pageid":47749536,"ns":0,"title":"Yoshua Bengio","revisions":[{"contentformat":"text/x-wiki","contentmodel":"wikitext","content":"{{Descripcion corto|Canadian computer scientist}}\n{{Use mdy dates|date=March 2019}}}"}]}]}}	`)),
		}, nil
	}
	wiki_client.Client = &getClientMock{} // without this line, the real api is fired

	response, err := WikiProvider.GetContent(wiki_domain.RequestQuery{Name: "Bob_Smith", Locale: "en"})
	assert.NotNil(t, response)
	assert.Nil(t, err)
	assert.Contains(t, "{{Descripcion corto|Canadian computer scientist}}\n{{Use mdy dates|date=March 2019}}}", response.Query.Pages[0].Revisions[0].Content)

	// Same as previous test
	assert.EqualValues(t, "20221220124836|1128498933", response.Continue.Rvcontinue)
	assert.EqualValues(t, "||", response.Continue.Continue)
	assert.EqualValues(t, "Subscribe to the mediawiki-api-announce mailing list at or notice of APIdeprecations and breaking changes. Use [[Special:ApiFeatureUsage]]to see usage of deprecated features by your application.", response.Warnings.Main.Warnings)
	assert.EqualValues(t, "Because \"rvslots\" waformat is deprecated, and in the future the new format will always be used.", response.Warnings.Revisions.Warnings)

	assert.EqualValues(t, 1, len(response.Query.Normalized))
	assert.EqualValues(t, false, response.Query.Normalized[0].Fromencoded)
	assert.EqualValues(t, "Yoshua_Bengio", response.Query.Normalized[0].From)
	assert.EqualValues(t, "Yoshua Bengio", response.Query.Normalized[0].To)

	assert.EqualValues(t, 1, len(response.Query.Pages))
	assert.EqualValues(t, 47749536, response.Query.Pages[0].Pageid)
	assert.EqualValues(t, 0, response.Query.Pages[0].Ns)
	assert.EqualValues(t, "Yoshua Bengio", response.Query.Pages[0].Title)

	assert.EqualValues(t, 1, len(response.Query.Pages))
	assert.EqualValues(t, "text/x-wiki", response.Query.Pages[0].Revisions[0].Contentformat)
	assert.EqualValues(t, "wikitext", response.Query.Pages[0].Revisions[0].Contentmodel)
}

func TestGetContentSummaryMissingPageRevisions(t *testing.T) {
	// The error we will get is from the "response" so we make the second parameter of the function is nil
	getContentMockFunc = func(url string) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(`{"continue":{"rvcontinue":"20221220124836|1128498933","continue":"||"},"warnings":{"main":{	"warnings":"Subscribe to the mediawiki-api-announce mailing list at or notice of APIdeprecations and breaking changes. Use [[Special:ApiFeatureUsage]]to see usage of deprecated features by your application."},"revisions":{"warnings":"Because \"rvslots\" waformat is deprecated, and in the future the new format will always be used."}},"query":{"normalized":[{"fromencoded":false,"from":"Yoshua_Bengio","to":"Yoshua Bengio"}],"pages":[{"pageid":47749536,"ns":0,"title":"Yoshua Bengio","revisions":[]}]}}`)),
		}, nil
	}
	wiki_client.Client = &getClientMock{} // without this line, the real api is fired

	response, err := WikiProvider.GetContentSummary(wiki_domain.RequestQuery{Name: "Bob_Smith", Locale: "en"})
	assert.Nil(t, response)
	assert.NotNil(t, err)
	assert.EqualValues(t, 404, err.Code)
	assert.EqualValues(t, "Missing page revisions in json response body", err.ErrorMessage)
}

func TestGetContentInvalidJsonError(t *testing.T) {
	// The error we will get is from the "response" so we make the second parameter of the function is nil
	getContentMockFunc = func(url string) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(`INVALID JSON`)),
		}, nil
	}
	wiki_client.Client = &getClientMock{} // without this line, the real api is fired

	request := wiki_domain.RequestQuery{Name: "Bob_Smith", Locale: "en"}
	got, err := WikiProvider.GetContentSummary(request)

	assert.Nil(t, got)
	assert.NotNil(t, err)

	assert.EqualValues(t, http.StatusInternalServerError, err.Code)
	assert.EqualValues(t, "error unmarshaling wiki fetch response", err.ErrorMessage)
}

func TestGetContentInvalidWikiPermissions(t *testing.T) {
	getContentMockFunc = func(url string) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusForbidden,
			Body:       io.NopCloser(strings.NewReader(`{"code": 403, "error": "permission denied"}`)),
		}, nil
	}
	wiki_client.Client = &getClientMock{} // without this line, the real api is fired

	got, err := WikiProvider.GetContent(wiki_domain.RequestQuery{Name: "Bob_Smith", Locale: "en"})

	assert.NotNil(t, err)
	assert.Nil(t, got)
	assert.EqualValues(t, http.StatusInternalServerError, err.Code)
	assert.EqualValues(t, "invalid json response body: 403", err.ErrorMessage)
}

func TestGetContentBadRequest(t *testing.T) {
	getContentMockFunc = func(url string) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusBadRequest,
			Body:       io.NopCloser(strings.NewReader(`{"code": 400, "error": "The given location is invalid"}`)),
		}, nil
	}
	wiki_client.Client = &getClientMock{} // without this line, the real api is fired

	response, err := WikiProvider.GetContent(wiki_domain.RequestQuery{Name: "Bob_Smith", Locale: "en"})
	assert.NotNil(t, err)
	assert.Nil(t, response)
	assert.EqualValues(t, http.StatusInternalServerError, err.Code)
	assert.EqualValues(t, "invalid json response body: 400", err.ErrorMessage)
}

func TestGetContentSummaryInvalidHttpBodyResponse(t *testing.T) {
	getContentMockFunc = func(url string) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusBadRequest,
			Body:       io.NopCloser(strings.NewReader(`{"code": 400, "error": "The given location is invalid"}`)),
		}, nil
	}
	wiki_client.Client = &getClientMock{} // without this line, the real api is fired

	response, err := WikiProvider.GetContentSummary(wiki_domain.RequestQuery{Name: "Bob_Smith", Locale: "en"})
	assert.NotNil(t, err)
	assert.Nil(t, response)
	assert.EqualValues(t, http.StatusInternalServerError, err.Code)
	assert.EqualValues(t, "invalid json response body: 400", err.ErrorMessage)
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
	wiki_client.Client = &getClientMock{} // without this line, the real api is fired

	response, err := WikiProvider.GetContent(wiki_domain.RequestQuery{Name: "Bob_Smith", Locale: "en"})
	assert.NotNil(t, err)
	assert.Nil(t, response)
	assert.EqualValues(t, http.StatusBadRequest, err.Code)
	assert.EqualValues(t, "error reading", err.ErrorMessage)
}

func TestGetContentSummaryInvalidWikiResponseBody2(t *testing.T) {
	getContentMockFunc = func(url string) (*http.Response, error) {
		invalidCloser, _ := os.Open("-asf3")
		return &http.Response{
			StatusCode: http.StatusBadRequest,
			Body:       invalidCloser,
		}, nil
	}
	wiki_client.Client = &getClientMock{} // without this line, the real api is fired

	response, err := WikiProvider.GetContentSummary(wiki_domain.RequestQuery{Name: "Bob_Smith", Locale: "en"})
	assert.Nil(t, response)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusBadRequest, err.Code)
	assert.EqualValues(t, "invalid argument", err.ErrorMessage)
}

// When the error response is invalid, here the code is supposed to be an integer, but a string was given.
// This can happen when the api owner changes some data types in the api
func TestGetContentSummary404(t *testing.T) {
	getContentMockFunc = func(url string) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(`{"XXX": "string code"}`)),
		}, nil
	}
	wiki_client.Client = &getClientMock{} // without this line, the real api is fired

	response, err := WikiProvider.GetContentSummary(wiki_domain.RequestQuery{Name: "Bob_Smith", Locale: "en"})
	assert.Nil(t, response)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusNotFound, err.Code)
	assert.EqualValues(t, "Missing page revisions in json response body", err.ErrorMessage)
}

// I did not add unit tests for GetExtract API call.  They would be very similar to the GetContent unit tests.
// Given more time I would add them in a commercial setting.

func Test_getClientMock_Get(t *testing.T) {
	type args struct {
		request string
	}
	tests := []struct {
		name    string
		cm      *getClientMock
		args    args
		want    *http.Response
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.cm.Get(tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("getClientMock.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getClientMock.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_mockReadCloser_Read(t *testing.T) {
	type args struct {
		p []byte
	}
	tests := []struct {
		name    string
		m       *mockReadCloser
		args    args
		wantN   int
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotN, err := tt.m.Read(tt.args.p)
			if (err != nil) != tt.wantErr {
				t.Errorf("mockReadCloser.Read() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotN != tt.wantN {
				t.Errorf("mockReadCloser.Read() = %v, want %v", gotN, tt.wantN)
			}
		})
	}
}

func Test_mockReadCloser_Close(t *testing.T) {
	tests := []struct {
		name    string
		m       *mockReadCloser
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.m.Close(); (err != nil) != tt.wantErr {
				t.Errorf("mockReadCloser.Close() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
