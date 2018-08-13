package nursehttp

import (
	"net/http"
)

type (
	BooruFetchFileRequest struct { // nolint: golint
		*BooruFetchFile
	}

	BooruFetchFile struct { // nolint: golint
		Type string `json:"type"`
		Key  string `json:"Key"`
	}

	BooruFetchTagRequest struct { // nolint: golint
		*BooruFetchTask
	}

	BooruFetchTask struct { // nolint: golint
		Domain string `json:"domain"`
		Tag    string `json:"tag"`
	}
)

func (s BooruFetchFileRequest) Bind(r *http.Request) error { // nolint: golint
	return nil
}

func (s BooruFetchTagRequest) Bind(r *http.Request) error { // nolint: golint
	return nil
}
