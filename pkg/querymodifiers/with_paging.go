package querymodifiers

type Page struct {
	Size             int
	Number           int
	IncludeTotalSize bool
}

func WithPaging(defaultSize int) func(*URLParams, *QueryModifiers) error {
	return func(urlParams *URLParams, reqParams *QueryModifiers) error {
		reqParams.Page = &Page{}

		var err error
		if reqParams.Page.IncludeTotalSize, err = urlParams.GetBool("includeTotalSize"); err != nil {
			return err
		}

		if reqParams.Page.Size, err = urlParams.GetInt("pageSize", defaultSize); err != nil {
			return err
		}

		reqParams.Page.Number, err = urlParams.GetInt("pageNumber", 1)
		return err
	}
}
