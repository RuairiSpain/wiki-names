package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
)

type ResponseService interface {
	GetWikiResponse(c *gin.Context, endpoint_pattern string) *http.Response
}
type ContentService interface {
	GetContent(c *gin.Context)
}
type ExtractService interface {
	GetExtract(c *gin.Context)
}
type WikiHandler struct {
	http *http.Server
}

type HTTPClient interface {
	Get(req *http.Request) (*http.Response, error)
}

func (w *WikiHandler) GetWikiResponse(c *gin.Context, endpoint_pattern string) *http.Response {
	var query RequestQuery
	if err := c.ShouldBindUri(&query); err != nil {
		log.Fatalln("Missing name in query string")
		c.JSON(400, gin.H{"msg": err})
		return nil
	}

	// Use a default language if not set in the URL
	if query.Locale == "" {
		query.Locale = "en"
	}

	endpoint := os.Getenv(endpoint_pattern)
	if endpoint == "" || !strings.Contains(endpoint, "PLACEHOLDER") {
		log.Fatalln("Missing valid env variable: " + endpoint_pattern)
		c.AbortWithStatus(500)
		return nil
	}

	// Replace placeholders with correct strings
	url := strings.Replace(endpoint, "PLACEHOLDER", query.Name, -1)
	url = strings.Replace(url, "LOCALE", query.Locale, -1)

	// Make GET request to Wikipedia API endpoint
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
		c.AbortWithError(500, err)
		return nil
	}

	return resp

}

func (w *WikiHandler) GetExtract(c *gin.Context) {
	resp := w.GetWikiResponse(c, "WIKI_EXTRACT_ENDPOINT")
	if resp == nil {
		return
	}

	// Parse the JSON from the Content endpoint
	defer resp.Body.Close()
	var result Extract
	err := json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		log.Fatalln(err)
		c.AbortWithError(500, err)
		return
	}

	// Validate the JSON hierarchy structure
	if len(result.Query.Pages) == 0 {
		c.AbortWithStatus(404)
		return
	}

	// Find the Extract property and return it's value
	c.String(200, result.Query.Pages[0].Extract)
}

func (w *WikiHandler) GetContent(c *gin.Context) {
	resp := w.GetWikiResponse(c, "WIKI_CONTENT_ENDPOINT")
	if resp == nil {
		return
	}

	// Parse the JSON from the Content endpoint
	defer resp.Body.Close()
	var result Content
	err := json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		log.Fatalln(err)
		c.AbortWithError(500, err)
		return
	}

	// Validate the JSON hierarchy structure
	if len(result.Query.Pages) == 0 ||
		len(result.Query.Pages[0].Revisions) == 0 {
		c.AbortWithStatus(404)
		return
	}

	// Find the Short Description placeholder text and extract the value using regex
	re := regexp.MustCompile(`{{Short description\|([^}][^}]+)}}`)
	matches := re.FindStringSubmatch(result.Query.Pages[0].Revisions[0].Content)
	if len(matches) < 1 {
		c.AbortWithStatus(404)
		return
	}

	c.String(200, matches[1])
}
