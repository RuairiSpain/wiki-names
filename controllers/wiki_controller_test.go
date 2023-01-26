package wiki_controller

import (
	"gin"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetContentSummary(t *testing.T) {
	type testCase struct {
		Name string

		C *gin.Context
	}

	validate := func(t *testing.T, tc *testCase) {
		t.Run(tc.Name, func(t *testing.T) {
			GetContentSummary(tc.C)
		})
	}

	validate(t, &testCase{})
}

func TestGetExtract(t *testing.T) {
	type testCase struct {
		Name string

		C *gin.Context
	}

	validate := func(t *testing.T, tc *testCase) {
		t.Run(tc.Name, func(t *testing.T) {
			GetExtract(tc.C)
		})
	}

	validate(t, &testCase{})
}
