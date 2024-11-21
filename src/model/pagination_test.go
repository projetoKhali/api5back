package model

import (
	"testing"

	"api5back/src/database"

	"github.com/stretchr/testify/require"
)

func TestParsePageRequest(t *testing.T) {
	for _, testCase := range []database.TestCase{
		{
			Name: "valid page and pageSize",
			Run: func(t *testing.T) {
				pageRequest := &PageRequest{
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
				pageRequest := &PageRequest{
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
				pageRequest := (*PageRequest)(nil)

				page, pageSize, err := ParsePageRequest(pageRequest)
				require.NoError(t, err)
				require.Equal(t, DefaultPage, page)
				require.Equal(t, DefaultPageSize, pageSize)
			},
		},
		{
			Name: "valid page and default pageSize",
			Run: func(t *testing.T) {
				pageRequest := &PageRequest{
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
				pageRequest := &PageRequest{
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
				pageRequest := &PageRequest{
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
		}
	}
}
