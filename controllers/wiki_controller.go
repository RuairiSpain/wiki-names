package wiki_controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	client "wiki-names/clients/wiki_client"
	domain "wiki-names/domains/wiki_domain"
	provider "wiki-names/providers/"
)

func GetWeather(c *gin.Context) {
	long, _ := strconv.ParseFloat(c.Param("longitude"), 64)
	lat, _ := strconv.ParseFloat(c.Param("latitude"), 64)
	request := domain.WikiRequest{
		Name: c.Param(""),
		Locale: "en"

	}
	result, apiError := provider.GetContentSummary(request)
	if apiError != nil {
		c.JSON(apiError.Status(), apiError)
		return
	}
	c.JSON(http.StatusOK, result)
}
