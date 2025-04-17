package services

import (
	"errors"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"server/internal/app/kanboard/dto"
	"server/internal/constant"
	"server/internal/global"
	"server/internal/repositories"
	md5 "server/pkg/MD5"
	"server/pkg/captcha"

	"github.com/gin-gonic/gin"
)

type PublicService struct {
	captcha      *captcha.Captcha
	resourceRepo *repositories.ResourceRepo
}

var publicService *PublicService

func NewPublicService() *PublicService {
	if publicService == nil {
		redisStore := captcha.NewRedisStore(constant.ADMIN_CAPTCHA, 1*time.Minute, global.Redis)
		publicService = &PublicService{
			captcha:      captcha.NewCaptcha(redisStore),
			resourceRepo: repositories.NewResourceRepo(),
		}
	}

	return publicService
}

func (p *PublicService) GetCaptcha() (*dto.CaptchaResponse, error) {
	id, b64s, _, err := p.captcha.Generate()
	if err != nil {
		return nil, err
	}
	response := &dto.CaptchaResponse{
		Captcha: b64s,
		Id:      id,
	}
	return response, nil
}

func (p *PublicService) VerifyCaptcha(id, answer string) bool {
	return p.captcha.Verify(id, answer)
}

func (p *PublicService) Upload(ctx *gin.Context, MD5 string, file *multipart.FileHeader) (dto.UploadResponse, error) {
	fileMD5, err := md5.GetFileMD5(file)
	if err != nil {
		return dto.UploadResponse{}, err
	}
	if MD5 != *fileMD5 {
		return dto.UploadResponse{}, errors.New("md5校验失败")
	}

	contentType := file.Header.Get("Content-Type")
	filetype := strings.Split(contentType, "/")[1]

	resource, err := p.resourceRepo.GetResourceByMd5(MD5)
	if err != nil {
		return dto.UploadResponse{}, nil
	}
	if resource.ID != 0 {
		return dto.UploadResponse{ID: resource.ID, URL: resource.StaticPath}, nil
	}

	if _, err := os.Stat(constant.FileConfig.Path); os.IsNotExist(err) {
		os.Mkdir(constant.FileConfig.Path, os.ModePerm)
	}

	filetypePath := filepath.Join(constant.FileConfig.Path, contentType)
	if _, err := os.Stat(filetypePath); os.IsNotExist(err) {
		os.Mkdir(filetypePath, os.ModePerm)
	}

	filename := fmt.Sprintf("%s.%s", MD5, filetype)
	filePath := filepath.Join(filetypePath, filename)
	staticFp := fmt.Sprintf("%s/%s", constant.FileConfig.Static, contentType)
	staticPath := filepath.Join(staticFp, filename)

	resource, err = p.resourceRepo.AddResource(MD5, contentType, filePath, staticPath)
	if err != nil {
		return dto.UploadResponse{}, err
	}

	ctx.SaveUploadedFile(file, filePath)

	return dto.UploadResponse{ID: resource.ID, URL: resource.StaticPath}, nil
}
