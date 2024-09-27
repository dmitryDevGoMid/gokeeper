package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/dmitryDevGoMid/gokeeper/server/internal/stuffing/model"
	"github.com/dmitryDevGoMid/gokeeper/server/internal/stuffing/repository/user"

	"github.com/gin-gonic/gin"
)

func (h *Handler) PasswordList(c *gin.Context) {
	//Перменная в которой будут данные тело запроса и данные клиента
	var request model.RequestBody
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userRequest := request.User

	var list *[]user.ResponseSaveData

	list, err := h.repoUser.GetPasswordByUser(c, &userRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	listPasswordsJson, err := json.Marshal(&list)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"marshal error": err.Error()})
		return
	}

	type ResponseLists struct {
		ID   string `json:"id"`
		Data string `json:"data"`
	}

	respList := &[]ResponseLists{}

	err = json.Unmarshal(listPasswordsJson, respList)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	/*if err != nil {
		fmt.Println("json.Unmarshal(listPasswordsJson, &ResponseLists{}) =>", err)
	} else {
		fmt.Println("respList=====>>>>>>>>", respList)
	}*/

	//fmt.Println("List Passwords =>", list)
	//fmt.Println("List Passwords JSON =>", string(listPasswordsJson))

	/*data, err := h.asimencrypt.EncryptByClientKeyParts(string(listPasswordsJson), request.User.PublicKey)
	if err != nil {
		log.Println("asimencrypt failed to encrypt", err)
	}
	fmt.Println("LIST PASSWORDS:", base64.StdEncoding.EncodeToString(data))*/

	//c.JSON(http.StatusOK, base64.StdEncoding.EncodeToString(data))
	c.JSON(http.StatusOK, *respList)
}
