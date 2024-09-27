package seeds

import (
	"context"
	"fmt"
	"math/rand"

	"api5back/ent"
	"api5back/src/property"

	"entgo.io/ent/dialect/sql"
	"github.com/jackc/pgx/v5/pgtype"
)

// possible first names
var firstNames = []string{
	"John",
	"Jane",
	"Michael",
	"Emily",
	"David",
	"Sarah",
	"James",
	"Jessica",
	"Robert",
	"Jennifer",
	"Walter",
	"Lisa",
	"Richard",
	"Mary",
	"Charles",
	"Karen",
}

// possible last names
var lastNames = []string{
	"Doe",
	"Smith",
	"Johnson",
	"Brown",
	"Williams",
	"Jones",
	"Miller",
	"Davis",
	"White",
	"Clark",
	"Moore",
	"Taylor",
	"Anderson",
	"Thomas",
	"Jackson",
	"Harris",
}

func randomName() [2]string {
	return [2]string{
		firstNames[rand.Intn(len(firstNames))],
		lastNames[rand.Intn(len(lastNames))],
	}
}

func DwProceduralHiringProcessCandidates(client *ent.Client) error {
	ctx := context.Background()

	// select random FactHiringProcess from the database (max 10)
	factHiringProcess, err := client.
		FactHiringProcess.
		Query().
		WithDimVacancy().
		Order(sql.OrderByRand()).
		Limit(10).
		All(ctx)
	if err != nil {
		return fmt.Errorf("failed to query FactHiringProcess: %v", err)
	}

	var candidatesToInsert []*ent.HiringProcessCandidateCreate

	// loop through the FactHiringProcess and create 5 to 10 candidates for each
	for _, fhp := range factHiringProcess {
		numberOfCandidates := rand.Intn(6) + 5

		for i := 0; i < numberOfCandidates; i++ {
			candidateName := randomName()
			candidateStatus := property.HiringProcessCandidateStatus(rand.Intn(4))
			factHiringProcessVacancy, err := fhp.Edges.DimVacancyOrErr()
			if err != nil {
				return fmt.Errorf("failed to get vacandy of FactHiringProcess: %v", err)
			}

			applyDate := factHiringProcessVacancy.OpeningDate.AddDate(
				0, 0, rand.Intn(
					int(factHiringProcessVacancy.
						ClosingDate.
						Sub(factHiringProcessVacancy.OpeningDate).
						Hours()/24)+1,
				),
			)
			applyDatePgType := &pgtype.Date{}
			if err := applyDatePgType.Scan(applyDate); err != nil {
				return fmt.Errorf("failed to generate random applyDate for candidate: %v", err)
			}

			candidatesToInsert = append(candidatesToInsert, client.
				HiringProcessCandidate.
				Create().
				SetFactHiringProcess(fhp).
				SetName(fmt.Sprintf(
					"%s %s",
					candidateName[0],
					candidateName[1],
				)).
				SetEmail(fmt.Sprintf(
					"%s.%s-%d@khali.com",
					candidateName[0],
					candidateName[1],
					rand.Intn(1000),
				)).
				SetPhone(fmt.Sprintf(
					"+1%010d",
					rand.Intn(10000000000),
				)).
				SetScore(rand.Float64()*100).
				SetApplyDate(applyDatePgType).
				SetStatus(candidateStatus))
		}
	}

	if _, err = client.
		HiringProcessCandidate.
		CreateBulk(candidatesToInsert...).
		Save(ctx); err != nil {
		return fmt.Errorf("failed to create candidate: %v", err)
	}

	return nil
}
