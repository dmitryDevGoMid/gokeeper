package handlers

import (
	"errors"
	"net/http"

	userRepository "github.com/dmitryDevGoMid/gokeeper/server/internal/stuffing/repository/user"

	"github.com/gin-gonic/gin"
)

func (h *Handler) Register(c *gin.Context) {
	var user userRepository.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if user.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": errors.New("password is required  for login")})
		return
	}

	hashPassword, err := h.security.EncryptPassword(c, user.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user.Password = hashPassword

	if err := h.repoUser.Create(c, &user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, user)
}
