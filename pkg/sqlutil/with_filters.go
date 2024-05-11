package sqlutil

import (
	"fmt"
	"regexp"
	"sort"
	"space-api/pkg/models"
	"space-api/pkg/urlparams"
	"strconv"
	"time"
)

type Operator string

const (
	Equals Operator = "eq"
	GTE    Operator = "gte"
	LTE    Operator = "lte"
)

type Filter struct {
	Field    Field
	Operator Operator
	Value    any
}

func WithFilters(fields Fields) func(*urlparams.URLParams, *QueryModifiers) error {
	return func(urlParams *urlparams.URLParams, reqParams *QueryModifiers) error {
		reqParams.Filters = make([]Filter, 0)
		regex, err := regexp.Compile(`^(?P<field>[a-zA-Z]+)\[(?P<operator>eq|lte|gte)]$`)
		if err != nil {
			return fmt.Errorf("compile regex: %w", err)
		}

		for key, values := range urlParams.Values {
			if len(values) == 0 {
				continue
			}

			submatches := regex.FindStringSubmatch(key)
			if len(submatches) != 3 {
				continue
			}

			operator := Operator(submatches[2])

			field, ok := fields.GetFieldByAPIName(submatches[1])
			if !ok {
				return models.NewInvalidInputError(submatches[1], "not a filterable field")
			}

			sqlValue, err := convertToSQLValue(field, values[0])
			if err != nil {
				return err
			}

			reqParams.Filters = append(reqParams.Filters, Filter{
				Field:    field,
				Operator: operator,
				Value:    sqlValue,
			})
		}

		sort.Slice(reqParams.Filters, func(i, j int) bool {
			return reqParams.Filters[i].Field.APIName < reqParams.Filters[j].Field.APIName
		})

		return nil
	}
}

func convertToSQLValue(field Field, apiValue string) (any, error) {
	if field.Type == String {
		return apiValue, nil
	}

	if field.Type == Integer {
		converted, err := strconv.Atoi(apiValue)
		if err != nil {
			return nil, models.NewInvalidInputError(field.APIName, "must be an integer")
		}
		return converted, nil
	}

	if field.Type == Date {
		converted, err := time.Parse("2006-01-02", apiValue)
		if err != nil {
			return nil, models.NewInvalidInputError(field.APIName, "date must be in the format YYYY-MM-DD")
		}
		return converted, nil
	}

	return nil, fmt.Errorf("invalid field type provided: %s", field.Type)
}
