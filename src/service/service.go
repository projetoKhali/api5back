package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"api5back/ent"
	"api5back/ent/dimprocess"
	"api5back/ent/dimvacancy"
	"api5back/ent/facthiringprocess"
	"api5back/src/processing"

	"github.com/jackc/pgx/v5/pgtype"
)

type MetricsService struct {
	dbClient *ent.Client
}

type MetricsData struct {
	VacancySummary    processing.VacancyStatusSummary      `json:"vacancyStatus"`
	CardInfos         processing.CardInfos                 `json:"cards"`
	AverageHiringTime processing.AverageHiringTimePerMonth `json:"averageHiringTime"`
}

type GetMetricsFilter struct {
	HiringProcessName string `json:"hiringProcess"`
	VacancyName       string `json:"vacancy"`
	StartDate         string `json:"startDate"`
	EndDate           string `json:"endDate"`
}

type DateRange struct {
	StartDate string `json:"startDate" form:"startDate" time_format:"2024-10-01T00:00:00Z"`
	EndDate   string `json:"endDate" form:"endDate" time_format:"2024-10-01T00:00:00Z"`
}

type VacancyTableFilter struct {
	Recruiters    []int      `json:"recruiters"`
	Processes     []int      `json:"processes"`
	Vacancies     []int      `json:"vacancies"`
	DateRange     *DateRange `json:"dateRange"`
	ProcessStatus []int      `json:"processStatus"`
	VacancyStatus []int      `json:"vacancyStatus"`
	Page          *int       `json:"page"`
	PageSize      *int       `json:"pageSize"`
}

type FactHiringProcessReturn struct {
	FactHiringProcess []*ent.FactHiringProcess `json:"factHiringProcess"`
	NumMaxPages       int                      `json:"numMaxPages"`
}

func NewMetricsService(dbclient *ent.Client) *MetricsService {
	return &MetricsService{dbClient: dbclient}
}

func (s *MetricsService) GetMetrics(
	ctx context.Context,
	filter GetMetricsFilter,
) (*MetricsData, error) {
	query := s.dbClient.
		FactHiringProcess.
		Query().
		WithDimVacancy().
		WithDimProcess().
		WithHiringProcessCandidates()

	if filter.HiringProcessName != "" {
		query = query.Where(
			facthiringprocess.HasDimProcessWith(
				dimprocess.TitleContains(
					filter.HiringProcessName,
				),
			),
		)
	}

	if filter.VacancyName != "" {
		query = query.Where(
			facthiringprocess.HasDimVacancyWith(
				dimvacancy.TitleContains(
					filter.VacancyName,
				),
			),
		)
	}

	if filter.StartDate != "" {
		hiringProcessStartDate, err := ParseStringToPgtypeDate(
			"2006-01-02",
			filter.StartDate,
		)
		if err != nil {
			return nil, fmt.Errorf(
				"could not parse `StartDate`: %w",
				err,
			)
		}

		query = query.Where(
			facthiringprocess.HasDimVacancyWith(
				dimvacancy.ClosingDateGTE(
					&hiringProcessStartDate,
				),
			),
		)
	}

	if filter.EndDate != "" {
		hiringProcessEndDate, err := ParseStringToPgtypeDate(
			"2006-01-02",
			filter.EndDate,
		)
		if err != nil {
			return nil, fmt.Errorf(
				"could not parse `EndDate`: %w",
				err,
			)
		}

		query = query.Where(
			facthiringprocess.HasDimVacancyWith(
				dimvacancy.OpeningDateLTE(
					&hiringProcessEndDate,
				),
			),
		)
	}

	hiringProcess, err := query.All(ctx)
	if err != nil {
		return nil, fmt.Errorf(
			"could not retrieve `FactHiringProcess` data: %w",
			err,
		)
	}

	var errors []error

	cardInfo, err := processing.ComputingCardInfo(hiringProcess)
	if err != nil {
		errors = append(errors, fmt.Errorf(
			"could not calculate `CardInfo` data: %w",
			err,
		))
	}

	vacancyInfo, err := processing.GenerateVacancyStatusSummary(hiringProcess)
	if err != nil {
		errors = append(errors, fmt.Errorf(
			"could not generate `VacancyStatus` summary: %w",
			err,
		))
	}

	averageHiringTime, err := processing.GenerateAverageHiringTimePerMonth(hiringProcess)
	if err != nil {
		errors = append(errors, fmt.Errorf(
			"could not generate `AvgHiringTime` data: %w",
			err,
		))
	}

	if len(errors) > 0 {
		return nil, fmt.Errorf(
			"failed to get metrics: %v",
			errors,
		)
	}

	return &MetricsData{
		CardInfos:         cardInfo,
		VacancySummary:    vacancyInfo,
		AverageHiringTime: averageHiringTime,
	}, nil
}

