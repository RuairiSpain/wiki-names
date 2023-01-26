package wiki_domain

type RequestQuery struct {
	Name   string `uri:"name" binding:"required"`
	Locale string `uri:"locale"`
}
type Response struct {
	ShortDescription string `json:"short_description"`
}

type ContentRevision struct {
	Contentformat string `json:"contentformat"`
	Contentmodel  string `json:"contentmodel"`
	Content       string `json:"content"`
}
type PageRevision struct {
	Pageid    int               `json:"pageid"`
	Ns        int               `json:"ns"`
	Title     string            `json:"title"`
	Revisions []ContentRevision `json:"revisions"`
}

type Normalize struct {
	Fromencoded bool   `json:"fromencoded"`
	From        string `json:"from"`
	To          string `json:"to"`
}
type PageExtract struct {
	Pageid  int    `json:"pageid"`
	Ns      int    `json:"ns"`
	Title   string `json:"title"`
	Extract string `json:"extract"`
}

type ContinueType struct {
	Rvcontinue string `json:"rvcontinue"`
	Continue   string `json:"continue"`
}

type WarningsType struct {
	Main      WarningsSimpleType `json:"main"`
	Revisions WarningsSimpleType `json:"revisions"`
}
type QueryPageRevisionType struct {
	Normalized []Normalize    `json:"normalized"`
	Pages      []PageRevision `json:"pages"`
}
type QueryPageExtractType struct {
	Normalized []Normalize   `json:"normalized"`
	Pages      []PageExtract `json:"pages"`
}

type WarningsSimpleType struct {
	Warnings string `json:"warnings"`
}
type Content struct {
	Continue ContinueType          `json:"continue"`
	Warnings WarningsType          `json:"warnings"`
	Query    QueryPageRevisionType `json:"query"`
}

type Extract struct {
	Batchcomplete bool                 `json:"batchcomplete"`
	Query         QueryPageExtractType `json:"query"`
}
