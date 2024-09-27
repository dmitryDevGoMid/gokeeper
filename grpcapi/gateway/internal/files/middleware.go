package files

import (
	"fmt"
	"net/http"

	"github.com/dmitryDevGoMid/gokeeper/grpcapi/gateway/internal/config"

	"github.com/gin-gonic/gin"
)

func CheckAuthByToken(cfg *config.Config) gin.HandlerFunc {

	return func(c *gin.Context) {

		token := c.GetHeader("Token")

		if token == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error"})
			c.Abort()
		}

		urlServerAuth := fmt.Sprintf("http://%s/api/user/check/token", cfg.GoKeeperServerAdress.GoKeeperAdress)

		// Выполняем HTTP-запрос к серверу авторизации с токеном
		req, err := http.NewRequest("POST", urlServerAuth, nil)
		req.Header.Set("Token", token)

		if err != nil {
			fmt.Println(urlServerAuth, err)
		}

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer resp.Body.Close()

		// Проверяем статус ответа
		if resp.StatusCode == http.StatusOK {
			c.Next()
		} else if resp.StatusCode == http.StatusUnauthorized {
			c.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
			c.Abort()
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error"})
			c.Abort()
		}

		c.Next()
	}
}
