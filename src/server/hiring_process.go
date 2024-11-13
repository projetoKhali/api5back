package server

import (
	"net/http"

	"api5back/ent"
	"api5back/src/model"
	"api5back/src/service"

	"github.com/gin-gonic/gin"
)

func HiringProcessDashboard(
	engine *gin.Engine,
	dbClient *ent.Client,
	dwClient *ent.Client,
) {
	v1 := engine.Group("/api/v1")
	{
		hiringProcess := v1.Group("/hiring-process")
		{
			hiringProcess.POST("/dashboard", Dashboard(dwClient))
			hiringProcess.POST("/table", VacancyTable(dwClient))
		}

		suggestions := v1.Group("/suggestions")
		{
			suggestions.POST("/recruiter", UserList(dwClient))
			suggestions.POST("/process", HiringProcessList((dwClient)))
			suggestions.POST("/vacancy", VacancyList(dwClient))
		}
	}
}

// Dashboard godoc
// @Summary dashboard
// @Schemes
// @Description show dashboard
// @Tags hiring-process
// @Accept json
// @Param body body model.FactHiringProcessFilter true "Metrics filter"
// @Produce json
// @Success 200 {string} Dashboard
// @Router /hiring-process/dashboard [post]
func Dashboard(
	dwClient *ent.Client,
) func(c *gin.Context) {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")

		// TODO: change to pointer
		var dashboardMetricsFilter model.FactHiringProcessFilter
		if err := c.ShouldBindJSON(&dashboardMetricsFilter); err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}

		metricsData, err := service.GetMetrics(
			c, dwClient, dashboardMetricsFilter,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		c.JSON(http.StatusOK, metricsData)
	}
}

// UserList godoc
// @Summary List users
// @Schemes
// @Description Return a list of users with id and name
// @Tags suggestions
// @Accept json
// @Produce json
// @Success 200 {array} model.Suggestion
// @Router /suggestions/recruiter/ [get]
func UserList(dwClient *ent.Client) func(c *gin.Context) {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		users, err := service.GetUsers(c, dwClient)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		c.JSON(http.StatusOK, users)
	}
}

// HiringProcessList godoc
// @Summary List hiring processes
// @Schemes
// @Description Return a list of hiring processes with id and title
// @Tags suggestions
// @Accept json
// @Param body body []int true "User IDs"
// @Produce json
// @Success 200 {array} model.Suggestion
// @Router /suggestions/process [post]
func HiringProcessList(
	dbClient *ent.Client,
) func(c *gin.Context) {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		var userIDs *[]int

		// Parse the body for user IDs
		if err := c.ShouldBindJSON(&userIDs); err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}

		processes, err := service.ListHiringProcesses(
			c, dbClient,
			userIDs,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		c.JSON(http.StatusOK, processes)
	}
}

// HiringProcessList godoc
// @Summary List hiring processes
// @Schemes
// @Description Return a list of hiring processes with id and title
// @Tags suggestions
// @Accept json
// @Param body body []int false "User IDs"
// @Produce json
// @Success 200 {array} model.Suggestion
// @Router /suggestions/vacancy [post]
func VacancyList(
	dwClient *ent.Client,
) func(c *gin.Context) {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		var processesIds *[]int
		if err := c.ShouldBindJSON(&processesIds); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		vacancies, err := service.GetVacancySuggestions(
			c, dwClient,
			processesIds,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, vacancies)
	}
}

// HiringProcessList godoc
// @Summary List hiring processes
// @Schemes
// @Description Return a list of vacancies with summarized information
// @Tags hiring-process
// @Accept json
// @Param body body model.FactHiringProcessFilter true "Metrics filter"
// @Produce json
// @Success 200 {array} model.Suggestion
// @Router /hiring-process/table [post]
func VacancyTable(
	dwClient *ent.Client,
) func(c *gin.Context) {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		var filter model.FactHiringProcessFilter
		if err := c.ShouldBindJSON(&filter); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		vacancies, err := service.GetVacancyTable(
			c, dwClient,
			filter,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, vacancies)
	}
}
