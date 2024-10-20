package server

import (
	"api5back/ent"
	"api5back/src/processing"
	"api5back/src/service"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SuggestionsResponse struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type TableResponse struct {
	Title             string   `json:"title"`
	NumPositions      int      `json:"numPositions"`
	NumCandidates     int      `json:"numCandidates"`
	CompetitionRate   *float32 `json:"competitionRate"`
	NumInterviewed    int      `json:"numInterviewed"`
	NumHired          int      `json:"numHired"`
	AverageHiringTime *float32 `json:"averageHiringTime"`
	NumFeedback       int      `json:"numFeedback"`
}

func HiringProcessDashboard(
	engine *gin.Engine,
	dbClient *ent.Client,
	dwClient *ent.Client,
) {
	v1 := engine.Group("/api/v1")
	{
		eg := v1.Group("/hiring-process")
		{
			eg.GET("/dashboard", Dashboard(dbClient, dwClient))
		}

		suggestions := v1.Group("/suggestions")
		{
			suggestions.GET("/recruiter", UserList(dwClient))
			suggestions.POST("/process", HiringProcessList((dwClient)))
			suggestions.POST("/vacancies", VacancyList(dwClient))
		}

		table := v1.Group("/table")
		{
			table.POST("/dashboard", VacancyTable(dwClient))
		}
	}
}

// Dashboard godoc
// @Summary dashboard
// @Schemes
// @Description show dashboard
// @Tags hiring-process
// @Accept json
// @Produce json
// @Success 200 {string} Dashboard
// @Router /hiring-process/dashboard [get]
func Dashboard(
	dbClient *ent.Client,
	dwClient *ent.Client,
) func(c *gin.Context) {
	return func(c *gin.Context) {
		MetricsService := service.NewMetricsService(dwClient)

		hiringProcessName := c.Query("hiringProcess")
		vacancyName := c.Query("vacancy")
		startDate := c.Query("startDate")
		endDate := c.Query("endDate")

		metricsData, err := MetricsService.GetMetrics(
			c,
			service.GetMetricsFilter{
				HiringProcessName: hiringProcessName,
				VacancyName:       vacancyName,
				StartDate:         startDate,
				EndDate:           endDate,
			},
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		c.JSON(http.StatusOK, metricsData)
	}
}

func TableData(
	dbClient *ent.Client,
	dwClient *ent.Client,
) func(c *gin.Context) {
	return func(c *gin.Context) {
		var userIDs []int

		// Parse the body for user IDs
		if err := c.ShouldBindJSON(&userIDs); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}
	}
}

// UserList godoc
// @Summary List users
// @Schemes
// @Description Return a list of users with id and name
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {array} SuggestionsResponse
// @Router /users/ [get]
func UserList(
	dwClient *ent.Client,
) func(c *gin.Context) {
	return func(c *gin.Context) {
		userService := service.NewUserService(dwClient)

		users, err := userService.GetUsers(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		var response []SuggestionsResponse
		for _, user := range users {
			response = append(response, SuggestionsResponse{
				Id:   user.ID,
				Name: user.Name,
			})
		}

		c.JSON(http.StatusOK, response)
	}
}

// HiringProcessList godoc
// @Summary List hiring processes
// @Schemes
// @Description Return a list of hiring processes with id and title
// @Tags hiring-process
// @Accept json
// @Produce json
// @Success 200 {array} SuggestionsResponse
// @Router /hiring-process [post]
func HiringProcessList(
	dbClient *ent.Client,
) func(c *gin.Context) {
	return func(c *gin.Context) {
		var userIDs []int

		// Parse the body for user IDs
		if err := c.ShouldBindJSON(&userIDs); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		hiringProcessService := service.NewHiringProcessService(dbClient)

		processes, err := hiringProcessService.ListHiringProcesses(c.Request.Context(), userIDs)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		var response []SuggestionsResponse
		for _, process := range processes {
			response = append(response, SuggestionsResponse{
				Id:   process.ID,
				Name: process.Title,
			})
		}

		c.JSON(http.StatusOK, response)
	}
}

// HiringProcessList godoc
// @Summary List hiring processes
// @Schemes
// @Description Return a list of hiring processes with id and title
// @Tags hiring-process
// @Accept json
// @Produce json
// @Success 200 {array} SuggestionsResponse
// @Router /suggestions/vacancies [post]
func VacancyList(
	dwClient *ent.Client,
) func(c *gin.Context) {
	return func(c *gin.Context) {
		var processesIds []int
		if err := c.ShouldBindJSON(&processesIds); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		vacancyService := service.NewVacancyService(dwClient)

		vacancies, err := vacancyService.GetVacancySuggestions(c.Request.Context(), processesIds)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		var response []SuggestionsResponse
		for _, vacancy := range vacancies {
			response = append(response, SuggestionsResponse{
				Id:   vacancy.ID,
				Name: vacancy.Title,
			})
		}

		c.JSON(http.StatusOK, response)
	}
}

// HiringProcessList godoc
// @Summary List hiring processes
// @Schemes
// @Description Return a list of hiring processes with id and title
// @Tags hiring-process
// @Accept json
// @Produce json
// @Success 200 {array} SuggestionsResponse
// @Router /suggestions/vacancies [post]
func VacancyTable(
	dwClient *ent.Client,
) func(c *gin.Context) {
	return func(c *gin.Context) {

		var filter service.VacancyTableFilter
		if err := c.ShouldBindJSON(&filter); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		vacancyServiceTablePointer := service.NewVacancyServiceTable(dwClient)

		vacancies, err := vacancyServiceTablePointer.GetVacancyTable(c.Request.Context(), filter)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		fmt.Println("Vagas retornadas:", vacancies)

		var response []TableResponse
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
			response = append(response, TableResponse{
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

		c.JSON(http.StatusOK, response)
	}
}
