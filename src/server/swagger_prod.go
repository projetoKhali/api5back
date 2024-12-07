//go:build production || integration
// +build production integration

package server

import (
	"api5back/ent"

	"github.com/gin-gonic/gin"
)

func Swagger(
	engine *gin.Engine,
	dbClient *ent.Client,
	dwClient *ent.Client,
) {

}
