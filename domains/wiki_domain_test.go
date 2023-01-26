package wiki_domain

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Content(t *testing.T) {
	request := Content{
		Continue: ContinueType{
			Rvcontinue: "json:rvcontinue",
			Continue:   "json:continue",
		},
		Warnings: WarningsType{
			Main: WarningsSimpleType{
				Warnings: "json:warnings",
			},
			Revisions: WarningsSimpleType{
				Warnings: "json:warnings",
			},
		},
		Query: QueryPageRevisionType{
			Normalized: []Normalize,{
				Fromencoded: true,
				From:        "json:from",
				To:          "json:to",
			},
			Pages: []RevisionPage{
				Pageid: "json:pageid",
				Ns:     1,
				Title:  "json:title",
				Revisions: []ContentRevision{
					Contentformat: "json:contentformat",
					Contentmodel:  "json:contentmodel",
					Content:       "json:content",
				},
			},
		},
	}
	bytes, err := json.Marshal(request)
	assert.Nil(t, err)
	assert.NotNil(t, bytes)

	var result Content
	if err = json.Unmarshal(bytes, &result); err != nil {
		assert.Nil(t, err)
	}

	assert.EqualValues(t, result.Continue.Rvcontinue, request.Continue.Rvcontinue)
	assert.EqualValues(t, result.Continue.Continue, request.Continue.Continue)
	assert.EqualValues(t, result.Warnings.Main.Warnings, request.Warnings.Main.Warnings)
	assert.EqualValues(t, result.Warnings.Revisions.Warnings, request.Warnings.Revisions.Warnings)

	assert.EqualValues(t, len(result.Query.Normalized), len(request.Query.Normalized))
	assert.EqualValues(t, result.Query.Normalized[0].Fromencoded, request.Query.Normalized[0].Fromencoded)
	assert.EqualValues(t, result.Query.Normalized[0].From, request.Query.Normalized[0].From)
	assert.EqualValues(t, result.Query.Normalized[0].To, request.Query.Normalized[0].To)

	assert.EqualValues(t, len(result.Query.Pages), len(request.Query.Pages))
	assert.EqualValues(t, result.Query.Pages[0].Pageid, request.Query.Pages[0].Pageid)
	assert.EqualValues(t, result.Query.Pages[0].Ns, request.Query.Pages[0].Ns)
	assert.EqualValues(t, result.Query.Pages[0].Title, request.Query.Pages[0].Title)

	assert.EqualValues(t, len(result.Query.Pages), len(request.Query.Pages))
	assert.EqualValues(t, result.Query.Pages[0].Revisions.Contentformat, request.Query.Pages[0].Revisions[0].Contentformat)
	assert.EqualValues(t, result.Query.Pages[0].Revisions.Contentmodel, request.Query.Pages[0].Revisions[0].Contentmodel)
	assert.EqualValues(t, result.Query.Pages[0].Revisions.Content, request.Query.Pages[0].Revisions[0].Content)
}

func Test_ContentError(t *testing.T) {
	request := WikiError{
		Code:         400,
		ErrorMessage: "Bad Request Error",
	}
	bytes, err := json.Marshal(request)
	assert.Nil(t, err)
	assert.NotNil(t, bytes)

	var errResult WikiError
	err = json.Unmarshal(bytes, &errResult)
	assert.Nil(t, err)
	assert.EqualValues(t, errResult.Code, request.Code)
	assert.EqualValues(t, errResult.ErrorMessage, request.ErrorMessage)
}
