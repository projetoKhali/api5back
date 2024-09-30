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
		eg := v1.Group("/hiring-process")
		{
			eg.GET("/dashboard", Dashboard(dbClient, dwClient))
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

		metricsData, err := MetricsService.GetMetrics(c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		c.JSON(http.StatusOK, metricsData)
	}
}
