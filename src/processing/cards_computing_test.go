package processing

import (
	"testing"
	"time"

	"api5back/ent"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
)

func TestComputingCardInfo(t *testing.T) {
	// Criação de dados fictícios de DimProcess
	process1 := &ent.DimProcess{
		Status:      1,
		InitialDate: &pgtype.Date{Time: time.Now().Add(-10 * 24 * time.Hour), Valid: true}, // Início há 10 dias
		FinishDate:  &pgtype.Date{Time: time.Now().Add(2 * 24 * time.Hour), Valid: true},   // Termina em 2 dias
	}
	process2 := &ent.DimProcess{
		Status:      2,
		InitialDate: &pgtype.Date{Time: time.Now().Add(-15 * 24 * time.Hour), Valid: true}, // Início há 15 dias
		FinishDate:  &pgtype.Date{Time: time.Now().Add(-5 * 24 * time.Hour), Valid: true},  // Terminou há 5 dias
	}
	process3 := &ent.DimProcess{
		Status:      3,
		InitialDate: &pgtype.Date{Time: time.Now().Add(-30 * 24 * time.Hour), Valid: true}, // Início há 30 dias
		FinishDate:  &pgtype.Date{Time: time.Now().Add(-20 * 24 * time.Hour), Valid: true}, // Terminou há 20 dias
	}

	// Criação de dados fictícios de FactHiringProcess
	hiringData := []*ent.FactHiringProcess{
		{
			Edges: ent.FactHiringProcessEdges{
				DimProcess: process1,
			},
			MetSumDurationHiringProces: 10,
		},
		{
			Edges: ent.FactHiringProcessEdges{
				DimProcess: process2,
			},
			MetSumDurationHiringProces: 15,
		},
		{
			Edges: ent.FactHiringProcessEdges{
				DimProcess: process3,
			},
			MetSumDurationHiringProces: 20,
		},
	}

	// Chama a função ComputingCardInfo com os dados de teste
	cardInfos, err := ComputingCardInfo(hiringData)

	// Verifica se não houve erro
	assert.NoError(t, err)

	// Verifica os valores retornados
	assert.Equal(t, 1, cardInfos.Open)
	assert.Equal(t, 1, cardInfos.InProgress)
	assert.Equal(t, 1, cardInfos.Closed)
	assert.Equal(t, 1, cardInfos.ApproachingDeadline)
	assert.Equal(t, 15, cardInfos.AverageHiringTime)
}

func TestComputingCardInfo_EmptyData(t *testing.T) {
	// Chama a função com uma lista vazia
	cardInfos, err := ComputingCardInfo([]*ent.FactHiringProcess{})

	// Verifica se não houve erro
	assert.NoError(t, err)

	// Verifica se os valores do cardInfos são os valores padrão (zero)
	assert.Equal(t, CardInfos{}, cardInfos)
}
