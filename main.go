package main

import (
	"api5back/src/database"
	"fmt"
	"net/http"

	docs "api5back/docs"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @BasePath /api/v1

// PingExample godoc
// @Summary ping example
// @Schemes
// @Description do ping
// @Tags example
// @Accept json
// @Produce json
// @Success 200 {string} Helloworld
// @Router /example/helloworld [get]
func Helloworld(g *gin.Context) {
	g.JSON(http.StatusOK, "helloworld")
}

func main() {
	dbClient, err := database.DatabaseSetup("DB")
	if err != nil {
		panic(fmt.Errorf("failed to setup normalized database: %v", err))
	}
	defer dbClient.Close()

	dwClient, err := database.DatabaseSetup("DW")
	if err != nil {
		panic(fmt.Errorf("failed to setup data warehouse: %v", err))
	}
	defer dwClient.Close()

	r := gin.Default()
	docs.SwaggerInfo.BasePath = "/api/v1"
	v1 := r.Group("/api/v1")
	{
		eg := v1.Group("/example")
		{
			eg.GET("/helloworld", Helloworld)
		}
	}
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	r.Run(":8080")
}
