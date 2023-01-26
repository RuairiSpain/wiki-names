package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// func Test_setupRouter(t *testing.T) {
// 	type args struct {
// 		handler *WikiHandler
// 	}
// 	tests := []struct {
// 		name string
// 		args args
// 		want *http.Server
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if got := setupRouter(tt.args.handler); !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("setupRouter() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

// func Test_main(t *testing.T) {
// 	tests := []struct {
// 		name string
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			main()
// 		})
// 	}
// }



func performRequest(r http.Server, method, path string) *httptest.ResponseRecorder {
   req, _ := http.NewRequest(method, path, nil)
   w := httptest.NewRecorder()
   r.Handler.ServeHTTP(w, req)
   return w
}
func TestHelloWorldSuccess(t *testing.T) {
   // Build our expected body
   expect := "Canadian computer scientist"
   // Grab our router
   handler := WikiHandler{http: &http.Server{}}
   router := SetupRouter(&handler)
   // Perform a GET request with that handler.
   w := performRequest(*router, "GET", "/search/Yoshua_Bengio")
   // Assert we encoded correctly,
   // the request gives a 200
   assert.Equal(t, http.StatusOK, w.Code)
   got := w.Body.String()
   assert.Equal(t, expect, got)
}

func TestHelloWorldFail(t *testing.T) {
   // Build our expected body
   body := "Canadian computer scientist"
   // Grab our router
   handler := WikiHandler{http: &http.Server{}}
   router := SetupRouter(&handler)
   // Perform a GET request with that handler.
   w := performRequest(*router, "GET", "/search/Yoshua_Bengio")
   // Assert we encoded correctly,
   // the request gives a 200
   assert.Equal(t, http.StatusOK, w.Code)
   response := w.Body.String()
   assert.Equal(t, body, response)
}