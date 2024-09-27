package handlers

import (
	"net/http"

	"github.com/dmitryDevGoMid/gokeeper/server/internal/stuffing/config/db"
	"github.com/dmitryDevGoMid/gokeeper/server/internal/stuffing/pkg/asimencrypt"
	"github.com/dmitryDevGoMid/gokeeper/server/internal/stuffing/pkg/security"
	filesRepository "github.com/dmitryDevGoMid/gokeeper/server/internal/stuffing/repository/files"
	userRepository "github.com/dmitryDevGoMid/gokeeper/server/internal/stuffing/repository/user"

	"github.com/gin-gonic/gin"
)

type PasswordRequest struct {
	ID          string `json:"id,omitempty"`
	Description string `json:"description"`
	Username    string `json:"username"`
	Password    string `json:"password"`
}

type CardRequest struct {
	ID          string `json:"id"`
	Description string `json:"description"`
	Number      string `json:"login"`
	Exp         string `json:"exp"`
	Cvc         string `json:"cvc"`
}

type Handlers interface {
	Login(c *gin.Context)
	ExchangeGet(c *gin.Context)
}

type Handler struct {
	repoUser    userRepository.UserRepository
	repoFile    filesRepository.FilesRepository
	asimencrypt asimencrypt.AsimEncrypt
	security    security.ISecurity
}

func NewIHandler(mongodb db.MongoDBClient, asimencrypt asimencrypt.AsimEncrypt) Handlers {
	repoUser := userRepository.NewUserRepository(mongodb)
	repoFile := filesRepository.NewFilesRepository(mongodb)
	return &Handler{repoUser: repoUser, repoFile: repoFile, asimencrypt: asimencrypt}
}

func NewHandler(repoUser userRepository.UserRepository, repoFile filesRepository.FilesRepository, asimencrypt asimencrypt.AsimEncrypt, security security.ISecurity) *Handler {
	//repoUser := userRepository.NewUserRepository(mongodb)
	//repoFile := filesRepository.NewFilesRepository(mongodb)
	//return &Handler{repoUser: repoUser, repoFile: repoFile, asimencrypt: asimencrypt}
	return &Handler{repoUser: repoUser, repoFile: repoFile, asimencrypt: asimencrypt, security: security}
}

func (h *Handler) UpdatePassword(c *gin.Context) {
	username := c.Param("username")
	var input struct {
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.repoUser.UpdatePassword(c, username, input.Password); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "password updated"})
}

func (h *Handler) DeleteByUsername(c *gin.Context) {
	username := c.Param("username")

	if err := h.repoUser.DeleteByUsername(c, username); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user deleted"})
}
