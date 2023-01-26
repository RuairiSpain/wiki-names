package app

import (
	"github.com/stretchr/testify/assert"
	"http"
	"testing"
)

func TestSetupRouter(t *testing.T) {
	type testCase struct {
		Name string

		Address string

		ExpectedServer *http.Server
	}

	validate := func(t *testing.T, tc *testCase) {
		t.Run(tc.Name, func(t *testing.T) {
			actualServer := SetupRouter(tc.Address)

			assert.Equal(t, tc.ExpectedServer, actualServer)
		})
	}

	validate(t, &testCase{})
}
