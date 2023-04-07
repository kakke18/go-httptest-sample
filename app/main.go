package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type LoginParam struct {
	UserName string `json:"user_name"`
	Password string `json:"password"`
}

const (
	port           = 8080
	authnServerURL = "http://authn:8081/login"
)

var (
	errInvalidParam   = fmt.Errorf("invalid param")
	errInternalServer = fmt.Errorf("internal server error")
)

func main() {
	r := gin.Default()

	r.GET("/ok", handleGetOk())
	r.POST("/login", handlePostLogin())

	if err := r.Run(fmt.Sprintf(":%d", port)); err != nil {
		panic(err)
	}
}

func handleGetOk() func(c *gin.Context) {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "OK"})
	}
}

type AuthnService interface {
	Authenticate(string, string) (bool, error)
}

func handlePostLogin() func(c *gin.Context) {
	return func(c *gin.Context) {
		var param LoginParam
		if err := c.ShouldBindJSON(&param); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request"})
			return
		}

		authnService := newAuthenticator()
		authnResult, err := authnService.Authenticate(param.UserName, param.Password)
		if err != nil {
			switch {
			case errors.Is(err, errInvalidParam):
				c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request"})
				return
			case errors.Is(err, errInternalServer):
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
				return
			default:
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
				return
			}
		}

		if !authnResult {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "OK"})
	}
}

func newAuthenticator() AuthnService {
	env := os.Getenv("ENV")
	if env == "test" {
		return &AuthnMock{}
	}

	return &Authn{}
}

type Authn struct{}

func (a *Authn) Authenticate(userName, password string) (bool, error) {
	reqBody, err := json.Marshal(map[string]string{
		"user_name": userName,
		"password":  password,
	})
	if err != nil {
		return false, err
	}

	resp, err := http.Post(authnServerURL, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if statusCode := resp.StatusCode; statusCode != http.StatusOK {
		switch statusCode {
		case http.StatusBadRequest:
			return false, errInvalidParam
		case http.StatusUnauthorized:
			return false, nil
		case http.StatusInternalServerError:
			return false, errInternalServer
		default:
			return false, errInternalServer
		}
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	var responseBody map[string]string
	if err = json.Unmarshal(body, &responseBody); err != nil {
		return false, err
	}

	status := responseBody["status"]

	return status == "Authenticated", nil
}

type AuthnMock struct{}

func (a *AuthnMock) Authenticate(userName, password string) (bool, error) {
	if userName != "user" || password != "password" {
		return false, nil
	}

	return true, nil
}
