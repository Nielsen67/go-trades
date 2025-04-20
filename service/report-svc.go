package service

import (
	"go-trades/entity"
	"go-trades/repository"
	"go-trades/utils"
	"time"

	"github.com/gin-gonic/gin"
)

type reportService struct {
	Repository repository.ReportRepository
}

type ReportService interface {
	GetReport(ctx *gin.Context, start time.Time, end time.Time) (*utils.Response, error)
}

func NewReportService(r repository.ReportRepository) ReportService {
	return &reportService{
		Repository: r,
	}
}

func (r *reportService) GetReport(ctx *gin.Context, start time.Time, end time.Time) (*utils.Response, error) {
	BestSelling, err := r.Repository.FindBestSelling(ctx, start, end)
	if err != nil {
		return nil, err
	}

	LowInventory, err := r.Repository.FindLowStock(ctx)
	if err != nil {
		return nil, err
	}

	OrderSummary, err := r.Repository.GenerateOrderSummary(ctx, start, end)
	if err != nil {
		return nil, err
	}

	data := entity.Report{
		BestSellingProducts: BestSelling,
		LowestInventories:   LowInventory,
		OrderSummary:        OrderSummary,
	}

	return &utils.Response{
		Status:  200,
		Message: "Success",
		Data:    data,
	}, nil
}
