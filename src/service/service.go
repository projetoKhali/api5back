package service

import (
	"context"
	"fmt"
	"time"

	"api5back/ent"
	"api5back/ent/dimprocess"
	"api5back/ent/dimuser"
	"api5back/ent/dimvacancy"
	"api5back/ent/facthiringprocess"
	"api5back/src/model"
	"api5back/src/processing"

	"github.com/jackc/pgx/v5/pgtype"
)

type DateRange struct {
	StartDate string `json:"startDate" form:"startDate" time_format:"2024-10-01T00:00:00Z"`
	EndDate   string `json:"endDate" form:"endDate" time_format:"2024-10-01T00:00:00Z"`
}

type FactHiringProcessFilter struct {
	Recruiters    []int      `json:"recruiters"`
	Processes     []int      `json:"processes"`
	Vacancies     []int      `json:"vacancies"`
	DateRange     *DateRange `json:"dateRange"`
	ProcessStatus []int      `json:"processStatus"`
	VacancyStatus []int      `json:"vacancyStatus"`
}

func GetMetrics(
	ctx context.Context,
	client *ent.Client,
	filter FactHiringProcessFilter,
) (*model.MetricsData, error) {
	query := client.
		FactHiringProcess.
		Query().
		WithDimVacancy().
		WithDimProcess().
		WithHiringProcessCandidates()

	if len(filter.Processes) > 0 {
		query = query.Where(
			facthiringprocess.HasDimProcessWith(
				dimprocess.IDIn(filter.Processes...),
			),
		)
	}

	if len(filter.Vacancies) > 0 {
		query = query.Where(
			facthiringprocess.HasDimVacancyWith(
				dimvacancy.IDIn(filter.Vacancies...),
			),
		)
	}

	if len(filter.Recruiters) > 0 {
		query = query.Where(
			facthiringprocess.HasDimUserWith(
				dimuser.IDIn(filter.Recruiters...),
			),
		)
	}

	if filter.DateRange != nil {
		if filter.DateRange.StartDate != "" {
			hiringProcessStartDate, err := ParseStringToPgtypeDate(
				"2006-01-02",
				filter.DateRange.StartDate,
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
			hiringProcessEndDate, err := ParseStringToPgtypeDate(
				"2006-01-02",
				filter.DateRange.EndDate,
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

	return &model.MetricsData{
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

func GetVacancyTable(
	ctx context.Context,
	client *ent.Client,
	filter FactHiringProcessFilter,
) ([]model.TableData, error) {
	query := client.FactHiringProcess.Query().WithDimProcess().WithDimVacancy()

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

	vacancies, err := query.All(ctx)
	if err != nil {
		return nil, err
	}

	var tableDatas []model.TableData
	for _, vacancy := range vacancies {

		numPositions := vacancy.Edges.DimVacancy.NumPositions
		var competitionRate *float32
		if numPositions > 0 {
			rate := float32(vacancy.MetTotalCandidatesApplied) / float32(numPositions)
			competitionRate = &rate
		} else {
			competitionRate = nil
		}

		hiringTime, err := processing.GenerateAverageHiringTimePerFactHiringProcess(vacancy)
		var averageHiringTime *float32
		if err != nil {
			averageHiringTime = nil
		} else {
			averageHiringTime = &(hiringTime)
		}

		numFeedback := vacancy.MetTotalFeedbackPositive + vacancy.MetTotalNegative + vacancy.MetTotalNeutral
		tableDatas = append(tableDatas, model.TableData{
			Title:             vacancy.Edges.DimVacancy.Title,
			NumPositions:      numPositions,
			NumCandidates:     vacancy.MetTotalCandidatesApplied,
			CompetitionRate:   competitionRate,
			NumInterviewed:    vacancy.MetTotalCandidatesInterviewed,
			NumHired:          vacancy.MetTotalCandidatesHired,
			AverageHiringTime: averageHiringTime,
			NumFeedback:       numFeedback,
		})
	}

	return tableDatas, nil
}
