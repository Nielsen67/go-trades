package service

import (
	"errors"
	"go-trades/entity"
	"go-trades/repository"
	"go-trades/utils"
	errorMessages "go-trades/utils/error-messages"
)

type inventoryService struct {
	Repository repository.InventoryRepository
}

type InventoryService interface {
	GetAllInventories() (*utils.Response, error)
	GetInventoryById(id uint) (*utils.Response, error)
	CreateInventory(req *entity.CreateInventoryRequest) (*utils.Response, error)
	UpdateInventory(id uint, req *entity.UpdateInventoryRequest) (*utils.Response, error)
	DeleteInventory(id uint) error
}

func NewInventoryService(r repository.InventoryRepository) InventoryService {
	return &inventoryService{
		Repository: r,
	}
}

func (s *inventoryService) GetAllInventories() (*utils.Response, error) {
	inventories, err := s.Repository.FindAll()
	if err != nil {
		return nil, err
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

	return &utils.Response{
		Status:  200,
		Message: "Success",
		Data:    data,
	}, nil
}

func (s *inventoryService) GetInventoryById(id uint) (*utils.Response, error) {
	inventory, err := s.Repository.FindById(id)
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

func (s *inventoryService) CreateInventory(req *entity.CreateInventoryRequest) (*utils.Response, error) {

	inventory := &entity.Inventory{
		ProductId: req.ProductId,
		Stock:     req.Stock,
		Location:  req.Location,
	}

	if err := s.Repository.CreateInventory(inventory); err != nil {
		return nil, err
	}

	savedInventory, err := s.Repository.FindById(inventory.ID)
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

func (s *inventoryService) UpdateInventory(id uint, req *entity.UpdateInventoryRequest) (*utils.Response, error) {

	inventory, err := s.Repository.FindById(id)
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

	if err := s.Repository.UpdateInventory(inventory); err != nil {
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

func (s *inventoryService) DeleteInventory(id uint) error {
	if err := s.Repository.DeleteInventory(id); err != nil {
		return err
	}

	return nil
}
