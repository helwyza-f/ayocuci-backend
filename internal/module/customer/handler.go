package customer

import (
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

// ========================
// CREATE CUSTOMER
// ========================
func (h *Handler) Create(c *gin.Context) {
	var req struct {
		Name  string `json:"name"`
		Phone string `json:"phone"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	outletID := c.GetUint("outlet_id")
	if outletID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "outlet context required"})
		return
	}

	if err := h.service.Create(outletID, req.Name, req.Phone); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "customer created"})
}

// ========================
// LIST CUSTOMERS
// ========================
func (h *Handler) List(c *gin.Context) {
	outletID := c.GetUint("outlet_id")
	if outletID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "outlet context required"})
		return
	}

	customers, err := h.service.GetAll(outletID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": customers})
}

// ========================
// GET CUSTOMER DETAIL
// ========================
func (h *Handler) GetByID(c *gin.Context) {
	outletID := c.GetUint("outlet_id")
	if outletID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "outlet context required"})
		return
	}

	id, _ := strconv.Atoi(c.Param("id"))
	customer, err := h.service.GetByID(outletID, uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": customer})
}

// ========================
// UPDATE CUSTOMER
// ========================
func (h *Handler) Update(c *gin.Context) {
	outletID := c.GetUint("outlet_id")
	if outletID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "outlet context required"})
		return
	}

	id, _ := strconv.Atoi(c.Param("id"))

	var req struct {
		Name  string `json:"name"`
		Phone string `json:"phone"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.Update(outletID, uint(id), req.Name, req.Phone); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "customer updated"})
}
