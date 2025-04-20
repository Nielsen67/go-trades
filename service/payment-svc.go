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
	GetAllPayments(ctx *gin.Context) (*utils.Response, error)
	CreatePayment(ctx *gin.Context, req *entity.PaymentRequest) (*utils.Response, error)
}

func NewPaymentService(db *gorm.DB, pr repository.PaymentRepository, or repository.OrderRepository) PaymentService {
	return &paymentService{
		db:                db,
		PaymentRepository: pr,
		OrderRepository:   or,
	}
}

func (s *paymentService) GetAllPayments(ctx *gin.Context) (*utils.Response, error) {
	payments, err := s.PaymentRepository.FindAll(ctx)
	if err != nil {
		return nil, err
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

	return &utils.Response{
		Status:  200,
		Message: "Success",
		Data:    data,
	}, nil
}

func (s *paymentService) CreatePayment(ctx *gin.Context, req *entity.PaymentRequest) (*utils.Response, error) {

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

	order, err := s.OrderRepository.FindById(ctx, req.OrderId)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	if order == nil {
		tx.Rollback()
		return nil, errors.New(errorMessages.ErrProductNotFound)
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
