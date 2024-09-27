package handlers

import (
	"encoding/base64"
	"encoding/json"
	"net/http"

	"github.com/dmitryDevGoMid/gokeeper/server/internal/stuffing/model"
	"github.com/dmitryDevGoMid/gokeeper/server/internal/stuffing/repository/user"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (h *Handler) PasswordSave(c *gin.Context) {
	update := false
	// Variable to hold the request body and client data
	var request model.RequestBody
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Variable to hold the incoming data
	newpassword := &PasswordRequest{}

	// Unmarshal JSON into the structure
	json.Unmarshal(request.Body, &newpassword)

	// Structure for saving to the database, the same for all, only the type differs
	savePasswordByData := &user.SaveData{}

	savePasswordByData.TypeData = "password"

	if newpassword.ID != "" {
		update = true
		docID, err := primitive.ObjectIDFromHex(newpassword.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		savePasswordByData.ID = docID
	}

	// Get the client's public key and encrypt the data
	data, err := h.asimencrypt.EncryptByClientKeyParts(string(request.Body), request.User.PublicKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Convert the encrypted data to base64 to save in the database
	savePasswordByData.Data = base64.StdEncoding.EncodeToString(data)
	savePasswordByData.User_ID = request.User.ID_User

	if !update {
		if err := h.repoUser.CreatePasswordByUser(c, savePasswordByData); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	} else {
		if err := h.repoUser.UpdatePasswordByKey(c, savePasswordByData); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	if !update {
		c.JSON(http.StatusCreated, gin.H{"result": "create is successfull"})
	} else {
		c.JSON(http.StatusOK, gin.H{"result": "update is successfull"})
	}
}
