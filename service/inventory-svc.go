package service

import (
	"errors"
	"go-trades/entity"
	"go-trades/repository"
	"go-trades/utils"
	errorMessages "go-trades/utils/error-messages"

	"github.com/gin-gonic/gin"
)

type inventoryService struct {
	InventoryRepository repository.InventoryRepository
	ProductRepository   repository.ProductRepository
}

type InventoryService interface {
	GetAllInventories(ctx *gin.Context, page, size int) (*utils.Response, int64, int64, error)
	GetInventoryById(ctx *gin.Context, id uint) (*utils.Response, error)
	CreateInventory(ctx *gin.Context, req *entity.CreateInventoryRequest) (*utils.Response, error)
	UpdateInventory(ctx *gin.Context, id uint, req *entity.UpdateInventoryRequest) (*utils.Response, error)
	DeleteInventory(ctx *gin.Context, id uint) error
}

func NewInventoryService(ir repository.InventoryRepository, pr repository.ProductRepository) InventoryService {
	return &inventoryService{
		InventoryRepository: ir,
		ProductRepository:   pr,
	}
}

func (s *inventoryService) GetAllInventories(ctx *gin.Context, page, size int) (*utils.Response, int64, int64, error) {
	inventories, totalSize, err := s.InventoryRepository.FindAll(ctx, page, size)
	if err != nil {
		return nil, 0, 0, err
	}
	data := make([]entity.InventoryDataResponse, len(inventories))
	for i, inventory := range inventories {
		data[i] = entity.InventoryDataResponse{
			ID:        inventory.ID,
			ProductId: inventory.ProductId,
			Stock:     inventory.Stock,
			Location:  inventory.Location,
		}
	}

	totalPage := utils.GetTotalPage(totalSize, size)

	return &utils.Response{
		Status:  200,
		Message: "Success",
		Data:    data,
	}, totalSize, totalPage, nil
}

func (s *inventoryService) GetInventoryById(ctx *gin.Context, id uint) (*utils.Response, error) {
	inventory, err := s.InventoryRepository.FindById(ctx, id)
	if err != nil {
		return nil, err
	}

	if inventory == nil {
		return nil, errors.New(errorMessages.ErrInventoryNotFound)
	}

	data := entity.InventoryDataResponse{
		ID:        inventory.ID,
		ProductId: inventory.ProductId,
		Stock:     inventory.Stock,
		Location:  inventory.Location,
	}
	return &utils.Response{
		Status:  200,
		Message: "Success",
		Data:    data,
	}, nil
}

func (s *inventoryService) CreateInventory(ctx *gin.Context, req *entity.CreateInventoryRequest) (*utils.Response, error) {

	product, err := s.ProductRepository.FindById(ctx, req.ProductId)
	if err != nil {
		return nil, err
	}
	if product == nil {
		return nil, errors.New(errorMessages.ErrProductNotFound)
	}

	inventory := &entity.Inventory{
		ProductId: req.ProductId,
		Stock:     req.Stock,
		Location:  req.Location,
	}

	if err := s.InventoryRepository.CreateInventory(ctx, inventory); err != nil {
		return nil, err
	}

	savedInventory, err := s.InventoryRepository.FindById(ctx, inventory.ID)
	if err != nil {
		return nil, errors.New("error loading inventory data")
	}

	data := entity.InventoryDataResponse{
		ID:        savedInventory.ID,
		ProductId: savedInventory.ProductId,
		Stock:     savedInventory.Stock,
		Location:  savedInventory.Location,
	}

	return &utils.Response{
		Status:  201,
		Message: "Inventory successfully created",
		Data:    data,
	}, nil
}

func (s *inventoryService) UpdateInventory(ctx *gin.Context, id uint, req *entity.UpdateInventoryRequest) (*utils.Response, error) {

	inventory, err := s.InventoryRepository.FindById(ctx, id)
	if err != nil {
		return nil, err
	}

	if inventory == nil {
		return nil, errors.New(errorMessages.ErrInventoryNotFound)
	}

	if req.Stock <= 0 {
		return nil, errors.New(errorMessages.ErrInventoryInvalidStock)
	}

	inventory.Stock = req.Stock

	if err := s.InventoryRepository.UpdateInventory(ctx, inventory); err != nil {
		return nil, err
	}

	data := entity.InventoryDataResponse{
		ID:        inventory.ID,
		ProductId: inventory.ProductId,
		Stock:     inventory.Stock,
		Location:  inventory.Location,
	}

	return &utils.Response{
		Status:  200,
		Message: "Inventory stock updated",
		Data:    data,
	}, nil
}

func (s *inventoryService) DeleteInventory(ctx *gin.Context, id uint) error {
	if err := s.InventoryRepository.DeleteInventory(ctx, id); err != nil {
		return err
	}

	return nil
}
