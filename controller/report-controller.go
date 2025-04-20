package controller

import (
	"go-trades/service"
	"time"

	"github.com/gin-gonic/gin"
)

type ReportController struct {
	Service service.ReportService
}

func NewReportController(s service.ReportService) *ReportController {
	return &ReportController{
		Service: s,
	}
}

func (c *ReportController) GetReport(ctx *gin.Context) {
	startStr := ctx.Query("startDate")
	endStr := ctx.Query("endDate")

	if startStr == "" || endStr == "" {
		ctx.JSON(400, gin.H{"error": "startDate and endDate are required"})
		return
	}

	start, err := time.Parse("2006-01-02", startStr)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid startDate format"})
		return
	}
	end, err := time.Parse("2006-01-02", endStr)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid endDate format"})
		return
	}

	resp, err := c.Service.GetReport(ctx, start, end)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, resp)
}
