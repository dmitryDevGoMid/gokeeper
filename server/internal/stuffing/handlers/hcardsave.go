package handlers

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"

	"github.com/dmitryDevGoMid/gokeeper/server/internal/stuffing/model"
	"github.com/dmitryDevGoMid/gokeeper/server/internal/stuffing/repository/user"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (h *Handler) CardSave(c *gin.Context) {
	update := false
	// Variable to hold the request body and client data
	var request model.RequestBody
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Variable to hold the incoming data
	newcard := &CardRequest{}

	// Unmarshal JSON into the structure
	json.Unmarshal(request.Body, &newcard)

	// Structure for saving to the database, the same for all, only the type differs
	saveCardByData := &user.SaveData{}

	saveCardByData.TypeData = "card"

	if newcard.ID != "" {
		update = true
		docID, err := primitive.ObjectIDFromHex(newcard.ID)
		if err != nil {
			log.Fatal(err)
		}
		saveCardByData.ID = docID
	}

	// Get the client's public key and encrypt the data
	data, err := h.asimencrypt.EncryptByClientKeyParts(string(request.Body), request.User.PublicKey)
	if err != nil {
		log.Println("asimencrypt failed to encrypt", err)
	}

	// Convert the encrypted data to base64 to save in the database
	saveCardByData.Data = base64.StdEncoding.EncodeToString(data)
	saveCardByData.User_ID = request.User.ID_User

	if !update {
		if err := h.repoUser.CreateCardByUser(c, saveCardByData); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	} else {
		if err := h.repoUser.UpdateCardByKey(c, saveCardByData); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	if !update {
		c.JSON(http.StatusCreated, gin.H{"result": "create is successful"})
	} else {
		c.JSON(http.StatusOK, gin.H{"result": "update is successful"})
	}
}
