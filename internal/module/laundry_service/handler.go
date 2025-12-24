package laundry_service

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service ServiceLogic
}

func NewHandler(service ServiceLogic) *Handler {
	return &Handler{service}
}

// ========================
// CREATE SERVICE
// ========================
func (h *Handler) Create(c *gin.Context) {
	var req struct {
		Name     string `json:"name"`
		Price    int64  `json:"price"`
		Estimate string `json:"estimate"`
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

	if err := h.service.Create(outletID, req.Name, req.Price, req.Estimate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "service created"})
}

// ========================
// LIST SERVICES
// ========================
func (h *Handler) List(c *gin.Context) {
	outletID := c.GetUint("outlet_id")
	if outletID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "outlet context required"})
		return
	}

	services, err := h.service.GetAll(outletID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": services})
}

// ========================
// UPDATE SERVICE
// ========================
func (h *Handler) Update(c *gin.Context) {
	outletID := c.GetUint("outlet_id")
	if outletID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "outlet context required"})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid service id"})
		return
	}

	var req struct {
		Name     string `json:"name"`
		Price    int64  `json:"price"`
		Estimate string `json:"estimate"`
		Active   bool   `json:"active"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.Update(
		outletID,
		uint(id),
		req.Name,
		req.Price,
		req.Estimate,
		req.Active,
	); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "service updated"})
}

// ========================
// DELETE SERVICE
// ========================
func (h *Handler) Delet(c *gin.Context) {
	outletID := c.GetUint("outlet_id")
	if outletID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "outlet context required"})
		return
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid service id"})
		return
	}

	if err := h.service.Delete(outletID, uint(id)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "service deleted"})
}