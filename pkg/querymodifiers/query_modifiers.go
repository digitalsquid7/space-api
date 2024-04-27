package querymodifiers

import "net/http"

type Option func(*URLParams, *QueryModifiers) error

type QueryModifiers struct {
	Page       *Page
	Filters    []Filter
	SortFields []SortField
}

// Load extracts query modifiers (paging, filtering, sorting) from the request URL parameters.
//
// Invalid user input is returned as InvalidInputError, so it can be handled separately are returned in a 400 response.
// The Go Options Pattern is used to let the caller configure which query modifiers to extract and validate.
func Load(req *http.Request, options ...Option) (*QueryModifiers, error) {
	requestParams := &QueryModifiers{}
	urlParams := NewURLParams(req.URL.Query())

	for _, option := range options {
		if err := option(urlParams, requestParams); err != nil {
			return nil, err
		}
	}

	return requestParams, nil
}
