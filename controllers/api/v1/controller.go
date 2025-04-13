package v1controller

import (
	"fmt"
	"net/http"
	"striplex/services"

	"github.com/gin-gonic/gin"
)

type V1 struct {
	basePath      string
	client        *http.Client
	wizarrService *services.WizarrService
}

func NewV1Controller(basePath string, client *http.Client, wizarrService *services.WizarrService) *V1 {
	return &V1{
		basePath:      basePath,
		client:        client,
		wizarrService: wizarrService,
	}
}
func (v *V1) GetRoutes(r *gin.RouterGroup) {
	stripe := r.Group("/stripe")
	{
		stripe.POST("/webhook", v.Webhook)
	}
}

func (s *V1) Webhook(c *gin.Context) {
	// Generate invite link using wizarr service
	link, err := s.wizarrService.GenerateInviteLink(0, 1)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "", gin.H{
			"error": "Failed to generate invite link",
		})
		return
	}

	// Create HTML content with the invite link
	htmlContent := fmt.Sprintf(`
	<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title>Your Wizarr Invite</title>
		<style>
			body {
				font-family: Arial, sans-serif;
				line-height: 1.6;
				max-width: 600px;
				margin: 0 auto;
				padding: 20px;
				text-align: center;
			}
			.container {
				border: 1px solid #ddd;
				border-radius: 8px;
				padding: 20px;
				box-shadow: 0 2px 8px rgba(0,0,0,0.1);
			}
			.invite-link {
				margin: 20px 0;
				padding: 12px;
				background-color: #f5f5f5;
				border-radius: 4px;
				word-break: break-all;
			}
			.button {
				display: inline-block;
				background-color: #4CAF50;
				color: white;
				padding: 10px 20px;
				text-decoration: none;
				border-radius: 4px;
				font-weight: bold;
			}
		</style>
	</head>
	<body>
		<div class="container">
			<h1>Your Wizarr Invite is Ready!</h1>
			<p>Click the button below to access your invite:</p>
			<div class="invite-link">
				<a href="%s" class="button">Access Your Invite</a>
			</div>
			<p>Or use this link:</p>
			<div class="invite-link">%s</div>
		</div>
	</body>
	</html>
	`, link, link)

	// Set content type and return the HTML
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(http.StatusOK, htmlContent)
}
