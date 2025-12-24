package outlet

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
func (h *Handler) Delete(c *gin.Context) {
    id, _ := strconv.Atoi(c.Param("id"))
    clientID := c.GetUint("client_id")

    if err := h.service.Delete(uint(id), clientID); err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }
    c.JSON(200, gin.H{"message": "outlet deactivated"})
}

func (h *Handler) Create(c *gin.Context) {
	var req struct {
		Name    string `json:"name"`
		Address string `json:"address"`
		Phone   string `json:"phone"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	clientID := c.GetUint("client_id")
	userID := c.GetUint("user_id")

	outlet, err := h.service.Create(
	clientID,
	userID,
	req.Name,
	req.Address,
	req.Phone,
)
if err != nil {
	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	return
}

c.JSON(http.StatusCreated, gin.H{
	"data": outlet,
})

}

func (h *Handler) List(c *gin.Context) {
	clientID := c.GetUint("client_id")
	role := c.GetString("role")

	outlets, err := h.service.List(clientID, role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": outlets})
}
