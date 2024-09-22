package main

import (
	"api5back/ent"
	"api5back/src/database"
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

// Função para popular os dados no banco
func SeedData(client *ent.Client) error {
	ctx := context.Background()

	userMap := make(map[string]int64)

	users := []struct {
		name       string
		occupation string
	}{
		{"Alice Santos", "Recruiter"},
		{"Bob Ferreira", "HR Manager"},
		{"Carla Mendes", "Software Engineer"},
		{"David Costa", "Data Analyst"},
		{"Eva Lima", "Product Manager"},
	}

	for _, user := range users {
		createdUser, err := client.DimUser.Create().
			SetName(user.name).
			SetOcupation(user.occupation).
			Save(ctx)
		if err != nil {
			return fmt.Errorf("failed to create user %s: %v", user.name, err)
		}
		userMap[user.name] = int64(createdUser.ID)
	}

	// Inserindo datas na tabela dim_datetime
	dates := []struct {
		date    time.Time
		year    int
		month   int
		weekday int
		day     int
		hour    int
		minute  int
		second  int
	}{
		{time.Date(2024, 7, 1, 0, 0, 0, 0, time.UTC), 2024, 7, 0, 1, 0, 0, 0},
		{time.Date(2024, 8, 1, 0, 0, 0, 0, time.UTC), 2024, 8, 0, 1, 0, 0, 0},
		{time.Date(2024, 9, 1, 0, 0, 0, 0, time.UTC), 2024, 9, 0, 1, 0, 0, 0},
		{time.Date(2024, 9, 2, 0, 0, 0, 0, time.UTC), 2024, 9, 1, 2, 0, 0, 0},
		{time.Date(2024, 9, 3, 0, 0, 0, 0, time.UTC), 2024, 9, 2, 3, 0, 0, 0},
	}

	for _, dt := range dates {
		dateValue := pgtype.Date{
			Time:  dt.date,
			Valid: true,
		}

		_, err := client.DimDatetime.Create().
			SetDate(&dateValue).
			SetYear(dt.year).
			SetMonth(dt.month).
			SetWeekday(dt.weekday).
			SetDay(dt.day).
			SetHour(dt.hour).
			SetMinute(dt.minute).
			SetSecond(dt.second).
			Save(ctx)
		if err != nil {
			return fmt.Errorf("failed to create date %s: %v", dt.date, err)
		}
	}

	// Inserindo vagas
	vacancies := []struct {
		title        string
		numPositions int
		reqID        int
		location     string
		openingDate  time.Time
		closingDate  time.Time
		userID       int64
	}{
		{"Software Engineer", 1, 1, "São Paulo", time.Date(2024, 9, 1, 0, 0, 0, 0, time.UTC), time.Date(2024, 9, 30, 0, 0, 0, 0, time.UTC), 1},
		{"Data Scientist", 2, 1, "Rio de Janeiro", time.Date(2024, 8, 15, 0, 0, 0, 0, time.UTC), time.Date(2024, 9, 15, 0, 0, 0, 0, time.UTC), 1},
		{"Product Manager", 2, 1, "Belo Horizonte", time.Date(2024, 7, 20, 0, 0, 0, 0, time.UTC), time.Date(2024, 8, 20, 0, 0, 0, 0, time.UTC), 1},
		{"HR Specialist", 1, 1, "São Paulo", time.Date(2024, 9, 5, 0, 0, 0, 0, time.UTC), time.Date(2024, 9, 25, 0, 0, 0, 0, time.UTC), 2},
		{"UX Designer", 2, 1, "Curitiba", time.Date(2024, 8, 10, 0, 0, 0, 0, time.UTC), time.Date(2024, 9, 10, 0, 0, 0, 0, time.UTC), 3},
		{"Software Engineer", 1, 2, "São Paulo", time.Date(2024, 9, 1, 0, 0, 0, 0, time.UTC), time.Date(2024, 9, 30, 0, 0, 0, 0, time.UTC), 1},
		{"UX Designer", 1, 2, "São Paulo", time.Date(2024, 9, 1, 0, 0, 0, 0, time.UTC), time.Date(2024, 9, 30, 0, 0, 0, 0, time.UTC), 2},
		{"Data Scientist", 1, 3, "Rio de Janeiro", time.Date(2024, 8, 15, 0, 0, 0, 0, time.UTC), time.Date(2024, 9, 15, 0, 0, 0, 0, time.UTC), 3},
		{"Product Manager", 1, 3, "Belo Horizonte", time.Date(2024, 7, 20, 0, 0, 0, 0, time.UTC), time.Date(2024, 8, 20, 0, 0, 0, 0, time.UTC), 4},
		{"HR Specialist", 1, 4, "São Paulo", time.Date(2024, 9, 5, 0, 0, 0, 0, time.UTC), time.Date(2024, 9, 25, 0, 0, 0, 0, time.UTC), 5},
	}

	for _, vacancy := range vacancies {
		_, err := client.DimVacancy.Create().
			SetTitle(vacancy.title).
			SetNumPositions(vacancy.numPositions).
			SetReqId(vacancy.reqID).
			SetLocation(vacancy.location).
			SetOpeningDate(vacancy.openingDate).
			SetClosingDate(int(vacancy.closingDate.Unix())).
			SetDimUsrId(int(vacancy.userID)).
			Save(ctx)
		if err != nil {
			return fmt.Errorf("failed to create vacancy %s: %v", vacancy.title, err)
		}
	}

	// Inserindo processos
	processes := []struct {
		title       string
		initialDate time.Time
		finishDate  time.Time
		userID      int64
	}{
		{"Desenvolvimento Ágil - Software Engineer", time.Date(2024, 8, 1, 0, 0, 0, 0, time.UTC), time.Date(2024, 9, 1, 0, 0, 0, 0, time.UTC), 1},
		{"Recrutamento e Seleção - HR Specialist", time.Date(2024, 8, 5, 0, 0, 0, 0, time.UTC), time.Date(2024, 9, 5, 0, 0, 0, 0, time.UTC), 2},
		{"Gestão de Produto - Product Manager", time.Date(2024, 7, 20, 0, 0, 0, 0, time.UTC), time.Date(2024, 8, 20, 0, 0, 0, 0, time.UTC), 3},
		{"Experiência do Usuário - UX Designer", time.Date(2024, 8, 10, 0, 0, 0, 0, time.UTC), time.Date(2024, 9, 10, 0, 0, 0, 0, time.UTC), 4},
		{"Análise de Dados - Data Scientist", time.Date(2024, 7, 15, 0, 0, 0, 0, time.UTC), time.Date(2024, 8, 15, 0, 0, 0, 0, time.UTC), 5},

		{"Desenvolvimento de Software - Software Engineer e UX Designer", time.Date(2024, 8, 15, 0, 0, 0, 0, time.UTC), time.Date(2024, 9, 15, 0, 0, 0, 0, time.UTC), 1},
		{"Análise de Dados e Relatórios - Data Scientist e Product Manager", time.Date(2024, 8, 20, 0, 0, 0, 0, time.UTC), time.Date(2024, 9, 20, 0, 0, 0, 0, time.UTC), 2},
		{"Processo de Recrutamento - HR Specialist e Software Engineer", time.Date(2024, 9, 1, 0, 0, 0, 0, time.UTC), time.Date(2024, 9, 30, 0, 0, 0, 0, time.UTC), 3},
		{"Estratégia de Produto - Product Manager e UX Designer", time.Date(2024, 9, 5, 0, 0, 0, 0, time.UTC), time.Date(2024, 10, 5, 0, 0, 0, 0, time.UTC), 4},
		{"Inovação em Dados - Data Scientist e HR Specialist", time.Date(2024, 8, 25, 0, 0, 0, 0, time.UTC), time.Date(2024, 9, 25, 0, 0, 0, 0, time.UTC), 5},
	}

	for _, process := range processes {
		_, err := client.DimProcess.Create().
			SetTitle(process.title).
			SetInitialDate(process.initialDate).
			SetFinishDate(process.finishDate).
			SetDimUsrId(int(process.userID)).
			Save(ctx)
		if err != nil {
			return fmt.Errorf("failed to create process %s: %v", process.title, err)
		}
	}

	// Inserindo dados na tabela fact_hiring_process
	facts := []struct {
		totalCandidatesApplied     int64
		totalCandidatesInterviewed int64
		totalCandidatesHired       int64
		sumDurationHiringProcess   int64
		sumSalaryInitial           int64
		totalFeedbackPositive      int64
		totalNeutral               int64
		totalNegative              int64
		processID                  int64
		vacancyID                  int64
		userID                     int64
		dateID                     int64
	}{
		{10, 5, 3, 30, 5000, 4, 2, 1, 1, 1, 1, 1},
		{12, 6, 4, 25, 5500, 5, 3, 2, 2, 1, 1, 1},
		{8, 4, 2, 20, 4500, 3, 2, 1, 3, 1, 1, 1},
		{15, 8, 5, 35, 6000, 6, 4, 2, 4, 1, 1, 1},
		{20, 10, 6, 40, 7000, 7, 5, 3, 5, 1, 1, 1},
	}

	for _, fact := range facts {
		_, err := client.FactHiringProcess.Create().
			SetMetTotalCandidatesApplied(int(fact.totalCandidatesApplied)).
			SetMetTotalCandidatesInterviewed(int(fact.totalCandidatesInterviewed)).
			SetMetTotalCandidatesHired(int(fact.totalCandidatesHired)).
			SetMetSumDurationHiringProces(int(fact.sumDurationHiringProcess)).
			SetMetSumSalaryInitial(int(fact.sumSalaryInitial)).
			SetMetTotalFeedbackPositive(int(fact.totalFeedbackPositive)).
			SetMetTotalNeutral(int(fact.totalNeutral)).
			SetMetTotalNegative(int(fact.totalNegative)).
			SetDimProcessID(int(fact.processID)).
			SetDimVacancyID(int(fact.vacancyID)).
			SetDimUserID(int(fact.userID)).
			SetDimDatetimeID(int(fact.dateID)).
			Save(ctx)
		if err != nil {
			return fmt.Errorf("failed to create fact hiring process: %v", err)
		}
	}

	fmt.Println("Database seeded successfully!")
	return nil
}

func SeedAll() error {
	dwPrefix := "DW"
	client, err := database.Setup(dwPrefix)
	if err != nil {
		return fmt.Errorf("failed to setup database: %v", err)
	}
	defer client.Close()

	if err := SeedData(client); err != nil {
		return fmt.Errorf("failed to seed database: %v", err)
	}

	fmt.Printf("Seeded database with prefix: %s\n", dwPrefix)

	return nil
}

// Ponto de entrada manual para rodar a seed
func main() {
	fmt.Println("Seeding all databases...")

	if err := SeedAll(); err != nil {
		panic(fmt.Errorf("failed to seed databases: %v", err))
	}

	fmt.Println("Successfully seeded databases.")
}
