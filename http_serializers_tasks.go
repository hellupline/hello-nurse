package main

import (
	"net/http"
)

type (
	BooruFetchFileRequest struct { // nolint
		*BooruFetchFile
	}

	BooruFetchFile struct { // nolint
		Type string `json:"type"`
		Key  string `json:"Key"`
	}

	BooruFetchTagRequest struct { // nolint
		*BooruFetchTask
	}

	BooruFetchTask struct { // nolint
		Domain string `json:"domain"`
		Tag    string `json:"tag"`
	}
)

func (s BooruFetchFileRequest) Bind(r *http.Request) error { // nolint
	return nil
}

func (s BooruFetchTagRequest) Bind(r *http.Request) error { // nolint
	return nil
}
