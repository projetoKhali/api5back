package property

import (
	"database/sql/driver"
	"fmt"
)

type DimProcessStatus int

const (
	DimProcessStatusOpen DimProcessStatus = iota + 1
	DimProcessStatusInProgress
	DimProcessStatusClosed
)

func (s DimProcessStatus) String() string {
	return [...]string{
		"Open",
		"In Progress",
		"Closed",
	}[s-1]
}

func (DimProcessStatus) Values() []string {
	return []string{
		DimProcessStatusOpen.String(),
		DimProcessStatusInProgress.String(),
		DimProcessStatusClosed.String(),
	}
}

func (s DimProcessStatus) Value() (driver.Value, error) {
	return s.String(), nil
}

func (s *DimProcessStatus) Scan(value interface{}) error {
	var valueStr string
	switch v := value.(type) {
	case nil:
		return nil
	case int:
		*s = DimProcessStatus(v)
		return nil
	case string:
		valueStr = v
	case []byte:
		valueStr = string(v)
	default:
		return fmt.Errorf("invalid dim_process status: %v", value)
	}

	for i, statusStr := range DimProcessStatus(0).Values() {
		if statusStr == string(valueStr) {
			*s = DimProcessStatus(i)
			return nil
		}
	}

	return fmt.Errorf("invalid dim_process status: %q", value)
}
