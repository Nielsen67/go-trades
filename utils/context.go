package utils

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const ctxTxKey = "db_tx"

func WithTx(ctx *gin.Context, tx *gorm.DB) {
	ctx.Set(ctxTxKey, tx)
}

func GetTx(ctx *gin.Context, db *gorm.DB) *gorm.DB {
	tx, exists := ctx.Get(ctxTxKey)
	if !exists {
		return db
	}
	if dbTx, ok := tx.(*gorm.DB); ok {
		return dbTx
	}
	return db
}
