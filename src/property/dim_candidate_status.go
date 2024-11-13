package property

import (
	"database/sql/driver"
	"fmt"
)

type DimCandidateStatus int

const (
	DimCandidateStatusInAnalysis DimCandidateStatus = iota
	DimCandidateStatusInterview
	DimCandidateStatusHired
	DimCandidateStatusRejected
)

func (s DimCandidateStatus) String() string {
	return [...]string{
		"In Analysis",
		"Interview",
		"Hired",
		"Rejected",
	}[s]
}

func (DimCandidateStatus) Values() []string {
	return []string{
		DimCandidateStatusInAnalysis.String(),
		DimCandidateStatusInterview.String(),
		DimCandidateStatusHired.String(),
		DimCandidateStatusRejected.String(),
	}
}

func (s DimCandidateStatus) Value() (driver.Value, error) {
	return s.String(), nil
}

func (s *DimCandidateStatus) Scan(value interface{}) error {
	var valueStr string
	switch v := value.(type) {
	case nil:
		return nil
	case int:
		*s = DimCandidateStatus(v)
		return nil
	case string:
		valueStr = v
	case []byte:
		valueStr = string(v)
	default:
		return fmt.Errorf("invalid dim_candidate status: %v", value)
	}

	for i, statusStr := range DimCandidateStatus(0).Values() {
		if statusStr == string(valueStr) {
			*s = DimCandidateStatus(i)
			return nil
		}
	}

	return fmt.Errorf("invalid dim_candidate status: %q", value)
}
