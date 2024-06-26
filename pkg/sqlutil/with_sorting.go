package sqlutil

import (
	"fmt"
	"regexp"
	"space-api/pkg/models"
	"space-api/pkg/urlparams"
	"strings"
)

type Direction string

const (
	Ascending  Direction = "asc"
	Descending Direction = "desc"
)

type SortField struct {
	Field     Field
	Direction Direction
}

func WithSorting(fields Fields) func(*urlparams.URLParams, *QueryModifiers) error {
	return func(urlParams *urlparams.URLParams, reqParams *QueryModifiers) error {
		reqParams.SortFields = make([]SortField, 0)

		sort := urlParams.Values.Get("sort")
		if sort == "" {
			return nil
		}

		regex := `^[a-zA-Z]+ (asc|desc)(, [a-zA-Z]+ (asc|desc))*$`
		matched, err := regexp.MatchString(regex, sort)
		if err != nil {
			return fmt.Errorf("check regex match: %w", err)
		}

		if !matched {
			return models.NewInvalidInputError("sort", "must match the regex: "+regex)
		}

		columns := strings.Split(sort, ", ")

		for _, column := range columns {
			parts := strings.Split(column, " ")

			field, ok := fields.GetFieldByAPIName(parts[0])
			if !ok {
				return models.NewInvalidInputError("sort", parts[0]+" is not a sortable field")
			}

			sortField := SortField{Field: field, Direction: Direction(parts[1])}
			reqParams.SortFields = append(reqParams.SortFields, sortField)
		}

		return nil
	}
}
