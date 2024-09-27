package handlers

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"

	"github.com/dmitryDevGoMid/gokeeper/server/internal/stuffing/model"
	"github.com/dmitryDevGoMid/gokeeper/server/internal/stuffing/repository/files"

	"github.com/gin-gonic/gin"
)

func (h *Handler) FilesList(c *gin.Context) {
	//Перменная в которой будут данные тело запроса и данные клиента
	var request model.RequestBody
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userRequest := request.User

	var list *[]files.Files
	list, err := h.repoFile.GetByUserIdListFiles(c, &userRequest)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	type ResponseLists struct {
		ID   string `json:"id"`
		Data string `json:"data"`
	}

	respList := []ResponseLists{}

	//Проходимся по полученным данным в цикле
	for _, val := range *list {
		//Преобразуем данные в json
		valJson, err := json.Marshal(val)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"marshal error": err.Error()})
			return
		}
		//Шифруем данные публичным ключом клиента
		data, err := h.asimencrypt.EncryptByClientKeyParts(string(valJson), request.User.PublicKey)
		if err != nil {
			log.Println("asimencrypt failed to encrypt", err)
		}
		//Сохраяем данные в массив и преобразуем в base64
		respList = append(respList, ResponseLists{ID: val.ID.Hex(), Data: base64.StdEncoding.EncodeToString(data)})
	}

	c.JSON(http.StatusOK, respList)
}
