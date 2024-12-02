package pagination

import (
	"testing"

	"api5back/src/model"
	"api5back/src/processing"

	"github.com/stretchr/testify/require"
)

func TestParsePageRequest(t *testing.T) {
	for i, testCase := range []struct {
		Name             string
		PageRequest      *model.PageRequest
		ExpectedError    bool
		ExpectedPage     int
		ExpectedPageSize int
	}{
		{
			Name: "valid page and pageSize",
			PageRequest: &model.PageRequest{
				Page:     &[]int{3}[0],
				PageSize: &[]int{20}[0],
			},
			ExpectedError:    false,
			ExpectedPage:     3,
			ExpectedPageSize: 20,
		},
		{
			Name: "nil Page and PageSize should fallback to defaults",
			PageRequest: &model.PageRequest{
				Page:     nil,
				PageSize: nil,
			},
			ExpectedError:    false,
			ExpectedPage:     DefaultPage,
			ExpectedPageSize: DefaultPageSize,
		},
		{
			Name:             "nil PageRequest should fallback to defaults",
			PageRequest:      (*model.PageRequest)(nil),
			ExpectedError:    false,
			ExpectedPage:     DefaultPage,
			ExpectedPageSize: DefaultPageSize,
		},
		{
			Name: "valid page and default pageSize",
			PageRequest: &model.PageRequest{
				Page:     &[]int{3}[0],
				PageSize: nil,
			},
			ExpectedError:    false,
			ExpectedPage:     3,
			ExpectedPageSize: DefaultPageSize,
		},
		{
			Name: "default page and valid pageSize",
			PageRequest: &model.PageRequest{
				Page:     nil,
				PageSize: &[]int{20}[0],
			},
			ExpectedError:    false,
			ExpectedPage:     DefaultPage,
			ExpectedPageSize: 20,
		},
		{
			Name: "negative page",
			PageRequest: &model.PageRequest{
				Page:     &[]int{-1}[0],
				PageSize: nil,
			},
			ExpectedError:    true,
			ExpectedPage:     0,
			ExpectedPageSize: 0,
		},
	} {
		if testResult := t.Run(testCase.Name, func(t *testing.T) {
			page, pageSize, err := ParsePageRequest(testCase.PageRequest)
			if testCase.ExpectedError {
				require.Error(t, err, "Expected error, got nil")
			} else {
				require.NoError(t, err, "Expected no error, got %v")
				require.Equalf(t,
					testCase.ExpectedPage, page,
					"Expected page %d, got %d",
					testCase.ExpectedPage, page,
				)
				require.Equalf(t,
					testCase.ExpectedPageSize, pageSize,
					"Expected pageSize %d, got %d",
					testCase.ExpectedPageSize, pageSize,
				)
			}
		}); !testResult {
			t.Errorf("Test case %d failed", i)
		}
	}
}

func TestParseOffsetAndTotalPages(t *testing.T) {
	for i, testCase := range []struct {
		Name               string
		Page               int
		PageSize           int
		TotalItems         int
		ExpectedOffset     int
		ExpectedTotalPages int
	}{
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
