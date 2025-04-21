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

type productService struct {
	ProductRepository  repository.ProductRepository
	CategoryRepository repository.CategoryRepository
}

type ProductService interface {
	GetAllProducts(ctx *gin.Context, page, size int, categoryId uint) (*utils.Response, int64, int64, error)
	GetProductById(ctx *gin.Context, id uint) (*utils.Response, error)
	CreateProduct(ctx *gin.Context, req *entity.CreateProductRequest) (*utils.Response, error)
	UpdateProduct(ctx *gin.Context, id uint, req *entity.UpdateProductRequest) (*utils.Response, error)
	DeleteProduct(ctx *gin.Context, id uint) error
}

func NewProductService(pr repository.ProductRepository, cr repository.CategoryRepository) ProductService {
	return &productService{
		ProductRepository:  pr,
		CategoryRepository: cr,
	}
}

func (s *productService) GetAllProducts(ctx *gin.Context, page, size int, categoryId uint) (*utils.Response, int64, int64, error) {

	var data []entity.ProductDataResponse
	var totalSize int64
	var err error

	if categoryId != 0 {
		data, totalSize, err = s.ProductRepository.FindByCategoryIdWithStock(ctx, page, size, categoryId)
		if err != nil {
			return nil, 0, 0, err
		}

	} else {
		data, totalSize, err = s.ProductRepository.FindAllWithStock(ctx, page, size)
		if err != nil {
			return nil, 0, 0, err
		}
	}

	totalPage := utils.GetTotalPage(totalSize, size)

	return &utils.Response{
		Status:  200,
		Message: "Success",
		Data:    data,
	}, totalSize, totalPage, nil
}

func (s *productService) GetProductById(ctx *gin.Context, id uint) (*utils.Response, error) {
	data, err := s.ProductRepository.FindByIdWithStock(ctx, id)
	if err != nil {
		return nil, err
	}

	if data == nil {
		return nil, errors.New(errorMessages.ErrProductNotFound)
	}

	return &utils.Response{
		Status:  200,
		Message: "Success",
		Data:    data,
	}, nil
}

func (s *productService) CreateProduct(ctx *gin.Context, req *entity.CreateProductRequest) (*utils.Response, error) {

	category, err := s.CategoryRepository.FindById(ctx, req.CategoryId)
	if err != nil {
		return nil, err
	}
	if category == nil {
		return nil, errors.New(errorMessages.ErrCategoryNotFound)
	}

	existingByName, err := s.ProductRepository.FindByName(ctx, req.Name)
	if err == nil && existingByName != nil {
		return nil, errors.New(errorMessages.ErrProductNameExists)
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) && err != nil {
		return nil, err
	}

	product := &entity.Product{
		CategoryId:  req.CategoryId,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
	}

	if err := s.ProductRepository.CreateProduct(ctx, product); err != nil {
		return nil, err
	}

	data := entity.ProductDataResponse{
		ID:          product.ID,
		CategoryId:  product.CategoryId,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		CreatedAt:   product.CreatedAt,
		UpdatedAt:   product.UpdatedAt,
	}

	return &utils.Response{
		Status:  201,
		Message: "Product successfully created",
		Data:    data,
	}, nil
}

func (s *productService) UpdateProduct(ctx *gin.Context, id uint, req *entity.UpdateProductRequest) (*utils.Response, error) {

	product, err := s.ProductRepository.FindById(ctx, id)
	if err != nil {
		return nil, err
	}

	if product == nil {
		return nil, errors.New(errorMessages.ErrProductNotFound)
	}

	category, err := s.CategoryRepository.FindById(ctx, req.CategoryId)
	if err != nil {
		return nil, err
	}
	if category == nil {
		return nil, errors.New(errorMessages.ErrCategoryNotFound)
	}

	existingByName, err := s.ProductRepository.FindByName(ctx, req.Name)
	if err != nil {
		return nil, err
	}
	if existingByName != nil && existingByName.ID != id {
		return nil, errors.New(errorMessages.ErrProductNameExists)
	}

	product.CategoryId = req.CategoryId
	product.Name = req.Name
	product.Description = req.Description
	product.Price = req.Price

	if err := s.ProductRepository.UpdateProduct(ctx, product); err != nil {
		return nil, err
	}

	data := entity.ProductDataResponse{
		ID:          product.ID,
		CategoryId:  product.CategoryId,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		CreatedAt:   product.CreatedAt,
		UpdatedAt:   product.UpdatedAt,
	}

	return &utils.Response{
		Status:  200,
		Message: "Product successfully updated",
		Data:    data,
	}, nil
}

func (s *productService) DeleteProduct(ctx *gin.Context, id uint) error {
	if err := s.ProductRepository.DeleteProduct(ctx, id); err != nil {
		return err
	}

	return nil
}
