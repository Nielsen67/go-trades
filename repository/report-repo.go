package repository

import (
	"go-trades/entity"
	"go-trades/utils"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type reportRepository struct {
	DB *gorm.DB
}

type ReportRepository interface {
	FindBestSelling(ctx *gin.Context, start time.Time, end time.Time) ([]entity.BestSellingProduct, error)
	FindLowStock(ctx *gin.Context) ([]entity.LowInventoryItem, error)
	GenerateOrderSummary(ctx *gin.Context, start time.Time, end time.Time) (*entity.OrderSummary, error)
}

func NewReportRepository(db *gorm.DB) ReportRepository {
	return &reportRepository{
		DB: db,
	}
}

func (r *reportRepository) FindBestSelling(ctx *gin.Context, start time.Time, end time.Time) ([]entity.BestSellingProduct, error) {
	db := utils.GetTx(ctx, r.DB)
	var results []entity.BestSellingProduct
	err := db.Raw(`
        SELECT 
            p.id AS product_id,
            p.name AS product_name,
            SUM(od.qty) AS qty_sold
        FROM 
            order_details od
            JOIN 
                orders o ON o.id = od.order_id
            JOIN 
                products p ON p.id = od.product_id
        WHERE 
            o.status IN (2,3,4) AND o.date BETWEEN ? AND ?
        GROUP BY 
            p.id, p.name
        ORDER BY 
            qty_sold DESC
        LIMIT 5
    `, start, end).Scan(&results).Error

	if err != nil {
		log.Printf("Error in FindBestSelling: %v", err)
	}
	log.Printf("FindBestSelling results: %+v", results)
	return results, err
}
func (r *reportRepository) FindLowStock(ctx *gin.Context) ([]entity.LowInventoryItem, error) {
	db := utils.GetTx(ctx, r.DB)
	var results []entity.LowInventoryItem
	err := db.Raw(`
	SELECT 
		i.product_id ,
		p.name AS product_name,
		i.stock
	FROM 
		inventories i
	JOIN 
		products p ON p.id = i.product_id
	ORDER BY 
		i.stock ASC
	LIMIT 5
	`).Scan(&results).Error

	return results, err
}

func (r *reportRepository) GenerateOrderSummary(ctx *gin.Context, start time.Time, end time.Time) (*entity.OrderSummary, error) {
	db := utils.GetTx(ctx, r.DB)
	var result entity.OrderSummary
	err := db.Raw(`
        SELECT 
            COUNT(o.id) AS total_orders,
            SUM(o.total) AS total_amount,
            COUNT(DISTINCT o.user_id) AS total_customers
        FROM 
            orders o
        WHERE 
            o.status = 2 AND o.date BETWEEN ? AND ?
    `, start, end).Scan(&result).Error

	return &result, err
}
