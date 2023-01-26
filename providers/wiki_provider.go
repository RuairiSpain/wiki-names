package wiki_provider

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"

	wiki_client "wiki-names/clients"
	wiki_domain "wiki-names/domains"
)

const (
	contentUrl = "https://en.wikipedia.org/w/api.php?action=query&prop=revisions&titles=%s&rvlimit=1&formatversion=2&format=json&rvprop=content"
	extractUrl = "https://%s.wikipedia.org/w/api.php?action=query&format=json&prop=extracts&titles=%s&formatversion=2&exsentences=2&exlimit=1&explaintext=1"
)

type WikiProviderStruct struct{}

type wikiServiceInterface interface {
	GetContent(request wiki_domain.RequestQuery) (*wiki_domain.Content, *wiki_domain.WikiError)
	GetContentSummary(request wiki_domain.RequestQuery) (*wiki_domain.Response, *wiki_domain.WikiError)
	GetExtract(request wiki_domain.RequestQuery) (*wiki_domain.Response, *wiki_domain.WikiError)
}

var WikiProvider wikiServiceInterface = &WikiProviderStruct{}

func (p *WikiProviderStruct) GetContent(request wiki_domain.RequestQuery) (*wiki_domain.Content, *wiki_domain.WikiError) {
	url := fmt.Sprintf(contentUrl, request.Name)
	response, err := wiki_client.Client.Get(url)
	if err != nil {
		log.Printf("error when trying to get wiki content %s", err.Error())
		return nil, &wiki_domain.WikiError{
			Code:         http.StatusBadRequest,
			ErrorMessage: err.Error(),
		}
	}
	bytes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, &wiki_domain.WikiError{
			Code:         http.StatusBadRequest,
			ErrorMessage: err.Error(),
		}
	}
	defer response.Body.Close()
	if err := checkStatusCode(response); err != nil {
		return nil, err
	}
	var result wiki_domain.Content
	if err := json.Unmarshal(bytes, &result); err != nil {
		log.Printf("error when trying to unmarshal wiki content successful response: %s", err.Error())
		return nil, &wiki_domain.WikiError{Code: http.StatusInternalServerError, ErrorMessage: "error unmarshaling wiki fetch response"}
	}

	return &result, nil
}

func (p *WikiProviderStruct) GetContentSummary(request wiki_domain.RequestQuery) (*wiki_domain.Response, *wiki_domain.WikiError) {
	result, err := p.GetContent(request)
	if err != nil {
		return nil, err
	}
	// Validate the JSON hierarchy structure
	if len(result.Query.Pages) == 0 ||
		len(result.Query.Pages[0].Revisions) == 0 {
		message := "Missing page revisions in json response body"
		log.Println(message)
		return nil, &wiki_domain.WikiError{
			Code:         http.StatusNotFound,
			ErrorMessage: message,
		}
	}
	// Find the Short Description placeholder text and extract the value using regex
	re := regexp.MustCompile(`{{Short description\|([^}][^}]+)}}`)
	matches := re.FindStringSubmatch(result.Query.Pages[0].Revisions[0].Content)
	if len(matches) < 1 {
		message := "Missing `Short description` in json response body"
		log.Println(message)
		/* ***NOTE***: If the "ShortDescription text is missing, I took the "executive decision"
		   to return a 206 status code with the full Content text in the HTTP BODY
			This needs to be fleshed out in a Agile daily meeting and align with PM and TechLead.
			I'm not overly comfortable with this inelegant solution.  If there is a standard for
			"short description"and it's translation, I would  do a version 2, with a mapping
			function to include the transaltions for the different language markup.

			As a version 1.5, I added an API endpoint to call Extract as an altertative for multilingual text,
			and get the first two sentences.

			This could be a solution for the multilinugal solution.  But I found that the # of
			sentence parameter has some strange behaviour, it seems to count periods (fullstops)
			without checking if it is part of a acronym or shorthand wording.  I'd recommend, we revisit
			the sentence count from "Extract" API to come up with a better way to delimit sentence
			boundaries.
		*/
		return nil, &wiki_domain.WikiError{
			Code:         http.StatusPartialContent,
			ErrorMessage: result.Query.Pages[0].Revisions[0].Content,
		}
	}
	return &wiki_domain.Response{ShortDescription: matches[1]}, nil
}

func (p *WikiProviderStruct) GetExtract(request wiki_domain.RequestQuery) (*wiki_domain.Response, *wiki_domain.WikiError) {
	url := fmt.Sprintf(extractUrl, request.Locale, request.Name)
	response, err := wiki_client.Client.Get(url)
	if err != nil {
		log.Printf("error when trying to get wiki extract %s", err.Error())
		return nil, &wiki_domain.WikiError{
			Code:         http.StatusBadRequest,
			ErrorMessage: err.Error(),
		}
	}
	bytes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, &wiki_domain.WikiError{
			Code:         http.StatusBadRequest,
			ErrorMessage: err.Error(),
		}
	}
	defer response.Body.Close()

	error := checkStatusCode(response)
	if error != nil {
		return nil, error
	}
	var result wiki_domain.Extract
	if err := json.Unmarshal(bytes, &result); err != nil {
		log.Printf("error when trying to unmarshal wiki extract successful response: %s", err.Error())
		return nil, &wiki_domain.WikiError{Code: http.StatusInternalServerError, ErrorMessage: "error unmarshaling wiki fetch response"}
	}

	description := result.Query.Pages[0].Extract
	return &wiki_domain.Response{ShortDescription: description}, nil
}

func checkStatusCode(response *http.Response) *wiki_domain.WikiError {
	// The api owner can decide to change datatypes, etc. When this happen, it might affect the error format returned
	if response.StatusCode > 299 {
		message := fmt.Sprintf("invalid json response body: %d", response.StatusCode)
		log.Println(message)
		return &wiki_domain.WikiError{
			Code:         http.StatusInternalServerError,
			ErrorMessage: message,
		}
	}
	return nil
}
