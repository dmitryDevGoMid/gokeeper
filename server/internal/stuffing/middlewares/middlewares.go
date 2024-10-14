package middlewares

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/dmitryDevGoMid/gokeeper/server/internal/stuffing/model"
	"github.com/dmitryDevGoMid/gokeeper/server/internal/stuffing/pkg/asimencrypt"
	"github.com/dmitryDevGoMid/gokeeper/server/internal/stuffing/pkg/jwttoken"

	"github.com/gin-gonic/gin"
)

// Миделвари для заголовков
func CheckToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.FullPath()
		if path != "/api/user/refresh/token" && path != "/api/user/login" && path != "/api/user/exchange/get" && path != "/api/user/exchange/set" && path != "/api/user/register" {
			_, err := jwttoken.Welcom(c.Request.Header.Get("Token"))
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"Message": "Unauthorized"})
				c.Abort()
				return
			}
		}
		c.Next()
	}
}

// Расшифровываем данные кроме пути: "/api/user/exchange"
func DecryptMiddleware(asimencrypt asimencrypt.AsimEncrypt) gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("DecryptMiddleware")

		if c.FullPath() != "/api/user/exchange/get" && c.FullPath() != "/api/user/refresh/token" {

			body, _ := io.ReadAll(c.Request.Body)

			decompressBody, _ := asimencrypt.DecryptOAEP(body)

			c.Request.Body = io.NopCloser(bytes.NewReader([]byte(decompressBody)))
		}

		c.Next()

	}
}

func ChaeckAndGetUserByToken(asimencrypt asimencrypt.AsimEncrypt) gin.HandlerFunc {
	return func(c *gin.Context) {
		requestBody := model.RequestBody{}
		fmt.Println("CheckToken")
		path := c.FullPath()
		if path != "/api/user/refresh/token" && path != "/api/user/login" && path != "/api/user/exchange/get" && path != "/api/user/exchange/set" && path != "/api/user/register" {
			user, err := jwttoken.Welcom(c.Request.Header.Get("Token"))
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"Message": "Unauthorized"})
				c.Abort()
				return
			}

			publicKey, err := base64.StdEncoding.DecodeString(c.Request.Header.Get("Public-Key"))
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"Message": "Unauthorized"})
				c.Abort()
				return
			}
			user.PublicKey = string(publicKey)

			requestBody.User = *user
			body, _ := io.ReadAll(c.Request.Body)

			requestBody.Body = body

			changeRequest, err := json.Marshal(requestBody)

			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"Message": "Unauthorized"})
				c.Abort()
				return
			}

			c.Request.Body = io.NopCloser(bytes.NewReader(changeRequest))
		}
		c.Next()
	}
}
