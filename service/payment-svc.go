package service

import (
	"errors"
	"go-trades/entity"
	"go-trades/repository"
	"go-trades/utils"
	errorMessages "go-trades/utils/error-messages"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type paymentService struct {
	db                *gorm.DB
	PaymentRepository repository.PaymentRepository
	OrderRepository   repository.OrderRepository
}

type PaymentService interface {
	GetAllPayments(ctx *gin.Context, page, size int) (*utils.Response, int64, int64, error)
	CreatePayment(ctx *gin.Context, userId uint, req *entity.PaymentRequest) (*utils.Response, error)
	GetUserPayments(ctx *gin.Context, userId uint, page, size int) (*utils.Response, int64, int64, error)
}

func NewPaymentService(db *gorm.DB, pr repository.PaymentRepository, or repository.OrderRepository) PaymentService {
	return &paymentService{
		db:                db,
		PaymentRepository: pr,
		OrderRepository:   or,
	}
}

func (s *paymentService) GetAllPayments(ctx *gin.Context, page, size int) (*utils.Response, int64, int64, error) {
	payments, totalSize, err := s.PaymentRepository.FindAll(ctx, page, size)
	if err != nil {
		return nil, 0, 0, err
	}
	data := make([]entity.PaymentDataResponse, len(payments))
	for i, payment := range payments {
		data[i] = entity.PaymentDataResponse{
			ID:        payment.ID,
			OrderId:   payment.OrderId,
			Method:    payment.Method,
			Amount:    payment.Amount,
			Status:    payment.Status,
			CreatedAt: payment.CreatedAt,
		}
	}

	totalPage := utils.GetTotalPage(totalSize, size)

	return &utils.Response{
		Status:  200,
		Message: "Success",
		Data:    data,
	}, totalSize, totalPage, nil
}

func (s *paymentService) GetUserPayments(ctx *gin.Context, userId uint, page, size int) (*utils.Response, int64, int64, error) {
	payments, totalSize, err := s.PaymentRepository.FindAllByUserId(ctx, userId, page, size)
	if err != nil {
		return nil, 0, 0, err
	}

	data := make([]entity.PaymentDataResponse, len(payments))
	for i, payment := range payments {
		data[i] = entity.PaymentDataResponse{
			ID:        payment.ID,
			OrderId:   payment.OrderId,
			Method:    payment.Method,
			Amount:    payment.Amount,
			Status:    payment.Status,
			CreatedAt: payment.CreatedAt,
		}
	}

	totalPage := utils.GetTotalPage(totalSize, size)

	return &utils.Response{
		Status:  200,
		Message: "Success",
		Data:    data,
	}, totalSize, totalPage, nil
}

func (s *paymentService) CreatePayment(ctx *gin.Context, userId uint, req *entity.PaymentRequest) (*utils.Response, error) {

	tx := s.db.Begin()
	utils.WithTx(ctx, tx)
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		} else if tx != nil {
			tx.Rollback()
		}
	}()

	order, err := s.OrderRepository.FindByUserIdWithId(ctx, userId, req.OrderId)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	if order == nil {
		tx.Rollback()
		return nil, errors.New(errorMessages.ErrOrderNotFound)
	}

	if order.Status != 1 {
		tx.Rollback()
		return nil, errors.New(errorMessages.ErrAlreadyPaid)
	}

	if order.Total != req.Amount {
		tx.Rollback()
		return nil, errors.New(errorMessages.ErrPaymentAmountInsufficient)
	}

	order.Status = 2
	if err := s.OrderRepository.UpdateOrder(ctx, order); err != nil {
		tx.Rollback()
		return nil, err
	}

	payment := &entity.Payment{
		OrderId: req.OrderId,
		Method:  req.Method,
		Amount:  req.Amount,
		Status:  2,
	}

	if err := s.PaymentRepository.CreatePayment(ctx, payment); err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()
	tx = nil

	data := entity.PaymentDataResponse{
		ID:        payment.ID,
		OrderId:   payment.OrderId,
		Method:    payment.Method,
		Amount:    payment.Amount,
		Status:    payment.Status,
		CreatedAt: payment.CreatedAt,
	}

	return &utils.Response{
		Status:  201,
		Message: "Inventory successfully created",
		Data:    data,
	}, nil
}
