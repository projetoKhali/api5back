package server

import (
	"net/http"

	"api5back/ent"

	"github.com/gin-gonic/gin"
)

type endpointGroup func(
	engine *gin.Engine,
	dbClient *ent.Client,
	dwClient *ent.Client,
)

func NewServer(
	dbClient *ent.Client,
	dwClient *ent.Client,
) *gin.Engine {
	engine := gin.Default()

	for _, endpointGroups := range []endpointGroup{
		Swagger,
		Base,
		HiringProcessDashboard,
	} {
		endpointGroups(
			engine,
			dbClient,
			dwClient,
		)
	}

	return engine
}

// @BasePath /api/v1
func Base(
	engine *gin.Engine,
	dbClient *ent.Client,
	dwClient *ent.Client,
) {
	v1 := engine.Group("/api/v1")
	{
		eg := v1.Group("/example")
		{
			eg.GET("/helloworld", Helloworld(dbClient, dwClient))
		}
	}
}

// PingExample godoc
// @Summary ping example
// @Schemes
// @Description do ping
// @Tags example
// @Accept json
// @Produce json
// @Success 200 {string} Helloworld
// @Router /example/helloworld [get]
func Helloworld(
	dbClient *ent.Client,
	dwClient *ent.Client,
) func(c *gin.Context) {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, "helloworld")
	}
}
