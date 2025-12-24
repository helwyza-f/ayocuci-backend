package employee

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service}
}

func (h *Handler) CreateUser(c *gin.Context) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	clientID := c.GetUint("client_id")

	userID, err := h.service.CreateUser(clientID, req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"user_id": userID,
	})
}

func (h *Handler) Assign(c *gin.Context) {
	var req struct {
		UserID uint   `json:"user_id"`
		Role   string `json:"role"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	outletID, _ := strconv.Atoi(c.Param("outlet_id"))
	requesterClientID := c.GetUint("client_id")

	log.Printf("Requester Client ID: %d", requesterClientID)

	if err := h.service.Assign(
		requesterClientID,
		req.UserID,
		uint(outletID),
		req.Role,
	); err != nil {
		c.JSON(403, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "employee assigned"})
}

func (h *Handler) Transfer(c *gin.Context) {
    var req struct {
        TargetOutletID uint   `json:"target_outlet_id"`
        Role           string `json:"role"`
    }

    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }

    userID, _ := strconv.Atoi(c.Param("user_id"))
    requesterClientID := c.GetUint("client_id")

    // Kita buat service baru: Transfer
    if err := h.service.Transfer(requesterClientID, uint(userID), req.TargetOutletID, req.Role); err != nil {
        c.JSON(403, gin.H{"error": err.Error()})
        return
    }

    c.JSON(200, gin.H{"message": "employee transferred successfully"})
}

func (h *Handler) ListEmployees(c *gin.Context) {
	outletID := c.GetUint("outlet_id")

	employees, err := h.service.ListEmployees(outletID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": employees})
}

// 🔥 ENDPOINT PENTING UNTUK PILIH OUTLET AKTIF
func (h *Handler) MyOutlets(c *gin.Context) {
	role := c.GetString("role")
	userID := c.GetUint("user_id")
	clientID := c.GetUint("client_id")

	outlets, err := h.service.MyOutlets(role, userID, clientID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": outlets})
}

