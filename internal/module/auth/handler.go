// internal/module/auth/handler.go
package auth

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/helwyza-f/ayocuci-backend/internal/common"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service}
}

func (h *Handler) Register(c *gin.Context) {
	var req struct {
		BusinessName string `json:"business_name"`
		Email        string `json:"email"`
		Password     string `json:"password"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.RegisterOwner(
		req.BusinessName,
		req.Email,
		req.Password,
	); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "owner registered",
	})
}

func (h *Handler) Login(c *gin.Context) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.service.Login(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	token, err := common.GenerateJWT(
		user.ID,
		user.ClientID,
		user.Role,
		os.Getenv("JWT_SECRET"),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user_id":   user.ID,
		"client_id": user.ClientID,
		"role":      user.Role,
	})
}

