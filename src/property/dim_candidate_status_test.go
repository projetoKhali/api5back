package property

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseDimCandidateStatus(t *testing.T) {
	for _, testCase := range []PropertyStatusTestCase{
		{
			Name:           "In Analysis",
			IntValue:       0,
			ExpectedPanic:  false,
			ExpectedError:  false,
			ExpectedStatus: "In Analysis",
		},
		{
			Name:           "Interview",
			IntValue:       1,
			ExpectedPanic:  false,
			ExpectedError:  false,
			ExpectedStatus: "Interview",
		},
		{
			Name:           "Hired",
			IntValue:       2,
			ExpectedPanic:  false,
			ExpectedError:  false,
			ExpectedStatus: "Hired",
		},
		{
			Name:           "Rejected",
			IntValue:       3,
			ExpectedPanic:  false,
			ExpectedError:  false,
			ExpectedStatus: "Rejected",
		},
		{
			Name:           "Invalid positive value",
			IntValue:       4,
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
			ExpectedPanic:  false,
			ExpectedError:  false,
			ExpectedStatus: "In Analysis",
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
				status := HiringProcessCandidateStatus(0)
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
