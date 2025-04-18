package service

import (
	"errors"
	"go-trades/entity"
	"go-trades/repository"
	"go-trades/utils"

	errorMessages "go-trades/utils/error-messages"

	"gorm.io/gorm"
)

type productService struct {
	Repository repository.ProductRepository
}

type ProductService interface {
	GetAllProducts(categoryId uint) (*utils.Response, error)
	GetProductById(id uint) (*utils.Response, error)
	CreateProduct(req *entity.CreateProductRequest) (*utils.Response, error)
	UpdateProduct(id uint, req *entity.UpdateProductRequest) (*utils.Response, error)
	DeleteProduct(id uint) error
}

func NewProductService(r repository.ProductRepository) ProductService {
	return &productService{
		Repository: r,
	}
}

func (s *productService) GetAllProducts(categoryId uint) (*utils.Response, error) {

	var products []entity.Product
	var err error

	if categoryId != 0 {
		products, err = s.Repository.FindByCategoryId(categoryId)
		if err != nil {
			return nil, err
		}

	} else {
		products, err = s.Repository.FindAll()
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

func (s *productService) GetProductById(id uint) (*utils.Response, error) {
	product, err := s.Repository.FindById(id)
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

func (s *productService) CreateProduct(req *entity.CreateProductRequest) (*utils.Response, error) {

	existingByName, err := s.Repository.FindByName(req.Name)
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

	if err := s.Repository.CreateProduct(product); err != nil {
		return nil, err
	}

	savedProduct, err := s.Repository.FindById(product.ID)
	if err != nil {
		return nil, errors.New("error loading product data")
	}

	productData := entity.ProductDataResponse{
		ID:          savedProduct.ID,
		CategoryId:  savedProduct.CategoryId,
		Name:        savedProduct.Name,
		Description: savedProduct.Description,
		Price:       savedProduct.Price,
		CreatedAt:   savedProduct.CreatedAt,
		UpdatedAt:   savedProduct.UpdatedAt,
	}

	return &utils.Response{
		Status:  201,
		Message: "Product successfully created",
		Data:    productData,
	}, nil
}

func (s *productService) UpdateProduct(id uint, req *entity.UpdateProductRequest) (*utils.Response, error) {

	product, err := s.Repository.FindById(id)
	if err != nil {
		return nil, err
	}

	if product == nil {
		return nil, errors.New(errorMessages.ErrProductNotFound)
	}

	existingByName, err := s.Repository.FindByName(req.Name)
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

	if err := s.Repository.UpdateProduct(product); err != nil {
		return nil, err
	}

	productData := entity.ProductDataResponse{
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
		Data:    productData,
	}, nil
}

func (s *productService) DeleteProduct(id uint) error {
	if err := s.Repository.DeleteProduct(id); err != nil {
		return err
	}

	return nil
}
