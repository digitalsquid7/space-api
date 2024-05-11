package urlparams

import (
	"net/url"
	"space-api/pkg/models"
	"strconv"
)

// URLParams is a decorator for retrieving URL parameters that provides a default value if not supplied in the request.
type URLParams struct {
	Values url.Values
}

func New(values url.Values) *URLParams {
	return &URLParams{Values: values}
}

func (u *URLParams) GetInt(field string, defaultValue int) (int, error) {
	value := u.Values.Get(field)
	if value == "" {
		return defaultValue, nil
	}

	converted, err := strconv.Atoi(value)
	if err != nil {
		return 0, models.NewInvalidInputError(field, "must be an integer")
	}

	return converted, nil
}

func (u *URLParams) GetBool(field string) (bool, error) {
	value := u.Values.Get(field)
	if value == "" {
		return false, nil
	}

	converted, err := strconv.ParseBool(value)
	if err != nil {
		return false, models.NewInvalidInputError(field, "must be a boolean")
	}

	return converted, nil
}
