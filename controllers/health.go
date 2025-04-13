package controllers

import (
	"net/http"
	"striplex/db"

	"github.com/gin-gonic/gin"
)

type HealthController struct{}

type HealthResponse struct {
	Status string `json:"status"`
}

func (h HealthController) Status(c *gin.Context) {
	var result int
	err := db.DB.Raw("SELECT 1").Scan(&result).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, HealthResponse{
			Status: "unhealthy",
		})
	} else {
		c.JSON(http.StatusOK, HealthResponse{
			Status: "healthy",
		})
	}
}
