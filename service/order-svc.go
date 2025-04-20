package service

import (
	"errors"
	"go-trades/entity"
	"go-trades/repository"
	"go-trades/utils"
	errorMessages "go-trades/utils/error-messages"
	status "go-trades/utils/status"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type orderService struct {
	db                  *gorm.DB
	OrderRepository     repository.OrderRepository
	ProductRepository   repository.ProductRepository
	InventoryRepository repository.InventoryRepository
}

type OrderService interface {
	GetAllOrders(ctx *gin.Context, status uint) (*utils.Response, error)
	GetOrderById(ctx *gin.Context, id uint) (*utils.Response, error)
	CreateOrder(ctx *gin.Context, req *entity.CreateOrderRequest) (*utils.Response, error)
	ProcessOrder(ctx *gin.Context, id uint) (*utils.Response, error)
	ConfirmOrder(ctx *gin.Context, id uint) (*utils.Response, error)
	CancelOrder(ctx *gin.Context, id uint) error
}

func NewOrderService(db *gorm.DB, or repository.OrderRepository, pr repository.ProductRepository, ir repository.InventoryRepository) OrderService {
	return &orderService{
		db:                  db,
		OrderRepository:     or,
		ProductRepository:   pr,
		InventoryRepository: ir,
	}
}

func (s *orderService) GetAllOrders(ctx *gin.Context, status uint) (*utils.Response, error) {
	var orders []entity.Order
	var err error

	if status != 0 {
		orders, err = s.OrderRepository.FindByStatus(ctx, status)
		if err != nil {
			return nil, err
		}

	} else {
		orders, err = s.OrderRepository.FindAll(ctx)
		if err != nil {
			return nil, err
		}
	}

	data := make([]entity.OrderDataResponse, len(orders))
	for i, order := range orders {
		data[i] = entity.OrderDataResponse{
			ID:                  order.ID,
			UserId:              order.UserId,
			Date:                order.Date,
			ShippingAddress:     order.ShippingAddress,
			Total:               order.Total,
			Status:              order.Status,
			OrderDetailResponse: make([]entity.OrderDetailResponse, len(order.OrderDetails)),
		}

		for j, orderDetail := range order.OrderDetails {
			data[i].OrderDetailResponse[j] = entity.OrderDetailResponse{
				ProductId: orderDetail.ProductId,
				Qty:       orderDetail.Qty,
				Subtotal:  orderDetail.Subtotal,
			}
		}
	}

	return &utils.Response{
		Status:  200,
		Message: "Success",
		Data:    data,
	}, nil
}

func (s *orderService) GetOrderById(ctx *gin.Context, id uint) (*utils.Response, error) {
	order, err := s.OrderRepository.FindById(ctx, id)
	if err != nil {
		return nil, err
	}

	if order == nil {
		return nil, errors.New(errorMessages.ErrCategoryNotFound)
	}

	data := entity.OrderDataResponse{
		ID:                  order.ID,
		UserId:              order.UserId,
		Date:                order.Date,
		ShippingAddress:     order.ShippingAddress,
		Total:               order.Total,
		Status:              order.Status,
		OrderDetailResponse: make([]entity.OrderDetailResponse, len(order.OrderDetails)),
	}

	for i, od := range order.OrderDetails {
		data.OrderDetailResponse[i] = entity.OrderDetailResponse{
			ProductId: od.ProductId,
			Qty:       od.Qty,
			Subtotal:  od.Subtotal,
		}
	}
	return &utils.Response{
		Status:  200,
		Message: "Success",
		Data:    data,
	}, nil
}

func (s *orderService) CreateOrder(ctx *gin.Context, req *entity.CreateOrderRequest) (*utils.Response, error) {

	productIds := make(map[uint]bool)
	for _, d := range req.OrderDetails {
		if productIds[d.ProductId] {
			return nil, errors.New(errorMessages.ErrOrderDuplicateProduct)
		}
		productIds[d.ProductId] = true
	}

	var total uint
	var orderDetails []entity.OrderDetail

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

	for _, detail := range req.OrderDetails {

		product, err := s.ProductRepository.FindById(ctx, detail.ProductId)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		if product == nil {
			tx.Rollback()
			return nil, errors.New(errorMessages.ErrProductNotFound)
		}

		total += detail.Qty * product.Price
		orderDetails = append(orderDetails, entity.OrderDetail{
			ProductId: detail.ProductId,
			Qty:       detail.Qty,
			Subtotal:  product.Price * detail.Qty,
		})

		inventory, err := s.InventoryRepository.FindFirstByProductId(ctx, detail.ProductId)
		if err != nil {
			return nil, err
		}

		if inventory.Stock < detail.Qty {
			return nil, errors.New(errorMessages.ErrInventoryInsufficientStock)
		}

		if err := s.InventoryRepository.UpdateInventoryForOrder(ctx, inventory, detail.Qty); err != nil {
			tx.Rollback()
			return nil, errors.New(errorMessages.ErrInventoryStockUpdate)
		}
	}
	order := entity.Order{
		UserId:          1,
		Date:            time.Now(),
		ShippingAddress: req.ShippingAddress,
		Total:           total,
		Status:          status.PENDING,
		OrderDetails:    orderDetails,
	}
	if err := s.OrderRepository.CreateOrder(ctx, &order); err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()
	tx = nil

	data := entity.OrderDataResponse{
		ID:                  order.ID,
		UserId:              order.UserId,
		Date:                order.Date,
		ShippingAddress:     order.ShippingAddress,
		Total:               order.Total,
		OrderDetailResponse: make([]entity.OrderDetailResponse, len(order.OrderDetails)),
	}

	for i, od := range order.OrderDetails {
		data.OrderDetailResponse[i] = entity.OrderDetailResponse{
			ProductId: od.ProductId,
			Qty:       od.Qty,
			Subtotal:  od.Subtotal,
		}
	}

	return &utils.Response{
		Status:  201,
		Message: "Success",
		Data:    data,
	}, nil
}

func (s *orderService) ProcessOrder(ctx *gin.Context, id uint) (*utils.Response, error) {
	product, err := s.OrderRepository.FindById(ctx, id)
	if err != nil {
		return nil, err
	}

	if product == nil {
		return nil, errors.New(errorMessages.ErrProductNotFound)
	}

	if product.Status != status.PAID {
		return nil, errors.New(errorMessages.ErrInvalidOrderStatus)
	}

	product.Status = status.PROCESSING
	if err := s.OrderRepository.UpdateOrder(ctx, product); err != nil {
		return nil, err
	}
	data := entity.OrderDataResponse{
		ID:                  product.ID,
		UserId:              product.UserId,
		Date:                product.Date,
		ShippingAddress:     product.ShippingAddress,
		Total:               product.Total,
		Status:              product.Status,
		OrderDetailResponse: make([]entity.OrderDetailResponse, len(product.OrderDetails)),
	}

	return &utils.Response{
		Status:  200,
		Message: "Success",
		Data:    data,
	}, nil
}

func (s *orderService) ConfirmOrder(ctx *gin.Context, id uint) (*utils.Response, error) {
	product, err := s.OrderRepository.FindById(ctx, id)
	if err != nil {
		return nil, err
	}

	if product == nil {
		return nil, errors.New(errorMessages.ErrProductNotFound)
	}

	if product.Status != status.PROCESSING {
		return nil, errors.New(errorMessages.ErrInvalidOrderStatus)
	}

	product.Status = status.DONE
	if err := s.OrderRepository.UpdateOrder(ctx, product); err != nil {
		return nil, err
	}
	data := entity.OrderDataResponse{
		ID:                  product.ID,
		UserId:              product.UserId,
		Date:                product.Date,
		ShippingAddress:     product.ShippingAddress,
		Total:               product.Total,
		Status:              product.Status,
		OrderDetailResponse: make([]entity.OrderDetailResponse, len(product.OrderDetails)),
	}

	return &utils.Response{
		Status:  200,
		Message: "Success",
		Data:    data,
	}, nil
}

func (s *orderService) CancelOrder(ctx *gin.Context, id uint) error {

	order, err := s.OrderRepository.FindById(ctx, id)
	if err != nil {
		return err
	}

	if order == nil {
		return errors.New(errorMessages.ErrProductNotFound)
	}

	if order.Status != status.PENDING {
		return errors.New(errorMessages.ErrOrderUncancelable)
	}

	if err := s.OrderRepository.DeleteOrder(ctx, id); err != nil {
		return err
	}

	return nil
}
