package pagination

import (
	"testing"

	"api5back/src/database"
	"api5back/src/model"

	"github.com/stretchr/testify/require"
)

func TestParsePageRequest(t *testing.T) {
	for _, testCase := range []database.TestCase{
		{
			Name: "valid page and pageSize",
			Run: func(t *testing.T) {
				pageRequest := &model.PageRequest{
					Page:     &[]int{3}[0],
					PageSize: &[]int{20}[0],
				}

				page, pageSize, err := ParsePageRequest(pageRequest)
				require.NoError(t, err)
				require.Equal(t, 3, page)
				require.Equal(t, 20, pageSize)
			},
		},
		{
			Name: "nil Page and PageSize should fallback to defaults",
			Run: func(t *testing.T) {
				pageRequest := &model.PageRequest{
					Page:     nil,
					PageSize: nil,
				}

				page, pageSize, err := ParsePageRequest(pageRequest)
				require.NoError(t, err)
				require.Equal(t, DefaultPage, page)
				require.Equal(t, DefaultPageSize, pageSize)
			},
		},
		{
			Name: "nil PageRequest should fallback to defaults",
			Run: func(t *testing.T) {
				pageRequest := (*model.PageRequest)(nil)

				page, pageSize, err := ParsePageRequest(pageRequest)
				require.NoError(t, err)
				require.Equal(t, DefaultPage, page)
				require.Equal(t, DefaultPageSize, pageSize)
			},
		},
		{
			Name: "valid page and default pageSize",
			Run: func(t *testing.T) {
				pageRequest := &model.PageRequest{
					Page:     &[]int{3}[0],
					PageSize: nil,
				}

				page, pageSize, err := ParsePageRequest(pageRequest)
				require.NoError(t, err)
				require.Equal(t, 3, page)
				require.Equal(t, DefaultPageSize, pageSize)
			},
		},
		{
			Name: "default page and valid pageSize",
			Run: func(t *testing.T) {
				pageRequest := &model.PageRequest{
					Page:     nil,
					PageSize: &[]int{20}[0],
				}

				page, pageSize, err := ParsePageRequest(pageRequest)
				require.NoError(t, err)
				require.Equal(t, DefaultPage, page)
				require.Equal(t, 20, pageSize)
			},
		},
		{
			Name: "negative page",
			Run: func(t *testing.T) {
				pageRequest := &model.PageRequest{
					Page:     &[]int{-1}[0],
					PageSize: nil,
				}

				page, pageSize, err := ParsePageRequest(pageRequest)
				require.Error(t, err)
				require.Equal(t, 0, page)
				require.Equal(t, 0, pageSize)
			},
		},
	} {
		if testResult := t.Run(testCase.Name, testCase.Run); !testResult {
			t.Fatalf("Test case failed")

func TestParseOffsetAndTotalPages(t *testing.T) {
	type TestCase struct {
		Name               string
		Page               int
		PageSize           int
		TotalItems         int
		ExpectedOffset     int
		ExpectedTotalPages int
	}

	for i, testCase := range []TestCase{
		{
			Name: "valid page and pageSize",
			Page: 3, PageSize: 20, TotalItems: 100,
			ExpectedOffset:     40,
			ExpectedTotalPages: 5,
		},
		{
			Name: "page exceeds total pages",
			Page: 6, PageSize: 20, TotalItems: 100,
			ExpectedOffset:     100,
			ExpectedTotalPages: 5,
		},
		{
			Name: "page equals total pages",
			Page: 5, PageSize: 20, TotalItems: 100,
			ExpectedOffset:     80,
			ExpectedTotalPages: 5,
		},
		{
			Name: "page and pageSize exceed total items",
			Page: 3, PageSize: 50, TotalItems: 100,
			ExpectedOffset:     100,
			ExpectedTotalPages: 2,
		},
		{
			Name: "page and pageSize equal total items",
			Page: 2, PageSize: 50, TotalItems: 100,
			ExpectedOffset:     50,
			ExpectedTotalPages: 2,
		},
		{
			Name: "page and pageSize less than total items",
			Page: 1, PageSize: 50, TotalItems: 100,
			ExpectedOffset:     0,
			ExpectedTotalPages: 2,
		},
		{
			Name: "page and pageSize less than total items",
			Page: 1, PageSize: 50, TotalItems: 100,
			ExpectedOffset:     0,
			ExpectedTotalPages: 2,
		},
		{
			Name: "total items exceeds page and pageSize by 1",
			Page: 1, PageSize: 10, TotalItems: 11,
			ExpectedOffset:     0,
			ExpectedTotalPages: 2,
		},
		{
			Name: "second page when total items exceeds page and pageSize by 1",
			Page: 2, PageSize: 10, TotalItems: 11,
			ExpectedOffset:     10,
			ExpectedTotalPages: 2,
		},
	} {
		if testResult := t.Run(testCase.Name, func(t *testing.T) {
			offset, totalPages := processing.ParseOffsetAndTotalPages(
				testCase.Page,
				testCase.PageSize,
				testCase.TotalItems,
			)
			require.Equalf(t,
				testCase.ExpectedOffset, offset,
				"Expected offset %d, got %d",
				testCase.ExpectedOffset, offset,
			)
			require.Equalf(t,
				testCase.ExpectedTotalPages, totalPages,
				"Expected total pages %d, got %d",
				testCase.ExpectedTotalPages, totalPages,
			)
		}); !testResult {
			t.Errorf("Test case %d failed", i)
		}
	}
}
