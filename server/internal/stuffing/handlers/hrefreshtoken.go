package handlers

import (
	"net/http"

	"github.com/dmitryDevGoMid/gokeeper/server/internal/stuffing/pkg/jwttoken"

	"github.com/gin-gonic/gin"
)

func (h *Handler) RefreshToken(c *gin.Context) {
	token, err := jwttoken.RefreshJWT(c.GetHeader("Token-Refresh"))
	if err != nil {
		c.JSON(http.StatusUnauthorized, nil)
		return
	}

	c.Header("Token", token)
	c.Header("Token-Refresh", c.GetHeader("Token-Refresh"))

	c.JSON(http.StatusOK, nil)
}
