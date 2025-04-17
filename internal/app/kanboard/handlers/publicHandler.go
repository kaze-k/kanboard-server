package handlers

import (
	"server/internal/app/kanboard/services"
	"server/internal/common"

	"github.com/gin-gonic/gin"
)

type PublicHandler struct {
	publicService *services.PublicService
}

var publicHandler *PublicHandler

func NewPublicHandler() *PublicHandler {
	if publicHandler == nil {
		publicHandler = &PublicHandler{
			publicService: services.NewPublicService(),
		}
	}

	return publicHandler
}

func (p PublicHandler) GetCaptcha(ctx *gin.Context) {
	data, err := p.publicService.GetCaptcha()
	if err != nil {
		common.Fail(ctx, common.RspOpts{
			Msg: err.Error(),
		})
		return
	}

	common.Ok(ctx, common.RspOpts{
		Data: data,
	})
}

func (p PublicHandler) Upload(ctx *gin.Context) {
	file, err := ctx.FormFile("file")
	if err != nil {
		common.Fail(ctx, common.RspOpts{
			Msg: err.Error(),
		})
		return
	}

	md5 := ctx.PostForm("md5")

	data, err := p.publicService.Upload(ctx, md5, file)
	if err != nil {
		common.Fail(ctx, common.RspOpts{
			Msg: err.Error(),
		})
		return
	}

	common.Ok(ctx, common.RspOpts{
		Data: data,
	})
}
