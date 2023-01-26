package wiki_provider

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"

	client "wiki-names/clients/wiki_client"
	domain "wiki-names/domains/wiki_domain"
)

const (
	contentUrl = "https://en.wikipedia.org/w/api.php?action=query&prop=revisions&titles=%s&rvlimit=1&formatversion=2&format=json&rvprop=content"
	extractUrl = "https://%s.wikipedia.org/w/api.php?action=query&format=json&prop=extracts&titles=%s&formatversion=2&exsentences=2&exlimit=1&explaintext=1"
)

type wikiProvider struct{}

type wikiServiceInterface interface {
	GetExtract(request domain.RequestQuery) (*domain.Wiki, *domain.WikiError)
	GetContent(request domain.RequestQuery) (*domain.Wiki, *domain.WikiError)
	GetContentSummary(request domain.RequestQuery) (*domain.WikiResponse, *domain.WikiError)
}

var WikiProvider wikiServiceInterface = &wikiProvider{}

func (p *wikiProvider) GetExtract(request domain.RequestQuery) (*domain.Extract, *domain.WikiError) {
	url := fmt.Sprintf(extractUrl, request.Name)
	response, err := client.ClientStruct.Get(url)
	if err != nil {
		log.Println(fmt.Sprintf("error when trying to get wiki extract %s", err.Error()))
		return nil, &domain.WikiError{
			Code:         http.StatusBadRequest,
			ErrorMessage: err.Error(),
		}
	}
	bytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, &domain.WikiError{
			Code:         http.StatusBadRequest,
			ErrorMessage: err.Error(),
		}
	}
	defer response.Body.Close()

	err = checkStatusCode(response)
	if err != nil {
		return nil, err
	}
	var result domain.Extract
	if err := json.Unmarshal(bytes, &result); err != nil {
		log.Println(fmt.Sprintf("error when trying to unmarshal wiki extract successful response: %s", err.Error()))
		return nil, &domain.WikiError{Code: http.StatusInternalServerError, ErrorMessage: "error unmarshaling wiki fetch response"}
	}
	return &result, nil
}

func (p *wikiProvider) GetContent(request domain.RequestQuery) (*domain.Content, *domain.WikiError) {
	url := fmt.Sprintf(extractUrl, request.Locale, request.Name)
	response, err := client.Client.Get(url)
	if err != nil {
		log.Println(fmt.Sprintf("error when trying to get wiki content %s", err.Error()))
		return nil, &domain.WikiError{
			Code:         http.StatusBadRequest,
			ErrorMessage: err.Error(),
		}
	}
	bytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, &domain.WikiError{
			Code:         http.StatusBadRequest,
			ErrorMessage: err.Error(),
		}
	}
	defer response.Body.Close()
	if err := checkStatusCode(response); err != nil {
		return nil, err
	}
	var result domain.Content
	if err := json.Unmarshal(bytes, &result); err != nil {
		log.Println(fmt.Sprintf("error when trying to unmarshal wiki content successful response: %s", err.Error()))
		return nil, &domain.WikiError{Code: http.StatusInternalServerError, ErrorMessage: "error unmarshaling wiki fetch response"}
	}

	return &result, nil
}

func checkStatusCode(response *http.Response) *domain.WikiError {
	// The api owner can decide to change datatypes, etc. When this happen, it might affect the error format returned
	if response.StatusCode > 299 {
		message := fmt.Sprintf("invalid json response body: %d", response.StatusCode)
		log.Println(message)
		return &domain.WikiError{
			Code:         http.StatusInternalServerError,
			ErrorMessage: message,
		}
	}
	return nil
}

func (p *wikiProvider) GetContentSummary(request domain.RequestQuery) (*domain.WikiResponse, *domain.WikiError) {
	result, err := p.GetContent(request)
	if err != nil {
		return nil, err
	}
	// Validate the JSON hierarchy structure
	if len(result.Query.Pages) == 0 ||
		len(result.Query.Pages[0].Revisions) == 0 {
		message := fmt.Sprintf("Missing page revisions in json response body")
		log.Println(message)
		return nil, &domain.WikiError{
			Code:         http.StatusNotFound,
			ErrorMessage: message,
		}
	}
	// Find the Short Description placeholder text and extract the value using regex
	re := regexp.MustCompile(`{{Short description\|([^}][^}]+)}}`)
	matches := re.FindStringSubmatch(result.Query.Pages[0].Revisions[0].Content)
	if len(matches) < 1 {
		message := fmt.Sprintf("Missing `Short description` in json response body")
		log.Println(message)
		/* NOTE: If the "ShortDescription text is missing, I took the "executive decision"
		   to return a 206 status code with the full Content text in the HTTP BODY
			This needs to be fleshed out in a Agile daily meeting and align with PM and TechLead.
			I'm not overly comfortable with this inelegant solution.  If there is a standard for
			"short description"and it's translation, I would  do a version 2, with a mapping
			function to include the transaltions for the different language markup.

			As a version 1.5, I added an API endpoint to call Extract and get the first two sentences.
			This could be a solution for the multilinugal solution.  But I found that the # of
			sentence parameter has some strange behaviour, it seems to count periods (fullstops)
			without checking if it is part of a acronym or shorthand wording.  I'd recommend, we revisit
			the sentence count from "Extract" API to come up with a better way to delimit sentence
			boundaries.
		*/
		return nil, &domain.WikiError{
			Code:         http.StatusPartialContent,
			ErrorMessage: result.Query.Pages[0].Revisions[0].Content,
		}
	}
	return &domain.WikiResponse{Text: matches[1]}, nil
}
