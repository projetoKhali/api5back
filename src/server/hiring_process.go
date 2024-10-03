package server

import (
	"api5back/ent"
	"api5back/src/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

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
			suggestions.GET("/process", HiringProcessList(dwClient))
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

// UserList godoc
// @Summary List users
// @Schemes
// @Description Return a list of users with id and name
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {array} map[string]interface{}
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
		var response []map[string]interface{}
		for _, user := range users {
			response = append(response, map[string]interface{}{
				"id":   user.ID,
				"name": user.Name,
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
// @Success 200 {array} map[string]interface{}
// @Router /hiring-process [get]
func HiringProcessList(
	dbClient *ent.Client,
) func(c *gin.Context) {
	return func(c *gin.Context) {
		hiringProcessService := service.NewHiringProcessService(dbClient)

		userIDParams := c.QueryArray("userID")
		var userIDs []int

		for _, userIDParam := range userIDParams {
			id, err := strconv.Atoi(userIDParam)
			if err == nil {
				userIDs = append(userIDs, id)
			}
		}

		processes, err := hiringProcessService.ListHiringProcesses(c.Request.Context(), userIDs)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		var response []map[string]interface{}
		for _, process := range processes {
			response = append(response, map[string]interface{}{
				"id":    process.ID,
				"title": process.Title,
			})
		}

		c.JSON(http.StatusOK, response)
	}
}
