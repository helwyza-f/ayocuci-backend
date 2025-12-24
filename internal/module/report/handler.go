package report

import (
	"log"
	"net/http"
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
// DAILY REPORT
// ========================
// GET /api/reports/daily?date=YYYY-MM-DD
func (h *Handler) Daily(c *gin.Context) {
	outletID := c.GetUint("outlet_id")
	if outletID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "outlet context required",
		})
		return
	}

	dateParam := c.Query("date") // optional: YYYY-MM-DD
	var date time.Time
	var err error

	loc, _ := time.LoadLocation("Asia/Jakarta")

	if dateParam == "" {
		date = time.Now().In(loc)
	} else {
		date, err = time.ParseInLocation("2006-01-02", dateParam, loc)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid date format, use YYYY-MM-DD",
			})
			return
		}
	}


	log.Print(date)

	report, err := h.service.Daily(outletID, date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": report,
	})
}

// ========================
// RANGE REPORT
// ========================
// GET /api/reports/range?start=YYYY-MM-DD&end=YYYY-MM-DD
func (h *Handler) Range(c *gin.Context) {
	outletID := c.GetUint("outlet_id")
	if outletID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "outlet context required",
		})
		return
	}

	startParam := c.Query("start")
	endParam := c.Query("end")

	if startParam == "" || endParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "start and end date are required",
		})
		return
	}

	start, err := time.Parse("2006-01-02", startParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid start date format, use YYYY-MM-DD",
		})
		return
	}

	end, err := time.Parse("2006-01-02", endParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid end date format, use YYYY-MM-DD",
		})
		return
	}

	// include full end day
	end = end.Add(24 * time.Hour)

	report, err := h.service.Range(outletID, start, end)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": report,
	})
}
