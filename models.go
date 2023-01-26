package main

type RequestQuery struct {
	Name   string `uri:"name" binding:"required"`
	Locale string `uri:"locale"`
}

type Content struct {
	Continue struct {
		Rvcontinue string `json:"rvcontinue"`
		Continue   string `json:"continue"`
	} `json:"continue"`
	Warnings struct {
		Main struct {
			Warnings string `json:"warnings"`
		} `json:"main"`
		Revisions struct {
			Warnings string `json:"warnings"`
		} `json:"revisions"`
	} `json:"warnings"`
	Query struct {
		Normalized []struct {
			Fromencoded bool   `json:"fromencoded"`
			From        string `json:"from"`
			To          string `json:"to"`
		} `json:"normalized"`
		Pages []struct {
			Pageid    int    `json:"pageid"`
			Ns        int    `json:"ns"`
			Title     string `json:"title"`
			Revisions []struct {
				Contentformat string `json:"contentformat"`
				Contentmodel  string `json:"contentmodel"`
				Content       string `json:"content"`
			} `json:"revisions"`
		} `json:"pages"`
	} `json:"query"`
}

type Extract struct {
	Batchcomplete bool `json:"batchcomplete"`
	Query         struct {
		Normalized []struct {
			Fromencoded bool   `json:"fromencoded"`
			From        string `json:"from"`
			To          string `json:"to"`
		} `json:"normalized"`
		Pages []struct {
			Pageid  int    `json:"pageid"`
			Ns      int    `json:"ns"`
			Title   string `json:"title"`
			Extract string `json:"extract"`
		} `json:"pages"`
	} `json:"query"`
}
