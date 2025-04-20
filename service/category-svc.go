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

type categoryService struct {
	Repository repository.CategoryRepository
}

type CategoryService interface {
	GetAllCategories(ctx *gin.Context) (*utils.Response, error)
	GetCategoryById(ctx *gin.Context, id uint) (*utils.Response, error)
	CreateCategory(ctx *gin.Context, req *entity.CategoryRequest) (*utils.Response, error)
	UpdateCategory(ctx *gin.Context, id uint, req *entity.CategoryRequest) (*utils.Response, error)
	DeleteCategory(ctx *gin.Context, id uint) error
}

func NewCategoryService(r repository.CategoryRepository) CategoryService {
	return &categoryService{
		Repository: r,
	}
}

func (s *categoryService) GetAllCategories(ctx *gin.Context) (*utils.Response, error) {
	categories, err := s.Repository.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	data := make([]entity.CategoryDataResponse, len(categories))
	for i, category := range categories {
		data[i] = entity.CategoryDataResponse{
			ID:   category.ID,
			Code: category.Code,
			Name: category.Name,
		}
	}

	return &utils.Response{
		Status:  200,
		Message: "Success",
		Data:    data,
	}, nil
}

func (s *categoryService) GetCategoryById(ctx *gin.Context, id uint) (*utils.Response, error) {
	category, err := s.Repository.FindById(ctx, id)
	if err != nil {
		return nil, err
	}

	if category == nil {
		return nil, errors.New(errorMessages.ErrCategoryNotFound)
	}

	data := entity.CategoryDataResponse{
		ID:   category.ID,
		Code: category.Code,
		Name: category.Name,
	}
	return &utils.Response{
		Status:  200,
		Message: "Success",
		Data:    data,
	}, nil
}

func (s *categoryService) CreateCategory(ctx *gin.Context, req *entity.CategoryRequest) (*utils.Response, error) {

	existingByName, err := s.Repository.FindByName(ctx, req.Name)
	if err == nil && existingByName != nil {
		return nil, errors.New(errorMessages.ErrCategoryNameExists)
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) && err != nil {
		return nil, err
	}

	existingByCode, err := s.Repository.FindByCode(ctx, req.Code)
	if err == nil && existingByCode != nil {
		return nil, errors.New(errorMessages.ErrCategoryCodeExists)
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) && err != nil {
		return nil, err
	}

	category := &entity.Category{
		Code: req.Code,
		Name: req.Name,
	}

	if err := s.Repository.CreateCategory(ctx, category); err != nil {
		return nil, err
	}

	savedCategory, err := s.Repository.FindById(ctx, category.ID)
	if err != nil {
		return nil, errors.New("error loading category data")
	}

	data := entity.CategoryDataResponse{
		ID:   savedCategory.ID,
		Code: savedCategory.Code,
		Name: savedCategory.Name,
	}

	return &utils.Response{
		Status:  201,
		Message: "Category successfully created",
		Data:    data,
	}, nil
}

func (s *categoryService) UpdateCategory(ctx *gin.Context, id uint, req *entity.CategoryRequest) (*utils.Response, error) {

	category, err := s.Repository.FindById(ctx, id)
	if err != nil {
		return nil, err
	}

	if category == nil {
		return nil, errors.New(errorMessages.ErrCategoryNotFound)
	}

	existingByName, err := s.Repository.FindByName(ctx, req.Name)
	if err != nil {
		return nil, err
	}
	if existingByName != nil && existingByName.ID != id {
		return nil, errors.New(errorMessages.ErrCategoryNameExists)
	}

	existingByCode, err := s.Repository.FindByCode(ctx, req.Code)
	if err != nil {
		return nil, err
	}
	if existingByCode != nil && existingByCode.ID != id {
		return nil, errors.New(errorMessages.ErrCategoryCodeExists)
	}

	category.Code = req.Code
	category.Name = req.Name

	if err := s.Repository.UpdateCategory(ctx, category); err != nil {
		return nil, err
	}

	data := entity.CategoryDataResponse{
		ID:   category.ID,
		Code: category.Code,
		Name: category.Name,
	}

	return &utils.Response{
		Status:  200,
		Message: "Category successfully updated",
		Data:    data,
	}, nil
}

func (s *categoryService) DeleteCategory(ctx *gin.Context, id uint) error {
	if err := s.Repository.DeleteCategory(ctx, id); err != nil {
		return err
	}

	return nil
}
