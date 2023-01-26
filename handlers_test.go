package main

import (
	"net/http"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestWikiHandler_GetWikiResponse(t *testing.T) {
	type fields struct {
		http *http.Server
	}
	type args struct {
		c                *gin.Context
		endpoint_pattern string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *http.Response
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &WikiHandler{
				http: tt.fields.http,
			}
			if got := w.GetWikiResponse(tt.args.c, tt.args.endpoint_pattern); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WikiHandler.GetWikiResponse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWikiHandler_GetExtract(t *testing.T) {
	type fields struct {
		http *http.Server
	}
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &WikiHandler{
				http: tt.fields.http,
			}
			w.GetExtract(tt.args.c)
		})
	}
}

func TestWikiHandler_GetContent(t *testing.T) {
	type fields struct {
		http *http.Server
	}
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &WikiHandler{
				http: tt.fields.http,
			}
			w.GetContent(tt.args.c)
		})
	}
}
