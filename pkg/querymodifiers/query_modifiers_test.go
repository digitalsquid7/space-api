package querymodifiers

import (
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var fields = Fields{
	fields: []Field{
		{
			SQLName: "id",
			APIName: "id",
			Type:    Integer,
		},
		{
			SQLName: "account_name",
			APIName: "accountName",
			Type:    String,
		},
		{
			SQLName: "birth_dt",
			APIName: "birthDate",
			Type:    Date,
		},
	},
}

func TestLoad(t *testing.T) {
	tests := []struct {
		name           string
		url            string
		options        []Option
		expectedResult *QueryModifiers
		expectedErr    error
	}{
		{
			name:    "apply page number and size",
			url:     "?pageNumber=5&pageSize=20",
			options: []Option{WithPaging(10)},
			expectedResult: &QueryModifiers{
				Page: &Page{
					Size:   20,
					Number: 5,
				},
			},
		},
		{
			name:    "apply default page size",
			url:     "?pageNumber=5",
			options: []Option{WithPaging(10)},
			expectedResult: &QueryModifiers{
				Page: &Page{
					Size:   10,
					Number: 5,
				},
			},
		},
		{
			name:    "apply include total size",
			url:     "?includeTotalSize=true",
			options: []Option{WithPaging(10)},
			expectedResult: &QueryModifiers{
				Page: &Page{
					Size:             10,
					Number:           1,
					IncludeTotalSize: true,
				},
			},
		},
		{
			name:    "apply sort",
			url:     "?sort=accountName asc, id desc",
			options: []Option{WithSorting(fields)},
			expectedResult: &QueryModifiers{
				SortFields: []SortField{
					{
						Field: Field{
							SQLName: "account_name",
							APIName: "accountName",
							Type:    String,
						},
						Direction: Ascending,
					},
					{
						Field: Field{
							SQLName: "id",
							APIName: "id",
							Type:    Integer,
						},
						Direction: Descending,
					},
				},
			},
		},
		{
			name:    "apply filter",
			url:     "?accountName[eq]=Squid&birthDate[lte]=2009-09-20&id[gte]=8477",
			options: []Option{WithFilters(fields)},
			expectedResult: &QueryModifiers{
				Filters: []Filter{
					{
						Field: Field{
							SQLName: "account_name",
							APIName: "accountName",
							Type:    String,
						},
						Operator: Equals,
						Value:    "Squid",
					},
					{
						Field: Field{
							SQLName: "birth_dt",
							APIName: "birthDate",
							Type:    Date,
						},
						Operator: LTE,
						Value:    time.Date(2009, 9, 20, 0, 0, 0, 0, time.UTC),
					},
					{
						Field: Field{
							SQLName: "id",
							APIName: "id",
							Type:    Integer,
						},
						Operator: GTE,
						Value:    8477,
					},
				},
			},
		},
		{
			name: "apply filter, sort and page",
			url:  "?accountName[eq]=Squid&pageSize=50&sort=birthDate desc",
			options: []Option{
				WithFilters(fields),
				WithSorting(fields),
				WithPaging(10),
			},
			expectedResult: &QueryModifiers{
				Filters: []Filter{
					{
						Field: Field{
							SQLName: "account_name",
							APIName: "accountName",
							Type:    String,
						},
						Operator: Equals,
						Value:    "Squid",
					},
				},
				Page: &Page{
					Size:   50,
					Number: 1,
				},
				SortFields: []SortField{
					{
						Field: Field{
							SQLName: "birth_dt",
							APIName: "birthDate",
							Type:    Date,
						},
						Direction: Descending,
					},
				},
			},
		},
		{
			name:        "invalid page number",
			url:         "?pageNumber=five",
			options:     []Option{WithPaging(10)},
			expectedErr: NewInvalidInputError("pageNumber", "must be an integer"),
		},
		{
			name:        "invalid page size",
			url:         "?pageSize=five",
			options:     []Option{WithPaging(10)},
			expectedErr: NewInvalidInputError("pageSize", "must be an integer"),
		},
		{
			name:        "invalid include total size",
			url:         "?includeTotalSize=invalidvalue",
			options:     []Option{WithPaging(10)},
			expectedErr: NewInvalidInputError("includeTotalSize", "must be a boolean"),
		},
		{
			name:        "invalid sort format",
			url:         "?sort=accountName asc,, id ",
			options:     []Option{WithSorting(fields)},
			expectedErr: NewInvalidInputError("sort", "must match the regex: ^[a-zA-Z]+ (asc|desc)(, [a-zA-Z]+ (asc|desc))*$"),
		},
		{
			name:        "invalid sort field",
			url:         "?sort=notReal desc",
			options:     []Option{WithSorting(fields)},
			expectedErr: NewInvalidInputError("sort", "notReal is not a sortable field"),
		},
		{
			name:        "invalid filter field",
			url:         "?notReal[eq]=Squid",
			options:     []Option{WithFilters(fields)},
			expectedErr: NewInvalidInputError("notReal", "not a filterable field"),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			testURL, err := url.Parse(test.url)
			require.NoError(t, err)
			req := &http.Request{URL: testURL}

			result, err := Load(req, test.options...)

			require.Equal(t, test.expectedErr, err)
			require.Equal(t, test.expectedResult, result)
		})
	}
}
