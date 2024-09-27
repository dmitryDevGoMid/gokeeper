package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Exchange struct {
	Key  string `json:"key"`
	Type string `json:"type"`
}

type ResponseExchange struct {
	Success bool   `json:"success"`
	Type    string `json:"type"`
}

func (h *Handler) ExchangeSet(c *gin.Context) {

	exchange := &Exchange{}

	if err := c.ShouldBindJSON(exchange); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fmt.Println(exchange)

	responseExchange := &ResponseExchange{Success: true, Type: "set_public_key"}

	c.JSON(http.StatusCreated, &responseExchange)
}