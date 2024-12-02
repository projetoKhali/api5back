package property

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseDimProcessStatus(t *testing.T) {
	for _, testCase := range []PropertyStatusTestCase{
		{
			Name:           "Open",
			IntValue:       1,
			ExpectedPanic:  false,
			ExpectedError:  false,
			ExpectedStatus: "Open",
		},
		{
			Name:           "In Progress",
			IntValue:       2,
			ExpectedPanic:  false,
			ExpectedError:  false,
			ExpectedStatus: "In Progress",
		},
		{
			Name:           "Closed",
			IntValue:       3,
			ExpectedPanic:  false,
			ExpectedError:  false,
			ExpectedStatus: "Closed",
		},
		{
			Name:           "Invalid positive value",
			IntValue:       4,
			ExpectedPanic:  true,
			ExpectedError:  false,
			ExpectedStatus: "",
		},
		{
			Name:           "Invalid zero value",
			IntValue:       0,
			ExpectedPanic:  true,
			ExpectedError:  false,
			ExpectedStatus: "",
		},
		{
			Name:           "Invalid negative value",
			IntValue:       -1,
			ExpectedPanic:  true,
			ExpectedError:  false,
			ExpectedStatus: "",
		},
		{
			Name:           "Invalid string value",
			IntValue:       "Invalid",
			ExpectedPanic:  false,
			ExpectedError:  true,
			ExpectedStatus: "",
		},
		{
			Name:           "Nil value",
			IntValue:       nil,
			ExpectedPanic:  true,
			ExpectedError:  false,
			ExpectedStatus: "",
		},
		{
			Name:           "Empty string value",
			IntValue:       "",
			ExpectedPanic:  false,
			ExpectedError:  true,
			ExpectedStatus: "",
		},
	} {
		t.Run(testCase.Name, func(t *testing.T) {
			testFunction := func() {
				status := DimProcessStatus(0)
				err := status.Scan(testCase.IntValue)
				if err != nil {
					require.Equal(t, testCase.ExpectedStatus, "")
				} else {
					require.Equal(t, testCase.ExpectedStatus, status.String())
				}
			}

			if testCase.ExpectedPanic {
				require.Panics(t, testFunction)
			} else {
				require.NotPanics(t, testFunction)
			}
		})
	}
}
