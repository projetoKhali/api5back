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
	switch v := value.(type) {
	case nil:
		return nil
	case string:
	case []byte:
		for i, str := range HiringProcessCandidateStatus(0).Values() {
			if str == string(v) {
				*s = HiringProcessCandidateStatus(i)
				break
			}
		}
		return nil
	default:
		return fmt.Errorf("invalid hiring_process_candidate status: %v", value)
	}

	return nil
}
