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

			fmt.Println("ENCRYPTE=====>>>>")
			fmt.Println("ENCRYPTE=====>>>>", string(body))

			decompressBody, _ := asimencrypt.DecryptOAEP(body)

			fmt.Println("decompressBody=====>>>>")
			fmt.Println("decompressBody=====>>>>", string(decompressBody))

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

			fmt.Println("user.PublicKey=======>", user.PublicKey)

			fmt.Println("UDERRSRSRSRSRS=>>>>", user)

			requestBody.User = *user
			body, _ := io.ReadAll(c.Request.Body)

			fmt.Println("ENCRYPTE=====>>>>")
			fmt.Println("ENCRYPTE=====>>>>", string(body))

			//decompressBody, _ := asimencrypt.DecryptOAEP(body)
			//requestBody.Body = decompressBody
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

/*func Authenticate(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//Выдергиваем токен из заголовка
		tokenString, err := security.ExtractToken(r)
		if err != nil {
			restutil.WriteError(w, http.StatusUnauthorized, restutil.ErrUnauthorized)
			return
		}
		//Парсим полученный токен
		token, err := security.ParseToken(tokenString)
		if err != nil {
			log.Println("parse error token", err.Error())
			//Возвращаем собщение со статусом не авторизован
			restutil.WriteError(w, http.StatusUnauthorized, restutil.ErrUnauthorized)
			return
		}

		if !token.Valid {
			log.Println("invalid token", tokenString)
			restutil.WriteError(w, http.StatusUnauthorized, restutil.ErrUnauthorized)
			return
		}

		next(w, r)

	}
}*/
