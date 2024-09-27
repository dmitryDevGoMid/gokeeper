package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/dmitryDevGoMid/gokeeper/server/internal/stuffing/model"
	"github.com/dmitryDevGoMid/gokeeper/server/internal/stuffing/repository/user"

	"github.com/gin-gonic/gin"
)

func (h *Handler) CardsList(c *gin.Context) {
	//Перменная в которой будут данные тело запроса и данные клиента
	var request model.RequestBody
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userRequest := request.User

	var list *[]user.ResponseSaveData
	list, err := h.repoUser.GetCardsByUser(c, &userRequest)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	listCardsJson, err := json.Marshal(&list)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"marshal error": err.Error()})
		return
	}

	type ResponseLists struct {
		ID   string `json:"id"`
		Data string `json:"data"`
	}

	respList := &[]ResponseLists{}

	err = json.Unmarshal(listCardsJson, respList)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"umarshal error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, *respList)
}