func ParseStringToPgtypeDate(
	layout string,
	dateString string,
) (pgDate pgtype.Date, err error) {
	t, err := time.Parse(layout, dateString)
	if err != nil {
		return pgDate, fmt.Errorf("failed to parse date: %v", err)
	}

	return pgtype.Date{
		Time:  t,
		Valid: true,
	}, nil
}

type VacancyServiceTable struct {
	dwClient *ent.Client
}

func NewVacancyServiceTable(client *ent.Client) *VacancyServiceTable {
	return &VacancyServiceTable{dwClient: client}
}

func (vs *VacancyServiceTable) GetVacancyTable(
	ctx context.Context,
	filter VacancyTableFilter,
) (*FactHiringProcessReturn, error) {
	query := vs.dwClient.FactHiringProcess.Query().WithDimProcess().WithDimVacancy()

	if len(filter.Recruiters) > 0 {
		query = query.Where(
			facthiringprocess.DimUserIdIn(filter.Recruiters...),
		)
	}

	if len(filter.Processes) > 0 {
		query = query.Where(
			facthiringprocess.DimProcessIdIn(filter.Processes...),
		)
	}

	if len(filter.Vacancies) > 0 {
		query = query.Where(
			facthiringprocess.DimVacancyIdIn(filter.Vacancies...),
		)
	}

	if filter.DateRange != nil {

		if filter.DateRange.StartDate != "" {

			startDateString := filter.DateRange.StartDate
			hiringProcessStartDate, err := ParseStringToPgtypeDate(
				"2006-01-02",
				startDateString,
			)
			if err != nil {
				return nil, fmt.Errorf(
					"could not parse `StartDate`: %w",
					err,
				)
			}

			query = query.Where(
				facthiringprocess.HasDimVacancyWith(
					dimvacancy.ClosingDateGTE(
						&hiringProcessStartDate,
					),
				),
			)
		}

		if filter.DateRange.EndDate != "" {

			endDateString := filter.DateRange.EndDate
			hiringProcessEndDate, err := ParseStringToPgtypeDate(
				"2006-01-02",
				endDateString,
			)
			if err != nil {
				return nil, fmt.Errorf(
					"could not parse `EndDate`: %w",
					err,
				)
			}

			query = query.Where(
				facthiringprocess.HasDimVacancyWith(
					dimvacancy.OpeningDateLTE(
						&hiringProcessEndDate,
					),
				),
			)
		}

	}

	if len(filter.ProcessStatus) > 0 {
		query = query.Where(
			facthiringprocess.HasDimProcessWith(
				dimprocess.StatusIn(filter.ProcessStatus...),
			),
		)
	}

	if len(filter.VacancyStatus) > 0 {
		query = query.Where(
			facthiringprocess.HasDimVacancyWith(
				dimvacancy.StatusIn(filter.VacancyStatus...),
			),
		)
	}

	if filter.Page == nil {
		defaultPage := 1
		filter.Page = &defaultPage
	}
	if filter.PageSize == nil {
		defaultPageSize := 10
		filter.PageSize = &defaultPageSize
	}
	if *filter.Page <= 0 || *filter.PageSize <= 0 {
		return nil, errors.New(
			"invalid page number or size",
		)
	}

	totalRecords, err := query.Count(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not get total records: %w", err)
	}

	numMaxPages := (totalRecords + *filter.PageSize - 1) / *filter.PageSize

	offset := (*filter.Page - 1) * *filter.PageSize
	query = query.Offset(offset).Limit(*filter.PageSize)

	vacancies, err := query.All(ctx)

	if err != nil {
		return nil, err
	}

	return &FactHiringProcessReturn{
		FactHiringProcess: vacancies,
		NumMaxPages:       numMaxPages,
	}, nil
}
