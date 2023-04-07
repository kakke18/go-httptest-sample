package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type LoginParam struct {
	UserName string `json:"user_name"`
	Password string `json:"password"`
}

var (
	validUserName = "user"
	validPassword = "password"
)

const (
	port = 8081
)

func main() {
	r := gin.Default()

	r.POST("/login", func(c *gin.Context) {
		var param LoginParam
		if err := c.ShouldBindJSON(&param); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request"})
			return
		}

		if param.UserName != validUserName || param.Password != validPassword {
			c.JSON(http.StatusUnauthorized, gin.H{"status": "Unauthorized"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "Authenticated"})
	})

	if err := r.Run(fmt.Sprintf(":%d", port)); err != nil {
		panic(err)
	}
}
