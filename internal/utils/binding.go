package utils

import (
	"server/internal/common"

	"github.com/gin-gonic/gin"
)

func BindRequest(ctx *gin.Context, obj any) error {
	if err := ctx.ShouldBind(obj); err != nil {
		common.Fail(ctx, common.RspOpts{
			Msg: err.Error(),
		})
		return err
	}

	return nil
}

func BindUri(ctx *gin.Context, obj any) error {
	if err := ctx.ShouldBindUri(obj); err != nil {
		common.Fail(ctx, common.RspOpts{
			Msg: err.Error(),
		})
		return err
	}
	return nil
}

func BindQuery(ctx *gin.Context, obj any) error {
	if err := ctx.ShouldBindQuery(obj); err != nil {
		common.Fail(ctx, common.RspOpts{
			Msg: err.Error(),
		})
		return err
	}
	return nil
}
