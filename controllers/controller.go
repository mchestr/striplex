package controllers

import (
	"encoding/json"
	"net/http"
	"striplex/config"
	"striplex/db"
	"striplex/model"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"

	apicontroller "striplex/controllers/api"
	plexcontroller "striplex/controllers/plex"
	stripecontroller "striplex/controllers/stripe"
)

type AppController struct {
	client *http.Client
}

func NewAppController(client *http.Client) *AppController {
	return &AppController{
		client: client,
	}
}

func (c *AppController) GetRoutes(r *gin.RouterGroup) {
	client := http.Client{
		Timeout: 10 * time.Second,
	}

	r.GET("/health", c.Health)
	r.GET("/whoami", c.WhoAmI)
	r.GET("/logout", c.Logout)
	r.GET("/", c.Index)

	api := r.Group("/api")
	{
		apiController := apicontroller.NewApiController(api.BasePath(), &client)
		apiController.GetRoutes(api)
	}

	plex := r.Group("/plex")
	{
		plexController := plexcontroller.NewPlexController(plex.BasePath(), &client)
		plexController.GetRoutes(plex)
	}

	stripe := r.Group("/stripe")
	{
		stripeController := stripecontroller.NewStripeController(stripe.BasePath(), &client)
		stripeController.GetRoutes(stripe)
	}
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
	var userInfoData model.UserInfo
	err := json.Unmarshal(userInfo.([]byte), &userInfoData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"error":  "Failed to parse user info: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"user":   userInfoData,
	})
}

// Logout clears the user session by deleting the user_info key
func (p AppController) Logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Delete("user_info")
	err := session.Save()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"error":  "Failed to save session: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Successfully logged out",
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
	// Get the default price ID from configuration
	priceID := config.Config.GetString("stripe.default_price_id")

	// Check if user is authenticated
	session := sessions.Default(c)
	userInfo := session.Get("user_info")
	isAuthenticated := userInfo != nil

	// Base HTML with styles
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
            margin-bottom: 2rem;
        }
        .checkout-btn {
            background: linear-gradient(90deg, #f38ba8, #fab387);
            color: #1e1e2e;
            border: none;
            padding: 0.8rem 1.5rem;
            border-radius: 8px;
            font-size: 1.2rem;
            font-weight: bold;
            cursor: pointer;
            transition: transform 0.2s, box-shadow 0.2s;
            box-shadow: 0 2px 10px rgba(243, 139, 168, 0.4);
            margin-right: 1rem;
        }
        .checkout-btn:hover {
            transform: translateY(-2px);
            box-shadow: 0 4px 15px rgba(243, 139, 168, 0.6);
        }
        .logout-btn {
            background: linear-gradient(90deg, #74c7ec, #89dceb);
            color: #1e1e2e;
            border: none;
            padding: 0.8rem 1.5rem;
            border-radius: 8px;
            font-size: 1.2rem;
            font-weight: bold;
            cursor: pointer;
            transition: transform 0.2s, box-shadow 0.2s;
            box-shadow: 0 2px 10px rgba(116, 199, 236, 0.4);
        }
        .logout-btn:hover {
            transform: translateY(-2px);
            box-shadow: 0 4px 15px rgba(116, 199, 236, 0.6);
        }
        .button-container {
            display: flex;
            justify-content: center;
            gap: 1rem;
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
        <div class="button-container">
`

	// Conditionally add buttons based on authentication status
	if isAuthenticated {
		// User is authenticated, show both buttons
		html += `
            <button class="checkout-btn" onclick="window.location.href='/stripe/checkout?price_id=` + priceID + `'">Subscribe Now</button>
            <button class="logout-btn" onclick="window.location.href='/logout'">Logout</button>
`
	} else {
		// User is not authenticated, show only sign in button
		html += `
            <button class="checkout-btn" onclick="window.location.href='/stripe/checkout?price_id=` + priceID + `'">Sign in to Subscribe</button>
`
	}

	// Close the HTML
	html += `
        </div>
    </div>
</body>
</html>
`

	c.Header("Content-Type", "text/html")
	c.String(http.StatusOK, html)
}
