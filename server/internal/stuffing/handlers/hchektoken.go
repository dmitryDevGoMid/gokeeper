package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) CheckToken(c *gin.Context) {
	c.JSON(http.StatusOK, nil)
}
