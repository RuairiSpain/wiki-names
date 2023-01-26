package wiki_controller

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	wiki_domain "wiki-names/domains"
	wiki_provider "wiki-names/providers"
)

func GetContentSummary(c *gin.Context) {
	var query wiki_domain.RequestQuery
	if err := c.ShouldBindUri(&query); err != nil {
		log.Fatalln("Missing name in query string")
		c.JSON(400, gin.H{"msg": err})
		return
	}
	result, apiError := wiki_provider.WikiProvider.GetContentSummary(query)
	if apiError != nil {
		c.JSON(apiError.Code, apiError)
		return
	}
	c.JSON(http.StatusOK, result)
}

func GetExtract(c *gin.Context) {
	var query wiki_domain.RequestQuery
	if err := c.ShouldBindUri(&query); err != nil {
		log.Fatalln("Missing name in query string")
		c.JSON(400, gin.H{"msg": err})
		return
	}
	// Use a default language if not set in the URL
	if query.Locale == "" {
		query.Locale = "en"
	}
	result, apiError := wiki_provider.WikiProvider.GetExtract(query)
	if apiError != nil {
		c.JSON(apiError.Code, apiError)
		return
	}
	c.JSON(http.StatusOK, result)
}
