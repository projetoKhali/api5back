//go:build !production && !integration
// +build !production,!integration

package server

import (
	docs "api5back/docs"
	"api5back/ent"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Swagger(
	engine *gin.Engine,
	dbClient *ent.Client,
	dwClient *ent.Client,
) {
	docs.SwaggerInfo.BasePath = "/api/v1"

	engine.GET(
		"/swagger/*any",
		ginSwagger.WrapHandler(swaggerfiles.Handler),
	)
}
