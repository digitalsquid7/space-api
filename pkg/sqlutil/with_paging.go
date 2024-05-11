package sqlutil

import (
	"space-api/pkg/models"
	"space-api/pkg/urlparams"
)

type Page struct {
	Size             int
	Number           int
	IncludeTotalSize bool
}

func WithPaging(defaultSize int) func(*urlparams.URLParams, *QueryModifiers) error {
	return func(urlParams *urlparams.URLParams, reqParams *QueryModifiers) error {
		reqParams.Page = &Page{}

		var err error
		if reqParams.Page.IncludeTotalSize, err = urlParams.GetBool("includeTotalSize"); err != nil {
			return err
		}

		if reqParams.Page.Size, err = urlParams.GetInt("pageSize", defaultSize); err != nil {
			return err
		}

		if reqParams.Page.Size < 1 || reqParams.Page.Size > 1000 {
			return models.NewInvalidInputError("pageSize", "must be >= 1 and <= 1000")
		}

		if reqParams.Page.Number, err = urlParams.GetInt("pageNumber", 1); err != nil {
			return err
		}

		if reqParams.Page.Number < 1 {
			return models.NewInvalidInputError("pageNumber", "must be >= 1")
		}

		return err
	}
}
