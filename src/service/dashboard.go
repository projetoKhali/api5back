package service

import (
	"context"
	"fmt"

	"api5back/ent"
	"api5back/ent/dimprocess"
	"api5back/ent/dimuser"
	"api5back/ent/dimvacancy"
	"api5back/ent/facthiringprocess"
	"api5back/src/model"
	"api5back/src/processing"
	"api5back/src/property"
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
			hiringProcessStartDate, err := processing.ParseStringToPgtypeDate(
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
			hiringProcessEndDate, err := processing.ParseStringToPgtypeDate(
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
			WithDimProcess().
			WithDimVacancy(func(query *ent.DimVacancyQuery) {
				query.WithHiringProcessCandidates()
			}),
		filter,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"could not apply filters: %w",
			err,
		)
	}

	hiringProcesses, err := query.All(ctx)
	if err != nil {
		return nil, fmt.Errorf(
			"could not retrieve `FactHiringProcess` data: %w",
			err,
		)
	}

	var errors []error

	cardInfo, err := processing.ComputingCardsInfo(hiringProcesses)
	if err != nil {
		errors = append(errors, fmt.Errorf(
			"could not calculate `CardInfo` data: %w",
			err,
		))
	}

	vacancyInfo, err := processing.GenerateVacancyStatusSummary(hiringProcesses)
	if err != nil {
		errors = append(errors, fmt.Errorf(
			"could not generate `VacancyStatus` summary: %w",
			err,
		))
	}

	var dimVacancies []*ent.DimVacancy
	for _, hp := range hiringProcesses {
		if hp.Edges.DimVacancy != nil {
			dimVacancies = append(dimVacancies, hp.Edges.DimVacancy)
		}
	}

	averageHiringTime, err := processing.GenerateAverageHiringTimePerMonth(dimVacancies)
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

func GetVacancyTable(
	ctx context.Context,
	client *ent.Client,
	filter model.FactHiringProcessFilter,
) (*model.Page[model.DashboardTableRow], error) {
	query, err := applyQueryFilters(
		client.
			FactHiringProcess.
			Query().
			WithDimProcess().
			WithDimVacancy(func(query *ent.DimVacancyQuery) {
				query.WithHiringProcessCandidates()
			}),
		filter,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"could not apply filters: %w",
			err,
		)
	}

	page, pageSize, err := model.ParsePageRequest(filter)
	if err != nil {
		return nil, err
	}

	totalRecords, err := query.Count(ctx)
	if err != nil {
		return nil, err
	}

	offset, numMaxPages := processing.ParseOffsetAndTotalPages(
		page,
		pageSize,
		totalRecords,
	)

	factHiringProcesses, err := query.
		Offset(offset).
		Limit(pageSize).
		All(ctx)
	if err != nil {
		return nil, err
	}

	var tableDatas []model.DashboardTableRow
	for _, factHiringProcess := range factHiringProcesses {

		numPositions := factHiringProcess.Edges.DimVacancy.NumPositions
		var competitionRate *float32
		if numPositions > 0 {
			rate := float32(factHiringProcess.MetTotalCandidatesApplied) / float32(numPositions)
			competitionRate = &rate
		} else {
			competitionRate = nil
		}

		hiringTime, err := processing.GenerateAverageHiringTimePerFactHiringProcess(factHiringProcess)
		var averageHiringTime *float32
		if err != nil {
			averageHiringTime = nil
		} else {
			averageHiringTime = &(hiringTime)
		}

		numFeedback := factHiringProcess.MetTotalFeedbackPositive + factHiringProcess.MetTotalNegative + factHiringProcess.MetTotalNeutral
		tableDatas = append(tableDatas, model.DashboardTableRow{
			ProcessTitle:      factHiringProcess.Edges.DimProcess.Title,
			VacancyTitle:      factHiringProcess.Edges.DimVacancy.Title,
			NumPositions:      numPositions,
			NumCandidates:     factHiringProcess.MetTotalCandidatesApplied,
			CompetitionRate:   competitionRate,
			NumInterviewed:    factHiringProcess.MetTotalCandidatesInterviewed,
			NumHired:          factHiringProcess.MetTotalCandidatesHired,
			AverageHiringTime: averageHiringTime,
			NumFeedback:       numFeedback,
		})
	}

	return &model.Page[model.DashboardTableRow]{
		Items:       tableDatas,
		NumMaxPages: numMaxPages,
	}, nil
}
