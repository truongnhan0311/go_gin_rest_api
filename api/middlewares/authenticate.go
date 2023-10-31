package middlewares

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"warranty-api/controllers"
	"warranty-api/models"
	"warranty-api/services"
)

func responseWithError(c *gin.Context, code int, message interface{}) {
	c.AbortWithStatusJSON(code, gin.H{"message": message})
}

func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		requiredToken := c.Request.Header["Authorization"]

		if len(requiredToken) == 0 {
			responseWithError(c, http.StatusForbidden, "Please login to your account")
		}

		userID, _ := services.DecodeToken(requiredToken[0])
		UserController := controllers.UserController()
		user, err := models.GetUserByEmail(UserController.Db, userID)

		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong giving you access"})
			return
		}
		//print(user.Email)
		if user.Email == "" {
			responseWithError(c, http.StatusNotFound, "User account not found")
			return
		}

		c.Set("User", user)
		c.Set("user-email", user.Email)
		c.Set("user-id", user.ID)

		if user.AgencyID == nil {
			c.Set("user-agency", nil)
		} else {
			c.Set("user-agency", *user.AgencyID)
		}

		if user.DailyID == nil {
			c.Set("user-daily", nil)
		} else {
			c.Set("user-daily", *user.DailyID)
		}

		c.Next()
	}
}
