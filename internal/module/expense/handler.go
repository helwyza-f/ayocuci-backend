package expense

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service}
}

// ========================
// CREATE EXPENSE
// ========================
func (h *Handler) Create(c *gin.Context) {
	var req struct {
		Category string `json:"category"`
		Amount   int64  `json:"amount"`
		Note     string `json:"note"`
		Date     string `json:"date"` // YYYY-MM-DD
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

	var date time.Time
	if req.Date != "" {
		parsed, err := time.Parse("2006-01-02", req.Date)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date format"})
			return
		}
		date = parsed
	}

	if err := h.service.Create(
		outletID,
		req.Category,
		req.Amount,
		req.Note,
		date,
	); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "expense created"})
}

// ========================
// LIST EXPENSES
// ========================
func (h *Handler) List(c *gin.Context) {
	outletID := c.GetUint("outlet_id")
	if outletID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "outlet context required"})
		return
	}

	expenses, err := h.service.GetAll(outletID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": expenses})
}

// ========================
// UPDATE EXPENSE
// ========================
func (h *Handler) Update(c *gin.Context) {
	outletID := c.GetUint("outlet_id")
	if outletID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "outlet context required"})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid expense id"})
		return
	}

	var req struct {
		Category string `json:"category"`
		Amount   int64  `json:"amount"`
		Note     string `json:"note"`
		Date     string `json:"date"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var date time.Time
	if req.Date != "" {
		parsed, err := time.Parse("2006-01-02", req.Date)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date format"})
			return
		}
		date = parsed
	}

	if err := h.service.Update(
		outletID,
		uint(id),
		req.Category,
		req.Amount,
		req.Note,
		date,
	); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "expense updated"})
}

// ========================
// DELETE EXPENSE
// ========================
func (h *Handler) Delete(c *gin.Context) {
	outletID := c.GetUint("outlet_id")
	if outletID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "outlet context required"})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid expense id"})
		return
	}

	if err := h.service.Delete(outletID, uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "expense deleted"})
}
