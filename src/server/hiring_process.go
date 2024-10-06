package server

import (
	"net/http"

	"api5back/ent"
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
			hiringProcess.GET("/dashboard", Dashboard(dwClient))
		}

		suggestions := v1.Group("/suggestions")
		{
			suggestions.GET("/recruiter", UserList(dwClient))
			suggestions.POST("/process", HiringProcessList(dwClient))
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
	dwClient *ent.Client,
) func(c *gin.Context) {
	return func(c *gin.Context) {
		MetricsService := service.NewMetricsService(dwClient)

		hiringProcessName := c.Query("hiringProcess")
		vacancyName := c.Query("vacancy")
		startDate := c.Query("startDate")
		endDate := c.Query("endDate")

		metricsData, err := service.GetMetrics(
			c, dwClient,
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

// UserList godoc
// @Summary List users
// @Schemes
// @Description Return a list of users with id and name
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {array} model.Suggestion
// @Router /users/ [get]
func UserList(
	dwClient *ent.Client,
) func(c *gin.Context) {
	return func(c *gin.Context) {
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
// @Tags hiring-process
// @Accept json
// @Produce json
// @Success 200 {array} model.Suggestion
// @Router /hiring-process [post]
func HiringProcessList(
	dbClient *ent.Client,
) func(c *gin.Context) {
	return func(c *gin.Context) {
		var userIDs []int

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
