package controllers

import (
	"encoding/json"
	"net/http"
	"striplex/db"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type AppController struct{}

func NewAppController() *AppController {
	return &AppController{}
}

// WhoAmI returns the authenticated user's information from the session
func (p AppController) WhoAmI(c *gin.Context) {
	userInfo := sessions.Default(c).Get("user_info")
	if userInfo == nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": "error",
			"error":  "Not authenticated",
		})
		return
	}
	var plexUser PlexUserResponse
	err := json.Unmarshal(userInfo.([]byte), &plexUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"error":  "Failed to parse user info: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"user":   plexUser,
	})
}

func (h AppController) Health(c *gin.Context) {
	var result int
	err := db.DB.Raw("SELECT 1").Scan(&result).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "unhealthy",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status": "healthy",
		})
	}
}

// Index returns a styled HTML landing page for Striplex
func (a *AppController) Index(c *gin.Context) {
	html := `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Striplex</title>
    <style>
        body {
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
            background-color: #1e1e2e;
            color: #cdd6f4;
            margin: 0;
            padding: 0;
            display: flex;
            flex-direction: column;
            justify-content: center;
            align-items: center;
            height: 100vh;
            overflow: hidden;
        }
        .container {
            text-align: center;
            padding: 2rem;
            border-radius: 10px;
            background-color: #313244;
            box-shadow: 0 4px 20px rgba(0, 0, 0, 0.3);
            max-width: 800px;
            width: 100%;
        }
        h1 {
            font-size: 4rem;
            margin-bottom: 0.5rem;
            background: linear-gradient(90deg, #f38ba8, #fab387, #f9e2af, #a6e3a1, #74c7ec, #cba6f7);
            -webkit-background-clip: text;
            background-clip: text;
            -webkit-text-fill-color: transparent;
            animation: gradient 10s ease infinite;
            background-size: 400% 400%;
        }
        p {
            font-size: 1.2rem;
            margin-bottom: 2rem;
            color: #bac2de;
        }
        .logo-container {
            display: flex;
            justify-content: center;
            align-items: center;
            margin-bottom: 2rem;
        }
        .logo {
            max-width: 200px;
        }
        .subtitle {
            display: inline-block;
            padding: 0.5rem 1rem;
            background-color: #45475a;
            border-radius: 10px;
            color: #cdd6f4;
            font-weight: bold;
        }
        @keyframes gradient {
            0% {
                background-position: 0% 50%;
            }
            50% {
                background-position: 100% 50%;
            }
            100% {
                background-position: 0% 50%;
            }
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>Striplex</h1>
        <div class="subtitle">Combining Stripe and Plex for seamless subscription management</div>
    </div>
</body>
</html>
`
	c.Header("Content-Type", "text/html")
	c.String(http.StatusOK, html)
}
