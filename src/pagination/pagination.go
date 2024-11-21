package pagination

import (
	"errors"

	"api5back/src/model"
)

func parsePageRequestError(err string) (int, int, error) {
	return 0, 0, errors.New(err)
}

var (
	DefaultPage     = 1
	DefaultPageSize = 10
)

func ParsePageRequest(pageRequest model.PageRequester) (int, int, error) {
	if pageRequest == nil {
		return DefaultPage, DefaultPageSize, nil
	}

	pr := (pageRequest).GetPageRequest()
	if pr == nil {
		return DefaultPage, DefaultPageSize, nil
	}

	page, pageSize := pr.Page, pr.PageSize

	if page == nil {
		page = &DefaultPage
	}
	if pageSize == nil {
		pageSize = &DefaultPageSize
	}

	if *page <= 0 {
		return parsePageRequestError("invalid page number")
	}

	if *pageSize <= 0 {
		return parsePageRequestError("invalid page size")
	}

	return *page, *pageSize, nil
}
