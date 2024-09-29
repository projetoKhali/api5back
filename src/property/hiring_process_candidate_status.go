package property

import (
	"database/sql/driver"
	"fmt"
)

type HiringProcessCandidateStatus int

const (
	HiringProcessCandidateStatusInAnalysis HiringProcessCandidateStatus = iota
	HiringProcessCandidateStatusInterview
	HiringProcessCandidateStatusHired
	HiringProcessCandidateStatusRejected
)

func (s HiringProcessCandidateStatus) String() string {
	return [...]string{
		"In Analysis",
		"Interview",
		"Hired",
		"Rejected",
	}[s]
}

func (HiringProcessCandidateStatus) Values() []string {
	return []string{
		HiringProcessCandidateStatusInAnalysis.String(),
		HiringProcessCandidateStatusInterview.String(),
		HiringProcessCandidateStatusHired.String(),
		HiringProcessCandidateStatusRejected.String(),
	}
}

func (s HiringProcessCandidateStatus) Value() (driver.Value, error) {
	return s.String(), nil
}

func (s *HiringProcessCandidateStatus) Scan(value interface{}) error {
	var valueStr string
	switch v := value.(type) {
	case nil:
		return nil
	case int:
		*s = HiringProcessCandidateStatus(v)
		return nil
	case string:
		valueStr = v
	case []byte:
		valueStr = string(v)
	default:
		return fmt.Errorf("invalid hiring_process_candidate status: %v", value)
	}

	for i, statusStr := range HiringProcessCandidateStatus(0).Values() {
		if statusStr == string(valueStr) {
			*s = HiringProcessCandidateStatus(i)
			return nil
		}
	}

	return fmt.Errorf("invalid hiring_process_candidate status: %q", value)
}
