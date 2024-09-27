package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/dmitryDevGoMid/gokeeper/server/internal/stuffing/model"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (h *Handler) CardDelete(c *gin.Context) {
	var docID primitive.ObjectID
	var err error

	//Перменная в которой будут данные тело запроса и данные клиента
	var request model.RequestBody
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//Переменная в которую положим прилетевшие данные
	card := &CardRequest{}

	//Кладем json в структуру
	json.Unmarshal(request.Body, &card)

	if card.ID != "" {
		docID, err = primitive.ObjectIDFromHex(card.ID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "id not object"})
			return
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is empty"})
		return
	}

	if err := h.repoUser.DelereCardById(c, docID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": "delete is successful"})
}
