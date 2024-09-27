package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/dmitryDevGoMid/gokeeper/server/internal/stuffing/model"
	"github.com/dmitryDevGoMid/gokeeper/server/internal/stuffing/repository/files"

	"github.com/gin-gonic/gin"
)

func (h *Handler) DeleteFile(c *gin.Context) {

	//Перменная в которой будут данные тело запроса и данные клиента
	var request model.RequestBody
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//Переменная в которую положим прилетевшие данные
	fileData := &files.Files{}

	//Кладем json в структуру
	json.Unmarshal(request.Body, &fileData)

	if fileData.ID_File != "" {
		if err := h.repoFile.DeleteFilesByID(c, fileData.ID_File); err != nil {
			//c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			//return
			switch err {
			case files.ErrDataNotFound:
				c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			}
			return
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is empty"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": "delete is successful"})
}
