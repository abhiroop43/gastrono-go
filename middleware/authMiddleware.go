package middleware

import "github.com/gin-gonic/gin"

func AuthenticationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientToken := c.Request.Header.Get("Authorization")

		if clientToken == "" {
			c.JSON(401, gin.H{"message": "Unauthorized"})
			c.Abort()
			return
		}

		//claims, err := helper
	}
}
