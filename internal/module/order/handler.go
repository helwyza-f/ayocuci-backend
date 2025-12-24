package order

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
// CREATE ORDER
// ========================
func (h *Handler) Create(c *gin.Context) {
	var req struct {
		CustomerID uint        `json:"customer_id"`
		Items      []ItemInput `json:"items"`
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

	if err := h.service.Create(outletID, req.CustomerID, req.Items); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "order created"})
}

// ========================
// LIST ORDERS (PAGINATION + FILTER)
// ========================
func (h *Handler) List(c *gin.Context) {
	outletID := c.GetUint("outlet_id")
	if outletID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "outlet context required"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	status := c.Query("status")

	var fromPtr *time.Time
	var toPtr *time.Time

	if from := c.Query("from"); from != "" {
		t, err := time.Parse("2006-01-02", from)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid from date"})
			return
		}
		fromPtr = &t
	}

	if to := c.Query("to"); to != "" {
		t, err := time.Parse("2006-01-02", to)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid to date"})
			return
		}
		t = t.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
		toPtr = &t
	}

	filter := OrderFilter{
		Status: status,
		From:   fromPtr,
		To:     toPtr,
		Page:   page,
		Limit:  limit,
	}

	orders, total, err := h.service.GetAll(outletID, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))

	c.JSON(http.StatusOK, gin.H{
		"data": orders,
		"meta": gin.H{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": totalPages,
		},
	})
}

// ========================
// UPDATE ORDER STATUS + AUDIT LOG
// ========================
func (h *Handler) UpdateStatus(c *gin.Context) {
	outletID := c.GetUint("outlet_id")
	clientID := c.GetUint("client_id")
	userID := c.GetUint("user_id")

	if outletID == 0 || clientID == 0 || userID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid auth context"})
		return
	}

	idParam := c.Param("id")
	orderID, err := strconv.Atoi(idParam)
	if err != nil || orderID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order id"})
		return
	}

	var req struct {
		Status string `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.UpdateStatus(
		clientID,
		outletID,
		userID,
		uint(orderID),
		req.Status,
	); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "status updated"})
}
