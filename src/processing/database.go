package processing

import (
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

func ParseStringToPgtypeDate(
	layout string,
	dateString string,
) (pgDate pgtype.Date, err error) {
	t, err := time.Parse(layout, dateString)
	if err != nil {
		return pgDate, fmt.Errorf("failed to parse date: %v", err)
	}

	return pgtype.Date{
		Time:  t,
		Valid: true,
	}, nil
}

func ParsePageAndPageSize(
	page, pageSize *int,
) (int, int, error) {
	if page == nil {
		defaultPage := 1
		page = &defaultPage
	}
	if pageSize == nil {
		defaultPageSize := 10
		pageSize = &defaultPageSize
	}

	if *page <= 0 {
		return 0, 0, errors.New(
			"invalid page number",
		)
	}

	if *pageSize <= 0 {
		return 0, 0, errors.New(
			"invalid page size",
		)
	}

	return *page, *pageSize, nil
}

func ParseOffsetAndTotalPages(
	page, pageSize, totalRecords int,
) (int, int) {
	offset := (page - 1) * pageSize
	numMaxPages := (totalRecords + pageSize - 1) / pageSize
	return offset, numMaxPages
}
