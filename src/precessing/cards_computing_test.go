package processing

import (
	"api5back/ent"
	"testing"
	"time"
)

// Função auxiliar para criar um DimProcess mock
func mockDimProcess(status int, initialDate, finishDate time.Time) *ent.DimProcess {
	return &ent.DimProcess{
		Status:      status,
		InitialDate: initialDate,
		FinishDate:  finishDate,
	}
}

// Função auxiliar para criar um FactHiringProcess mock
func mockFactHiringProcess(duration int) *ent.FactHiringProcess {
	return &ent.FactHiringProcess{
		MetSumDurationHiringProces: duration,
	}
}

// Testa a função ComputingCardInfo
func TestComputingCardInfo(t *testing.T) {
	processes := []*ent.DimProcess{
		mockDimProcess(1, time.Now().Add(-5*time.Hour), time.Now().Add(1*time.Hour)),   // Open, approaching deadline
		mockDimProcess(2, time.Now().Add(-10*time.Hour), time.Now().Add(-1*time.Hour)), // Expired
		mockDimProcess(3, time.Now().Add(-8*time.Hour), time.Now().Add(-2*time.Hour)),  // Closed
	}

	hirings := []*ent.FactHiringProcess{
		mockFactHiringProcess(10),
		mockFactHiringProcess(20),
	}

	cardInfo, err := ComputingCardInfo(processes, hirings)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Verifica se o cardInfo retornado está correto
	if cardInfo.openProcess != 1 {
		t.Errorf("Expected 1 open process, got %d", cardInfo.openProcess)
	}

	if cardInfo.expirededProcess != 1 {
		t.Errorf("Expected 1 expired process, got %d", cardInfo.expirededProcess)
	}

	if cardInfo.closeProcess != 1 {
		t.Errorf("Expected 1 closed process, got %d", cardInfo.closeProcess)
	}

	if cardInfo.approachingDeadlineProcess != 1 {
		t.Errorf("Expected 1 approaching deadline process, got %d", cardInfo.approachingDeadlineProcess)
	}

	if cardInfo.averageHiringTime != 15 { // (10 + 20) / 2 = 15
		t.Errorf("Expected average hiring time 15, got %d", cardInfo.averageHiringTime)
	}
}

// Testa a função getProcessCardInfo
func TestGetProcessCardInfo(t *testing.T) {
	processes := []*ent.DimProcess{
		mockDimProcess(1, time.Now().Add(-5*time.Hour), time.Now().Add(1*time.Hour)),   // Open, approaching deadline
		mockDimProcess(2, time.Now().Add(-10*time.Hour), time.Now().Add(-1*time.Hour)), // Expired
		mockDimProcess(3, time.Now().Add(-8*time.Hour), time.Now().Add(-2*time.Hour)),  // Closed
	}

	var cardInfos CardInfos
	err := getProcessCardInfo(processes, &cardInfos)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Verifica se o cardInfos retornado está correto
	if cardInfos.openProcess != 1 {
		t.Errorf("Expected 1 open process, got %d", cardInfos.openProcess)
	}

	if cardInfos.expirededProcess != 1 {
		t.Errorf("Expected 1 expired process, got %d", cardInfos.expirededProcess)
	}

	if cardInfos.closeProcess != 1 {
		t.Errorf("Expected 1 closed process, got %d", cardInfos.closeProcess)
	}

	if cardInfos.approachingDeadlineProcess != 1 {
		t.Errorf("Expected 1 approaching deadline process, got %d", cardInfos.approachingDeadlineProcess)
	}
}

// Testa a função getAverageHiringTime
func TestGetAverageHiringTime(t *testing.T) {
	hirings := []*ent.FactHiringProcess{
		mockFactHiringProcess(10),
		mockFactHiringProcess(20),
	}

	var cardInfos CardInfos
	err := getAverageHiringTime(hirings, &cardInfos)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Verifica se o tempo médio de contratação está correto
	if cardInfos.averageHiringTime != 15 { // (10 + 20) / 2 = 15
		t.Errorf("Expected average hiring time 15, got %d", cardInfos.averageHiringTime)
	}
}

// Testa erro de hiring data vazio
func TestGetAverageHiringTimeEmptyData(t *testing.T) {
	var hirings []*ent.FactHiringProcess

	var cardInfos CardInfos
	err := getAverageHiringTime(hirings, &cardInfos)
	if err == nil {
		t.Error("Expected error for empty hiring data, got none")
	}
}
