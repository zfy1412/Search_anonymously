package send

import (

	"github.com/gin-gonic/gin"

)

func DataAuthority() gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.Request.Header.Get("token")
		if key == "" {
			//c.JSON(http.StatusBadRequest, gin.H{})
			c.Abort()
			return
		}
		claims, err := ParseToken(key)
		if err != nil {
			//c.JSON(http.StatusBadRequest, gin.H{})
			c.Abort()
			return
		}
		if Ask(claims.User.Id,claims.User.Name,claims.User.Password)==1 {
			//c.JSON(http.StatusBadRequest, gin.H{})
			c.Abort()
			return
		}
		c.Next()
	}
}
