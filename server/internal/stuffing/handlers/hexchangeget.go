package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) ExchangeGet(c *gin.Context) {

	type Exchange struct {
		Type string `json:"type"`
	}
	exchange := &Exchange{}

	if err := c.ShouldBindJSON(exchange); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fmt.Println(exchange)

	type ResponseExchange struct {
		Key  string `json:"key"`
		Type string `json:"type"`
	}

	responseExchange := &ResponseExchange{Key: string(h.asimencrypt.GetBytePublic()), Type: "get_public_key"}

	c.JSON(http.StatusCreated, &responseExchange)
}
