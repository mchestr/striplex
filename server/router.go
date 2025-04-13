package server

import (
	"striplex/config"
	"striplex/controllers"

	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	gin.SetMode(config.GetConfig().GetString("server.mode"))
	router := gin.Default()

	health := new(controllers.HealthController)

	router.GET("/health", health.Status)
	return router

}
