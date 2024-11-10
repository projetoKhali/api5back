package seeds

import (
	"context"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"time"

	"api5back/ent"
	"api5back/src/property"

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

	for _, user := range users {
		_, err := client.DimUser.Create().
			SetDbId(user.DbId).
			SetName(user.Name).
			SetOccupation(user.Occupation).
			Save(ctx)
		if err != nil {
			return fmt.Errorf("failed to create user %s: %v", user.Name, err)
		}
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
			OpeningDate: &pgtype.Date{Time: time.Date(2024, 12, 1, 0, 0, 0, 0, time.UTC), Valid: true},
			ClosingDate: &pgtype.Date{Time: time.Date(2024, 12, 30, 0, 0, 0, 0, time.UTC), Valid: true},
			Status:      property.DimVacancyStatusOpen,
		},
		{
			DbId: 2, Title: "Data Scientist",
			DimUsrId: 1, NumPositions: 2, ReqId: 1,
			Location:    "Rio de Janeiro",
			OpeningDate: &pgtype.Date{Time: time.Date(2024, 6, 15, 0, 0, 0, 0, time.UTC), Valid: true},
			ClosingDate: &pgtype.Date{Time: time.Date(2024, 7, 15, 0, 0, 0, 0, time.UTC), Valid: true},
			Status:      property.DimVacancyStatusOpen,
		},
		{
			DbId: 3, Title: "HR Specialist",
			DimUsrId: 2, NumPositions: 1, ReqId: 1,
			Location:    "São Paulo",
			OpeningDate: &pgtype.Date{Time: time.Date(2024, 3, 5, 0, 0, 0, 0, time.UTC), Valid: true},
			ClosingDate: &pgtype.Date{Time: time.Date(2024, 3, 25, 0, 0, 0, 0, time.UTC), Valid: true},
			Status:      property.DimVacancyStatusOpen,
		},
		{
			DbId: 4, Title: "UX Designer",
			DimUsrId: 3, NumPositions: 2, ReqId: 1,
			Location:    "Curitiba",
			OpeningDate: &pgtype.Date{Time: time.Date(2024, 8, 10, 0, 0, 0, 0, time.UTC), Valid: true},
			ClosingDate: &pgtype.Date{Time: time.Date(2024, 9, 10, 0, 0, 0, 0, time.UTC), Valid: true},
			Status:      property.DimVacancyStatusOpen,
		},
		{
			DbId: 5, Title: "Software Engineer",
			DimUsrId: 1, NumPositions: 1, ReqId: 2,
			Location:    "São Paulo",
			OpeningDate: &pgtype.Date{Time: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), Valid: true},
			ClosingDate: &pgtype.Date{Time: time.Date(2024, 2, 30, 0, 0, 0, 0, time.UTC), Valid: true},
			Status:      property.DimVacancyStatusOpen,
		},
		{
			DbId: 6, Title: "UX Designer",
			DimUsrId: 5, NumPositions: 1, ReqId: 2,
			Location:    "São Paulo",
			OpeningDate: &pgtype.Date{Time: time.Date(2024, 10, 1, 0, 0, 0, 0, time.UTC), Valid: true},
			ClosingDate: &pgtype.Date{Time: time.Date(2024, 11, 30, 0, 0, 0, 0, time.UTC), Valid: true},
			Status:      property.DimVacancyStatusOpen,
		},
		{
			DbId: 7, Title: "Data Scientist",
			DimUsrId: 4, NumPositions: 1, ReqId: 3,
			Location:    "Rio de Janeiro",
			OpeningDate: &pgtype.Date{Time: time.Date(2024, 9, 15, 0, 0, 0, 0, time.UTC), Valid: true},
			ClosingDate: &pgtype.Date{Time: time.Date(2024, 10, 15, 0, 0, 0, 0, time.UTC), Valid: true},
			Status:      property.DimVacancyStatusOpen,
		},
		{
			DbId: 8, Title: "Product Manager",
			DimUsrId: 5, NumPositions: 1, ReqId: 3,
			Location:    "Belo Horizonte",
			OpeningDate: &pgtype.Date{Time: time.Date(2024, 4, 20, 0, 0, 0, 0, time.UTC), Valid: true},
			ClosingDate: &pgtype.Date{Time: time.Date(2024, 5, 20, 0, 0, 0, 0, time.UTC), Valid: true},
			Status:      property.DimVacancyStatusOpen,
		},
		{
			DbId: 9, Title: "HR Specialist",
			DimUsrId: 3, NumPositions: 2, ReqId: 4,
			Location:    "São Paulo",
			OpeningDate: &pgtype.Date{Time: time.Date(2024, 4, 5, 0, 0, 0, 0, time.UTC), Valid: true},
			ClosingDate: &pgtype.Date{Time: time.Date(2024, 4, 25, 0, 0, 0, 0, time.UTC), Valid: true},
			Status:      property.DimVacancyStatusInAnalysis,
		},
		{
			DbId: 10, Title: "Data Engineer",
			DimUsrId: 1, NumPositions: 3, ReqId: 5,
			Location:    "Sergipe",
			OpeningDate: &pgtype.Date{Time: time.Date(2024, 8, 25, 0, 0, 0, 0, time.UTC), Valid: true},
			ClosingDate: &pgtype.Date{Time: time.Date(2024, 9, 25, 0, 0, 0, 0, time.UTC), Valid: true},
			Status:      property.DimVacancyStatusClosed,
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
			SetStatus(vacancy.Status).
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
			Status:      property.DimProcessStatusOpen,
		},
		{
			DbId: 2, Title: "Recrutamento e Seleção - HR Specialist",
			InitialDate: &pgtype.Date{Time: time.Date(2024, 8, 5, 0, 0, 0, 0, time.UTC), Valid: true},
			FinishDate:  &pgtype.Date{Time: time.Date(2024, 9, 5, 0, 0, 0, 0, time.UTC), Valid: true},
			DimUsrId:    2,
			Status:      property.DimProcessStatusOpen,
		},
		{
			DbId: 3, Title: "Gestão de Produto - Product Manager",
			InitialDate: &pgtype.Date{Time: time.Date(2024, 7, 20, 0, 0, 0, 0, time.UTC), Valid: true},
			FinishDate:  &pgtype.Date{Time: time.Date(2024, 8, 20, 0, 0, 0, 0, time.UTC), Valid: true},
			DimUsrId:    3,
			Status:      property.DimProcessStatusOpen,
		},
		{
			DbId: 4, Title: "Experiência do Usuário - UX Designer",
			InitialDate: &pgtype.Date{Time: time.Date(2024, 8, 10, 0, 0, 0, 0, time.UTC), Valid: true},
			FinishDate:  &pgtype.Date{Time: time.Date(2024, 9, 10, 0, 0, 0, 0, time.UTC), Valid: true},
			DimUsrId:    4,
			Status:      property.DimProcessStatusOpen,
		},
		{
			DbId: 5, Title: "Análise de Dados - Data Scientist",
			InitialDate: &pgtype.Date{Time: time.Date(2024, 7, 15, 0, 0, 0, 0, time.UTC), Valid: true},
			FinishDate:  &pgtype.Date{Time: time.Date(2024, 8, 15, 0, 0, 0, 0, time.UTC), Valid: true},
			DimUsrId:    5,
			Status:      property.DimProcessStatusOpen,
		},
		{
			DbId: 6, Title: "Desenvolvimento de Software - Software Engineer e UX Designer",
			InitialDate: &pgtype.Date{Time: time.Date(2024, 8, 15, 0, 0, 0, 0, time.UTC), Valid: true},
			FinishDate:  &pgtype.Date{Time: time.Date(2024, 9, 15, 0, 0, 0, 0, time.UTC), Valid: true},
			DimUsrId:    1,
			Status:      property.DimProcessStatusOpen,
		},
		{
			DbId: 7, Title: "Análise de Dados e Relatórios - Data Scientist e Product Manager",
			InitialDate: &pgtype.Date{Time: time.Date(2024, 8, 20, 0, 0, 0, 0, time.UTC), Valid: true},
			FinishDate:  &pgtype.Date{Time: time.Date(2024, 9, 20, 0, 0, 0, 0, time.UTC), Valid: true},
			DimUsrId:    2,
			Status:      property.DimProcessStatusOpen,
		},
		{
			DbId: 8, Title: "Processo de Recrutamento - HR Specialist e Software Engineer",
			InitialDate: &pgtype.Date{Time: time.Date(2024, 9, 1, 0, 0, 0, 0, time.UTC), Valid: true},
			FinishDate:  &pgtype.Date{Time: time.Date(2024, 9, 30, 0, 0, 0, 0, time.UTC), Valid: true},
			DimUsrId:    3,
			Status:      property.DimProcessStatusOpen,
		},
		{
			DbId: 9, Title: "Estratégia de Produto - Product Manager e UX Designer",
			InitialDate: &pgtype.Date{Time: time.Date(2024, 9, 5, 0, 0, 0, 0, time.UTC), Valid: true},
			FinishDate:  &pgtype.Date{Time: time.Date(2024, 10, 5, 0, 0, 0, 0, time.UTC), Valid: true},
			DimUsrId:    4,
			Status:      property.DimProcessStatusInProgress,
		},
		{
			DbId: 10, Title: "Inovação em Dados - Data Scientist e HR Specialist",
			InitialDate: &pgtype.Date{Time: time.Date(2024, 8, 25, 0, 0, 0, 0, time.UTC), Valid: true},
			FinishDate:  &pgtype.Date{Time: time.Date(2024, 9, 25, 0, 0, 0, 0, time.UTC), Valid: true},
			DimUsrId:    5,
			Status:      property.DimProcessStatusClosed,
		},
	}

	for _, process := range processes {
		_, err := client.DimProcess.Create().
			SetDbId(process.DbId).
			SetTitle(process.Title).
			SetInitialDate(process.InitialDate).
			SetFinishDate(process.FinishDate).
			SetDimUsrId(int(process.DimUsrId)).
			SetStatus(process.Status).
			Save(ctx)
		if err != nil {
			return fmt.Errorf("failed to create process %s: %v", process.Title, err)
		}
	}

	// Inserindo dados na tabela fact_hiring_process
	facts := []ent.FactHiringProcess{
		{
			DimUserId:                     1,
			DimProcessId:                  1,
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
			DimUserId:                     2,
			DimProcessId:                  2,
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
			DimUserId:                     3,
			DimProcessId:                  3,
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
			DimUserId:                     4,
			DimProcessId:                  4,
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
			DimUserId:                     5,
			DimProcessId:                  5,
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
		{
			DimUserId:                     2,
			DimProcessId:                  6,
			DimVacancyId:                  6,
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
			DimUserId:                     3,
			DimProcessId:                  7,
			DimVacancyId:                  7,
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
			DimUserId:                     4,
			DimProcessId:                  8,
			DimVacancyId:                  8,
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
			DimUserId:                     5,
			DimProcessId:                  9,
			DimVacancyId:                  9,
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
			DimUserId:                     1,
			DimProcessId:                  10,
			DimVacancyId:                  10,
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

	logCurrentCandidate := 1
	var logTotalCandidates int

	for _, fact := range facts {
		logTotalCandidates += fact.MetTotalCandidatesApplied
	}

	for factId, fact := range facts {
		factProgressString := fmt.Sprintf(
			"Creating fact hiring process %d/%d (%d%%)",
			factId+1,
			len(facts),
			int(float64(factId+1)/float64(len(facts))*100),
		)

		fmt.Printf("%s\n", factProgressString)

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

		factCurrentCandidateDbId := 1
		factTotalCandidatesRemaining := fact.MetTotalCandidatesApplied - fact.MetTotalCandidatesHired - fact.MetTotalCandidatesInterviewed

		for candidateStatus, candidateCategory := range [4]int{
			factTotalCandidatesRemaining / 2,
			fact.MetTotalCandidatesInterviewed,
			fact.MetTotalCandidatesHired,
			factTotalCandidatesRemaining / 2,
		} {
			for j := 0; j < candidateCategory; j++ {
				fmt.Printf(
					"%s • Creating fact candidate %d/%d (%d%%)\n",
					factProgressString,
					logCurrentCandidate,
					logTotalCandidates,
					int(float64(logCurrentCandidate)/float64(logTotalCandidates)*100),
				)

				t := float64(j) / float64(fact.MetTotalCandidatesApplied-1)
				candidateName := generateName(factCurrentCandidateDbId)
				candidateApplyDate := lerpDate(
					*vacancies[fact.DimVacancyId-1].OpeningDate,
					*vacancies[fact.DimVacancyId-1].ClosingDate,
					t,
				)

				candidateBuilder := client.HiringProcessCandidate.Create().
					SetDbId(factCurrentCandidateDbId).
					SetName(candidateName).
					SetEmail(fmt.Sprintf("%s_%d@mail.com", candidateName, factCurrentCandidateDbId)).
					SetPhone(pseudoRandomPhoneIdempotent(factCurrentCandidateDbId)).
					SetScore(pseudoRandomScoreIdempotent(factCurrentCandidateDbId)).
					SetFactHiringProcessID(factId + 1).
					SetApplyDate(candidateApplyDate).
					SetStatus(property.HiringProcessCandidateStatus(candidateStatus))

				if property.HiringProcessCandidateStatus(candidateStatus) > property.HiringProcessCandidateStatusInAnalysis {
					candidateBuilder.SetUpdatedAt(lerpDate(
						*candidateApplyDate,
						*vacancies[fact.DimVacancyId-1].ClosingDate,
						0.5,
					))
				}

				_, err := candidateBuilder.Save(ctx)
				if err != nil {
					return fmt.Errorf("failed to create fact candidate: %v", err)
				}

				factCurrentCandidateDbId++
				logCurrentCandidate++

				fmt.Print("\033[F\033[K") // \033[F moves up a line, \033[K clears the line
			}
		}

		fmt.Print("\033[F\033[K") // \033[F moves up a line, \033[K clears the line
	}

	return nil
}

func generateName(index int) string {
	// Convert the index to bytes and hash it
	indexBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(indexBytes, uint64(index))
	hash := sha256.Sum256(indexBytes)

	// Use first byte for first name and second byte for last name
	firstName := firstNames[int(hash[0])%len(firstNames)]
	lastName := lastNames[int(hash[1])%len(lastNames)]

	// Combine into full name
	return fmt.Sprintf("%s %s", firstName, lastName)
}

func pseudoRandomPhoneIdempotent(index int) string {
	// Convert the index to bytes
	indexBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(indexBytes, uint64(index))

	// Compute SHA-256 hash of the index
	hash := sha256.Sum256(indexBytes)

	// Use the first 5 bytes of the hash to generate a consistent 10-digit number
	phonePart1 := int(binary.BigEndian.Uint16(hash[:2]) % 1000)   // 3 digits
	phonePart2 := int(binary.BigEndian.Uint16(hash[2:4]) % 1000)  // 3 digits
	phonePart3 := int(binary.BigEndian.Uint16(hash[4:6]) % 10000) // 4 digits

	// Format into phone number style and return as a string
	return fmt.Sprintf("11%03d%03d%04d", phonePart1, phonePart2, phonePart3)
}

func pseudoRandomScoreIdempotent(index int) float64 {
	// Convert the index to bytes
	indexBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(indexBytes, uint64(index))

	// Compute SHA-256 hash of the index
	hash := sha256.Sum256(indexBytes)

	// Use the first 8 bytes of the hash to get a consistent integer
	hashInt := binary.BigEndian.Uint64(hash[:8])

	// Map hashInt to a float between 0.0 and 100.0
	return float64(hashInt%10000) / 100.0
}

// lerpDate linearly interpolates between two pgtype.Date values based on t (0 to 1).
func lerpDate(date1, date2 pgtype.Date, t float64) *pgtype.Date {
	t1 := date1.Time
	t2 := date2.Time

	if t < 0 {
		t = 0
	} else if t > 1 {
		t = 1
	}

	// Calculate the interpolated duration between t1 and t2
	interpolated := t1.Add(time.Duration(float64(t2.Sub(t1)) * t))

	return &pgtype.Date{Time: interpolated, Valid: true}
}
