package seeds

import (
	"context"
	"fmt"
	"time"

	"api5back/ent"

	"github.com/jackc/pgx/v5/pgtype"
)

// Função para popular os dados no banco
func DataWarehouse(client *ent.Client) error {
	ctx := context.Background()

	users := []ent.DimUser{
		{DbId: 1, Name: "Alice Santos", Occupation: "Recruiter"},
		{DbId: 2, Name: "Bob Ferreira", Occupation: "HR Manager"},
		{DbId: 3, Name: "Carla Mendes", Occupation: "Software Engineer"},
		{DbId: 4, Name: "David Costa", Occupation: "Data Analyst"},
		{DbId: 5, Name: "Eva Lima", Occupation: "Product Manager"},
	}

	var userIDs []int64

	for _, user := range users {
		createdUser, err := client.DimUser.Create().
			SetDbId(user.DbId).
			SetName(user.Name).
			SetOccupation(user.Occupation).
			Save(ctx)
		if err != nil {
			return fmt.Errorf("failed to create user %s: %v", user.Name, err)
		}
		userIDs = append(userIDs, int64(createdUser.DbId))
	}

	// Inserindo datas na tabela dim_datetime
	dates := []ent.DimDatetime{
		{
			Date: &pgtype.Date{Time: time.Date(2024, 7, 1, 0, 0, 0, 0, time.UTC), Valid: true},
			Year: 2024, Month: 7, Weekday: 1, Day: 0, Hour: 0, Minute: 0, Second: 0,
		},
		{
			Date: &pgtype.Date{Time: time.Date(2024, 8, 1, 0, 0, 0, 0, time.UTC), Valid: true},
			Year: 2024, Month: 8, Weekday: 1, Day: 0, Hour: 0, Minute: 0, Second: 0,
		},
		{
			Date: &pgtype.Date{Time: time.Date(2024, 9, 1, 0, 0, 0, 0, time.UTC), Valid: true},
			Year: 2024, Month: 9, Weekday: 1, Day: 0, Hour: 0, Minute: 0, Second: 0,
		},
		{
			Date: &pgtype.Date{Time: time.Date(2024, 9, 1, 0, 0, 0, 0, time.UTC), Valid: true},
			Year: 2024, Month: 9, Weekday: 2, Day: 0, Hour: 0, Minute: 0, Second: 0,
		},
		{
			Date: &pgtype.Date{Time: time.Date(2024, 9, 1, 0, 0, 0, 0, time.UTC), Valid: true},
			Year: 2024, Month: 9, Weekday: 3, Day: 0, Hour: 0, Minute: 0, Second: 0,
		},
	}

	for _, date := range dates {
		_, err := client.DimDatetime.Create().
			SetDate(date.Date).
			SetYear(date.Year).
			SetMonth(date.Month).
			SetWeekday(date.Weekday).
			SetDay(date.Day).
			SetHour(date.Hour).
			SetMinute(date.Minute).
			SetSecond(date.Second).
			Save(ctx)
		if err != nil {
			return fmt.Errorf("failed to create date: %v", err)
		}
	}

	// Inserindo vagas
	vacancies := []ent.DimVacancy{
		{
			DbId: 1, Title: "Software Engineer",
			DimUsrId: 1, NumPositions: 1, ReqId: 1,
			Location:    "São Paulo",
			OpeningDate: &pgtype.Date{Time: time.Date(2024, 9, 1, 0, 0, 0, 0, time.UTC), Valid: true},
			ClosingDate: &pgtype.Date{Time: time.Date(2024, 9, 30, 0, 0, 0, 0, time.UTC), Valid: true},
		},
		{
			DbId: 2, Title: "Data Scientist",
			DimUsrId: 1, NumPositions: 2, ReqId: 1,
			Location:    "Rio de Janeiro",
			OpeningDate: &pgtype.Date{Time: time.Date(2024, 6, 15, 0, 0, 0, 0, time.UTC), Valid: true},
			ClosingDate: &pgtype.Date{Time: time.Date(2024, 7, 15, 0, 0, 0, 0, time.UTC), Valid: true},
		},
		{
			DbId: 3, Title: "HR Specialist",
			DimUsrId: 2, NumPositions: 1, ReqId: 1,
			Location:    "São Paulo",
			OpeningDate: &pgtype.Date{Time: time.Date(2024, 3, 5, 0, 0, 0, 0, time.UTC), Valid: true},
			ClosingDate: &pgtype.Date{Time: time.Date(2024, 3, 25, 0, 0, 0, 0, time.UTC), Valid: true},
		},
		{
			DbId: 4, Title: "UX Designer",
			DimUsrId: 3, NumPositions: 2, ReqId: 1,
			Location:    "Curitiba",
			OpeningDate: &pgtype.Date{Time: time.Date(2024, 8, 10, 0, 0, 0, 0, time.UTC), Valid: true},
			ClosingDate: &pgtype.Date{Time: time.Date(2024, 9, 10, 0, 0, 0, 0, time.UTC), Valid: true},
		},
		{
			DbId: 5, Title: "Software Engineer",
			DimUsrId: 1, NumPositions: 1, ReqId: 2,
			Location:    "São Paulo",
			OpeningDate: &pgtype.Date{Time: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), Valid: true},
			ClosingDate: &pgtype.Date{Time: time.Date(2024, 2, 30, 0, 0, 0, 0, time.UTC), Valid: true},
		},
		{
			DbId: 6, Title: "UX Designer",
			DimUsrId: 5, NumPositions: 1, ReqId: 2,
			Location:    "São Paulo",
			OpeningDate: &pgtype.Date{Time: time.Date(2024, 10, 1, 0, 0, 0, 0, time.UTC), Valid: true},
			ClosingDate: &pgtype.Date{Time: time.Date(2024, 11, 30, 0, 0, 0, 0, time.UTC), Valid: true},
		},
		{
			DbId: 7, Title: "Data Scientist",
			DimUsrId: 4, NumPositions: 1, ReqId: 3,
			Location:    "Rio de Janeiro",
			OpeningDate: &pgtype.Date{Time: time.Date(2024, 9, 15, 0, 0, 0, 0, time.UTC), Valid: true},
			ClosingDate: &pgtype.Date{Time: time.Date(2024, 10, 15, 0, 0, 0, 0, time.UTC), Valid: true},
		},
		{
			DbId: 8, Title: "Product Manager",
			DimUsrId: 5, NumPositions: 1, ReqId: 3,
			Location:    "Belo Horizonte",
			OpeningDate: &pgtype.Date{Time: time.Date(2024, 4, 20, 0, 0, 0, 0, time.UTC), Valid: true},
			ClosingDate: &pgtype.Date{Time: time.Date(2024, 5, 20, 0, 0, 0, 0, time.UTC), Valid: true},
		},
		{
			DbId: 9, Title: "HR Specialist",
			DimUsrId: 3, NumPositions: 1, ReqId: 4,
			Location:    "São Paulo",
			OpeningDate: &pgtype.Date{Time: time.Date(2023, 12, 5, 0, 0, 0, 0, time.UTC), Valid: true},
			ClosingDate: &pgtype.Date{Time: time.Date(2024, 1, 25, 0, 0, 0, 0, time.UTC), Valid: true},
		},
	}

	for _, vacancy := range vacancies {
		_, err := client.DimVacancy.Create().
			SetDbId(vacancy.DbId).
			SetTitle(vacancy.Title).
			SetNumPositions(vacancy.NumPositions).
			SetReqId(vacancy.ReqId).
			SetLocation(vacancy.Location).
			SetOpeningDate(vacancy.OpeningDate).
			SetClosingDate(vacancy.ClosingDate).
			SetDimUsrId(vacancy.DimUsrId).
			Save(ctx)
		if err != nil {
			return fmt.Errorf("failed to create vacancy %s: %v", vacancy.Title, err)
		}
	}

	// Inserindo processos
	processes := []ent.DimProcess{
		{
			DbId: 1, Title: "Desenvolvimento Ágil - Software Engineer",
			InitialDate: &pgtype.Date{Time: time.Date(2024, 8, 1, 0, 0, 0, 0, time.UTC), Valid: true},
			FinishDate:  &pgtype.Date{Time: time.Date(2024, 9, 1, 0, 0, 0, 0, time.UTC), Valid: true},
			DimUsrId:    1,
		},
		{
			DbId: 2, Title: "Recrutamento e Seleção - HR Specialist",
			InitialDate: &pgtype.Date{Time: time.Date(2024, 8, 5, 0, 0, 0, 0, time.UTC), Valid: true},
			FinishDate:  &pgtype.Date{Time: time.Date(2024, 9, 5, 0, 0, 0, 0, time.UTC), Valid: true},
			DimUsrId:    2,
		},
		{
			DbId: 3, Title: "Gestão de Produto - Product Manager",
			InitialDate: &pgtype.Date{Time: time.Date(2024, 7, 20, 0, 0, 0, 0, time.UTC), Valid: true},
			FinishDate:  &pgtype.Date{Time: time.Date(2024, 8, 20, 0, 0, 0, 0, time.UTC), Valid: true},
			DimUsrId:    3,
		},
		{
			DbId: 4, Title: "Experiência do Usuário - UX Designer",
			InitialDate: &pgtype.Date{Time: time.Date(2024, 8, 10, 0, 0, 0, 0, time.UTC), Valid: true},
			FinishDate:  &pgtype.Date{Time: time.Date(2024, 9, 10, 0, 0, 0, 0, time.UTC), Valid: true},
			DimUsrId:    4,
		},
		{
			DbId: 5, Title: "Análise de Dados - Data Scientist",
			InitialDate: &pgtype.Date{Time: time.Date(2024, 7, 15, 0, 0, 0, 0, time.UTC), Valid: true},
			FinishDate:  &pgtype.Date{Time: time.Date(2024, 8, 15, 0, 0, 0, 0, time.UTC), Valid: true},
			DimUsrId:    5,
		},
		{
			DbId: 6, Title: "Desenvolvimento de Software - Software Engineer e UX Designer",
			InitialDate: &pgtype.Date{Time: time.Date(2024, 8, 15, 0, 0, 0, 0, time.UTC), Valid: true},
			FinishDate:  &pgtype.Date{Time: time.Date(2024, 9, 15, 0, 0, 0, 0, time.UTC), Valid: true},
			DimUsrId:    1,
		},
		{
			DbId: 7, Title: "Análise de Dados e Relatórios - Data Scientist e Product Manager",
			InitialDate: &pgtype.Date{Time: time.Date(2024, 8, 20, 0, 0, 0, 0, time.UTC), Valid: true},
			FinishDate:  &pgtype.Date{Time: time.Date(2024, 9, 20, 0, 0, 0, 0, time.UTC), Valid: true},
			DimUsrId:    2,
		},
		{
			DbId: 8, Title: "Processo de Recrutamento - HR Specialist e Software Engineer",
			InitialDate: &pgtype.Date{Time: time.Date(2024, 9, 1, 0, 0, 0, 0, time.UTC), Valid: true},
			FinishDate:  &pgtype.Date{Time: time.Date(2024, 9, 30, 0, 0, 0, 0, time.UTC), Valid: true},
			DimUsrId:    3,
		},
		{
			DbId: 9, Title: "Estratégia de Produto - Product Manager e UX Designer",
			InitialDate: &pgtype.Date{Time: time.Date(2024, 9, 5, 0, 0, 0, 0, time.UTC), Valid: true},
			FinishDate:  &pgtype.Date{Time: time.Date(2024, 10, 5, 0, 0, 0, 0, time.UTC), Valid: true},
			DimUsrId:    4,
		},
		{
			DbId: 10, Title: "Inovação em Dados - Data Scientist e HR Specialist",
			InitialDate: &pgtype.Date{Time: time.Date(2024, 8, 25, 0, 0, 0, 0, time.UTC), Valid: true},
			FinishDate:  &pgtype.Date{Time: time.Date(2024, 9, 25, 0, 0, 0, 0, time.UTC), Valid: true},
			DimUsrId:    5,
		},
	}

	for _, process := range processes {
		_, err := client.DimProcess.Create().
			SetDbId(process.DbId).
			SetTitle(process.Title).
			SetInitialDate(process.InitialDate).
			SetFinishDate(process.FinishDate).
			SetDimUsrId(int(process.DimUsrId)).
			Save(ctx)
		if err != nil {
			return fmt.Errorf("failed to create process %s: %v", process.Title, err)
		}
	}

	// Inserindo dados na tabela fact_hiring_process
	facts := []ent.FactHiringProcess{
		{
			DimProcessId:                  1,
			DimUserId:                     int(userIDs[0]),
			DimVacancyId:                  1,
			DimDateId:                     1,
			MetTotalCandidatesApplied:     10,
			MetTotalCandidatesInterviewed: 5,
			MetTotalCandidatesHired:       3,
			MetSumDurationHiringProces:    30,
			MetSumSalaryInitial:           5000,
			MetTotalFeedbackPositive:      4,
			MetTotalNeutral:               2,
			MetTotalNegative:              1,
		},
		{
			DimProcessId:                  2,
			DimUserId:                     int(userIDs[1]),
			DimVacancyId:                  2,
			DimDateId:                     1,
			MetTotalCandidatesApplied:     12,
			MetTotalCandidatesInterviewed: 6,
			MetTotalCandidatesHired:       4,
			MetSumDurationHiringProces:    25,
			MetSumSalaryInitial:           5500,
			MetTotalFeedbackPositive:      5,
			MetTotalNeutral:               3,
			MetTotalNegative:              2,
		},
		{
			DimProcessId:                  3,
			DimUserId:                     int(userIDs[2]),
			DimVacancyId:                  3,
			DimDateId:                     1,
			MetTotalCandidatesApplied:     8,
			MetTotalCandidatesInterviewed: 4,
			MetTotalCandidatesHired:       2,
			MetSumDurationHiringProces:    20,
			MetSumSalaryInitial:           4500,
			MetTotalFeedbackPositive:      3,
			MetTotalNeutral:               2,
			MetTotalNegative:              1,
		},
		{
			DimProcessId:                  4,
			DimUserId:                     int(userIDs[3]),
			DimVacancyId:                  4,
			DimDateId:                     1,
			MetTotalCandidatesApplied:     15,
			MetTotalCandidatesInterviewed: 8,
			MetTotalCandidatesHired:       5,
			MetSumDurationHiringProces:    35,
			MetSumSalaryInitial:           6000,
			MetTotalFeedbackPositive:      6,
			MetTotalNeutral:               4,
			MetTotalNegative:              2,
		},
		{
			DimProcessId:                  5,
			DimUserId:                     int(userIDs[4]),
			DimVacancyId:                  5,
			DimDateId:                     1,
			MetTotalCandidatesApplied:     20,
			MetTotalCandidatesInterviewed: 10,
			MetTotalCandidatesHired:       6,
			MetSumDurationHiringProces:    40,
			MetSumSalaryInitial:           7000,
			MetTotalFeedbackPositive:      7,
			MetTotalNeutral:               5,
			MetTotalNegative:              3,
		},
	}

	for _, fact := range facts {
		_, err := client.FactHiringProcess.Create().
			SetMetTotalCandidatesApplied(fact.MetTotalCandidatesApplied).
			SetMetTotalCandidatesInterviewed(fact.MetTotalCandidatesInterviewed).
			SetMetTotalCandidatesHired(fact.MetTotalCandidatesHired).
			SetMetSumDurationHiringProces(fact.MetSumDurationHiringProces).
			SetMetSumSalaryInitial(fact.MetSumSalaryInitial).
			SetMetTotalFeedbackPositive(fact.MetTotalFeedbackPositive).
			SetMetTotalNeutral(fact.MetTotalNeutral).
			SetMetTotalNegative(fact.MetTotalNegative).
			SetDimProcessID(fact.DimProcessId).
			SetDimVacancyID(fact.DimVacancyId).
			SetDimUserID(fact.DimUserId).
			SetDimDatetimeID(fact.DimDateId).
			Save(ctx)
		if err != nil {
			return fmt.Errorf("failed to create fact hiring process: %v", err)
		}
	}

	return nil
}
