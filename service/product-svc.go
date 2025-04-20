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
	Repository repository.ProductRepository
}

type ProductService interface {
	GetAllProducts(ctx *gin.Context, categoryId uint) (*utils.Response, error)
	GetProductById(ctx *gin.Context, id uint) (*utils.Response, error)
	CreateProduct(ctx *gin.Context, req *entity.CreateProductRequest) (*utils.Response, error)
	UpdateProduct(ctx *gin.Context, id uint, req *entity.UpdateProductRequest) (*utils.Response, error)
	DeleteProduct(ctx *gin.Context, id uint) error
}

func NewProductService(r repository.ProductRepository) ProductService {
	return &productService{
		Repository: r,
	}
}

func (s *productService) GetAllProducts(ctx *gin.Context, categoryId uint) (*utils.Response, error) {

	var products []entity.Product
	var err error

	if categoryId != 0 {
		products, err = s.Repository.FindByCategoryId(ctx, categoryId)
		if err != nil {
			return nil, err
		}

	} else {
		products, err = s.Repository.FindAll(ctx)
		if err != nil {
			return nil, err
		}
	}

	data := make([]entity.ProductDataResponse, len(products))
	for i, product := range products {
		data[i] = entity.ProductDataResponse{
			ID:          product.ID,
			CategoryId:  product.CategoryId,
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
			CreatedAt:   product.CreatedAt,
			UpdatedAt:   product.UpdatedAt,
		}
	}

	return &utils.Response{
		Status:  200,
		Message: "Success",
		Data:    data,
	}, nil
}

func (s *productService) GetProductById(ctx *gin.Context, id uint) (*utils.Response, error) {
	product, err := s.Repository.FindById(ctx, id)
	if err != nil {
		return nil, err
	}

	if product == nil {
		return nil, errors.New(errorMessages.ErrProductNotFound)
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
		Message: "Success",
		Data:    data,
	}, nil
}

func (s *productService) CreateProduct(ctx *gin.Context, req *entity.CreateProductRequest) (*utils.Response, error) {

	existingByName, err := s.Repository.FindByName(ctx, req.Name)
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

	if err := s.Repository.CreateProduct(ctx, product); err != nil {
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

	product, err := s.Repository.FindById(ctx, id)
	if err != nil {
		return nil, err
	}

	if product == nil {
		return nil, errors.New(errorMessages.ErrProductNotFound)
	}

	existingByName, err := s.Repository.FindByName(ctx, req.Name)
	if err != nil {
		return nil, err
	}
	if existingByName != nil && existingByName.ID != id {
		return nil, errors.New(errorMessages.ErrCategoryNameExists)
	}

	product.CategoryId = req.CategoryId
	product.Name = req.Name
	product.Description = req.Description
	product.Price = req.Price

	if err := s.Repository.UpdateProduct(ctx, product); err != nil {
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
	if err := s.Repository.DeleteProduct(ctx, id); err != nil {
		return err
	}

	return nil
}
