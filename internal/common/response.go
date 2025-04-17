package common

import (
	"net/http"
	"server/internal/constant"

	"github.com/gin-gonic/gin"
)

type BaseRspWithOmitempty struct {
	Status int    `json:"status,omitempty"`
	Code   int    `json:"code,omitempty"`
	Msg    string `json:"message,omitempty"`
	Data   any    `json:"data,omitempty"`
}

type BaseRsp struct {
	Status int    `json:"status"`
	Code   int    `json:"code"`
	Msg    string `json:"message"`
	Data   any    `json:"data"`
}

type RspOptsWithCode struct {
	Code int
	Msg  string
	Data any
}

type RspOpts struct {
	Msg  string
	Data any
}

func newBaseRspWithOmitempty(status int, code int, message string, data any) BaseRspWithOmitempty {
	return BaseRspWithOmitempty{
		Status: status,
		Code:   code,
		Msg:    message,
		Data:   data,
	}
}

func newBaseRsp(status int, code int, message string, data any) BaseRsp {
	return BaseRsp{
		Status: status,
		Code:   code,
		Msg:    message,
		Data:   data,
	}
}

func ResponseWithOmitempty(ctx *gin.Context, baseResponseWithOmitempty BaseRspWithOmitempty) {
	ctx.AbortWithStatusJSON(
		baseResponseWithOmitempty.Status,
		baseResponseWithOmitempty,
	)
}

func Response(ctx *gin.Context, baseResponse BaseRsp) {
	ctx.AbortWithStatusJSON(
		baseResponse.Status,
		baseResponse,
	)
}

func Ok(ctx *gin.Context, responseOpts RspOpts) {
	baseRsp := newBaseRsp(http.StatusOK, constant.OK, responseOpts.Msg, responseOpts.Data)

	ctx.AbortWithStatusJSON(
		baseRsp.Status,
		baseRsp,
	)
}

func Fail(ctx *gin.Context, responseOpts RspOpts) {
	baseRsp := newBaseRspWithOmitempty(http.StatusBadRequest, constant.FAIL, responseOpts.Msg, responseOpts.Data)

	ctx.AbortWithStatusJSON(
		baseRsp.Status,
		baseRsp,
	)
}

func NotFound(ctx *gin.Context, responseOpts RspOpts) {
	baseRsp := newBaseRspWithOmitempty(http.StatusNotFound, constant.FAIL, responseOpts.Msg, responseOpts.Data)

	ctx.AbortWithStatusJSON(
		baseRsp.Status,
		baseRsp,
	)
}

func ServerError(ctx *gin.Context, responseOpts RspOpts) {
	baseRsp := newBaseRspWithOmitempty(http.StatusInternalServerError, constant.FAIL, responseOpts.Msg, responseOpts.Data)

	ctx.AbortWithStatusJSON(
		baseRsp.Status,
		baseRsp,
	)
}

func BindRequest(ctx *gin.Context, request any) error {
	err := ctx.ShouldBind(&request)
	if err != nil {
		Fail(ctx, RspOpts{
			Msg: err.Error(),
		})
	}
	return err
}
