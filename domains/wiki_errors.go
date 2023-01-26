package wiki_domain

import (
	"encoding/json"
	"net/http"
)

type WikiErrorInterface interface {
	Status() int
	Message() string
}
type WikiError struct {
	Code      int           `json:"code"`
	ErrorMessage     string  `json:"error"`
}

func (w *WikiError) Status() int {
	return w.Code
}
func (w *WikiError) Message() string {
	return w.ErrorMessage
}

func NewWikiError(statusCode int, message string) WikiErrorInterface {
	return &WikiError{
		Code:         statusCode,
		ErrorMessage: message,
	}
}
func NewBadRequestError(message string) WikiErrorInterface {
	return &WikiError{
		Code: http.StatusBadRequest,
		ErrorMessage: message,
	}
}

func NewForbiddenError(message string) WikiErrorInterface {
	return &WikiError{
		Code: http.StatusForbidden,
		ErrorMessage: message,
	}
}

func NewApiErrFromBytes(body []byte) (WikiErrorInterface, error) {
	var result WikiError
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}
	return &result, nil
}