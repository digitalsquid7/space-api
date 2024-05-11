package sqlutil

import (
	"net/http"
	"space-api/pkg/urlparams"
)

type Option func(*urlparams.URLParams, *QueryModifiers) error

type QueryModifiers struct {
	Page       *Page
	Filters    []Filter
	SortFields []SortField
}

// LoadQueryModifiers extracts query modifiers (paging, filtering, sorting) from the request URL parameters.
//
// Invalid user input is returned as InvalidInputError, so it can be handled separately are returned in a 400 response.
// The Go Options Pattern is used to let the caller configure which query modifiers to extract and validate.
func LoadQueryModifiers(req *http.Request, options ...Option) (*QueryModifiers, error) {
	requestParams := &QueryModifiers{}
	urlParams := urlparams.New(req.URL.Query())

	for _, option := range options {
		if err := option(urlParams, requestParams); err != nil {
			return nil, err
		}
	}

	return requestParams, nil
}
