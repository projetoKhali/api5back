package property

import (
	"database/sql/driver"
	"fmt"
)

type DimVacancyStatus int

const (
	DimVacancyStatusOpen DimVacancyStatus = iota + 1
	DimVacancyStatusInAnalysis
	DimVacancyStatusClosed
)

func (s DimVacancyStatus) String() string {
	return [...]string{
		"Open",
		"In Analysis",
		"Closed",
	}[s-1]
}

func (DimVacancyStatus) Values() []string {
	return []string{
		DimVacancyStatusOpen.String(),
		DimVacancyStatusInAnalysis.String(),
		DimVacancyStatusClosed.String(),
	}
}

func (s DimVacancyStatus) Value() (driver.Value, error) {
	return s.String(), nil
}

func (s *DimVacancyStatus) Scan(value interface{}) error {
	var valueStr string
	switch v := value.(type) {
	case nil:
		return nil
	case int:
		*s = DimVacancyStatus(v)
		return nil
	case string:
		valueStr = v
	case []byte:
		valueStr = string(v)
	default:
		return fmt.Errorf("invalid dim_vacancy status: %v", value)
	}

	for i, statusStr := range DimVacancyStatus(0).Values() {
		if statusStr == string(valueStr) {
			*s = DimVacancyStatus(i)
			return nil
		}
	}

	return fmt.Errorf("invalid dim_vacancy status: %q", value)
}
