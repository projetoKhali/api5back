package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"api5back/ent"
	"api5back/ent/dimprocess"
	"api5back/ent/dimuser"
	"api5back/ent/dimvacancy"
	"api5back/ent/facthiringprocess"
	"api5back/src/model"
	"api5back/src/processing"
	"api5back/src/property"

	"github.com/jackc/pgx/v5/pgtype"
)

func applyQueryFilters(
	query *ent.FactHiringProcessQuery,
	filter model.FactHiringProcessFilter,
) (*ent.FactHiringProcessQuery, error) {
	if filter.Recruiters != nil && len(filter.Recruiters) > 0 {
		query = query.Where(
			facthiringprocess.HasDimUserWith(
				dimuser.IDIn(filter.Recruiters...),
			),
		)
	}
	if filter.Processes != nil && len(filter.Processes) > 0 {
		query = query.Where(
			facthiringprocess.HasDimProcessWith(
				dimprocess.IDIn(filter.Processes...),
			),
		)
	}

	if filter.Vacancies != nil && len(filter.Vacancies) > 0 {
		query = query.Where(
			facthiringprocess.HasDimVacancyWith(
				dimvacancy.IDIn(filter.Vacancies...),
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

	if filter.ProcessStatus != nil && len(filter.ProcessStatus) > 0 {
		var processStatuses []property.DimProcessStatus
		for _, status := range filter.ProcessStatus {
			processStatuses = append(processStatuses, property.DimProcessStatus(status))
		}

		query = query.Where(
			facthiringprocess.HasDimProcessWith(
				dimprocess.StatusIn(processStatuses...),
			),
		)
	}

	if filter.VacancyStatus != nil && len(filter.VacancyStatus) > 0 {
		var vacancyStatuses []property.DimVacancyStatus
		for _, status := range filter.VacancyStatus {
			vacancyStatuses = append(vacancyStatuses, property.DimVacancyStatus(status))
		}

		query = query.Where(
			facthiringprocess.HasDimVacancyWith(
				dimvacancy.StatusIn(vacancyStatuses...),
			),
		)
	}

	return query, nil
}

func GetMetrics(
	ctx context.Context,
	client *ent.Client,
	filter model.FactHiringProcessFilter,
) (*model.DashboardMetrics, error) {
	query, err := applyQueryFilters(
		client.
			FactHiringProcess.
			Query().
			WithDimVacancy().
			WithDimProcess().
			WithHiringProcessCandidates(),
		filter,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"could not apply filters: %w",
			err,
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

	cardInfo, err := processing.ComputingCardsInfo(hiringProcess)
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

	return &model.DashboardMetrics{
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

func GetVacancyTable(
	ctx context.Context,
	client *ent.Client,
	filter model.FactHiringProcessFilter,
) (*model.DashboardTablePage, error) {
	query, err := applyQueryFilters(
		client.
			FactHiringProcess.
			Query().
			WithDimProcess().
			WithDimVacancy().
			WithHiringProcessCandidates(),
		filter,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"could not apply filters: %w",
			err,
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

	var tableDatas []model.DashboardTableRow
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
		tableDatas = append(tableDatas, model.DashboardTableRow{
			ProcessTitle:      vacancy.Edges.DimProcess.Title,
			VacancyTitle:      vacancy.Edges.DimVacancy.Title,
			NumPositions:      numPositions,
			NumCandidates:     vacancy.MetTotalCandidatesApplied,
			CompetitionRate:   competitionRate,
			NumInterviewed:    vacancy.MetTotalCandidatesInterviewed,
			NumHired:          vacancy.MetTotalCandidatesHired,
			AverageHiringTime: averageHiringTime,
			NumFeedback:       numFeedback,
		})
	}

	return &model.DashboardTablePage{
		Items:       tableDatas,
		NumMaxPages: numMaxPages,
	}, nil
}
