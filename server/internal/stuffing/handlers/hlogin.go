package handlers

import (
	"fmt"
	"net/http"

	"github.com/dmitryDevGoMid/gokeeper/server/internal/stuffing/pkg/jwttoken"
	userRepository "github.com/dmitryDevGoMid/gokeeper/server/internal/stuffing/repository/user"

	"github.com/gin-gonic/gin"
)

func (h *Handler) Login(c *gin.Context) {

	// Declare a variable to hold the user data
	var user *userRepository.User

	// Bind the JSON request body to the user variable
	if err := c.ShouldBindJSON(&user); err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Retrieve the user from the repository by username
	userFind, err := h.repoUser.GetByUsername(c, user.Username)
	if err != nil {
		switch err {
		case userRepository.ErrDataNotFound:
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	// Verify the user's password
	err = h.security.VerifyPassword(c, userFind.Password, user.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, nil)
		return
	}

	// Generate a JWT token for the user
	token, err := jwttoken.SetToken(userFind)
	if err != nil {
		c.JSON(http.StatusUnauthorized, nil)
		return
	}

	// Generate a refresh token for the user
	tokenRefresh, err := jwttoken.CreateRefreshToken(userFind)
	if err != nil {
		c.JSON(http.StatusUnauthorized, nil)
		return
	}

	// Set the tokens in the response headers
	c.Header("Token", token)
	c.Header("TokenRefresh", tokenRefresh)

	// Update the user's tokens and clear the password field
	userFind.TokenRefresh = tokenRefresh
	userFind.Token = token
	userFind.Password = ""

	// Return the user data with tokens in the response body
	c.JSON(http.StatusOK, userFind)
}
