package service

import (
	"archive/zip"
	"errors"
	"go-trades/entity"
	"go-trades/repository"
	"go-trades/utils"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	maxFileSize     = 2 * 1024 * 1024 // 2MB
	uploadDirectory = "uploads"
)

type productImageService struct {
	DB                     *gorm.DB
	ProductImageRepository repository.ProductImageRepository
	ProductRepository      repository.ProductRepository
}

type ProductImageService interface {
	UploadProductImage(ctx *gin.Context, productId uint, image *multipart.FileHeader) error
	DownloadProductImages(ctx *gin.Context, productId uint) (string, error)
}

func NewProductImageService(db *gorm.DB, pir repository.ProductImageRepository, pr repository.ProductRepository) ProductImageService {
	return &productImageService{
		ProductImageRepository: pir,
		ProductRepository:      pr,
		DB:                     db,
	}
}

func (s *productImageService) UploadProductImage(ctx *gin.Context, productId uint, image *multipart.FileHeader) error {

	if image.Size > maxFileSize {
		return errors.New("maximum file size is 2MB")
	}

	if !utils.IsValidExtension(image) {
		return errors.New("please upload file with extension PNG, JPG, JPEG")
	}

	product, err := s.ProductRepository.FindById(ctx, productId)
	if err != nil {
		return err
	}
	if product == nil {
		return errors.New("product not found")
	}

	tx := s.DB.Begin()
	utils.WithTx(ctx, tx)

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		} else if tx != nil {
			tx.Rollback()
		}
	}()

	imagePath := filepath.Join(uploadDirectory, image.Filename)
	if err := ctx.SaveUploadedFile(image, imagePath); err != nil {
		return errors.New("error uploading file")
	}

	productImage := &entity.ProductImage{
		ProductId:  productId,
		ImageUrl:   imagePath,
		FileName:   image.Filename,
		UploadedAt: time.Now(),
	}

	if err := s.ProductImageRepository.CreateProductImage(ctx, productImage); err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	tx = nil

	return nil
}

func (s *productImageService) DownloadProductImages(ctx *gin.Context, productId uint) (string, error) {
	product, err := s.ProductRepository.FindById(ctx, productId)
	if err != nil {
		return "", err
	}
	if product == nil {
		return "", errors.New("product not found")
	}

	images, err := s.ProductImageRepository.FindAllByProductId(ctx, productId)
	if err != nil {
		return "", err
	}
	if len(images) == 0 {
		return "", errors.New("no images found for this product")
	}

	uuidStr := uuid.New().String()
	zipFileName := filepath.Join(uploadDirectory, uuidStr+".zip")
	zipFile, err := os.Create(zipFileName)
	if err != nil {
		return "", errors.New("error creating ZIP file")
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	for _, img := range images {
		file, err := os.Open(img.ImageUrl)
		if err != nil {
			return "", errors.New("error opening image file: " + img.FileName)
		}

		zipEntry, err := zipWriter.Create(img.FileName)
		if err != nil {
			file.Close()
			return "", errors.New("error adding file to ZIP: " + img.FileName)
		}

		if _, err := io.Copy(zipEntry, file); err != nil {
			file.Close()
			return "", errors.New("error writing file to ZIP: " + img.FileName)
		}
		file.Close()
	}

	go func() {
		time.Sleep(30 * time.Second)
		os.Remove(zipFileName)
	}()

	return zipFileName, nil
}
