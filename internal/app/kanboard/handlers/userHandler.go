package handlers

import (
	"server/internal/app/kanboard/dto"
	"server/internal/app/kanboard/services"
	"server/internal/common"
	"server/internal/utils"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService   *services.UserService
	publicService *services.PublicService
}

var userHandler *UserHandler

func NewUserHandler() *UserHandler {
	if userHandler == nil {
		userHandler = &UserHandler{
			userService:   services.NewUserService(),
			publicService: services.NewPublicService(),
		}
	}

	return userHandler
}

func (u UserHandler) Login(ctx *gin.Context) {
	var loginRequest dto.UserLoginRequest

	if err := utils.BindRequest(ctx, &loginRequest); err != nil {
		return
	}

	captchaResult := u.publicService.VerifyCaptcha(loginRequest.CaptchaID, loginRequest.CaptchaAnswer)
	if !captchaResult {
		common.Fail(ctx, common.RspOpts{
			Msg: "验证码错误或者过期",
		})
		return
	}

	data, token, err := u.userService.Login(loginRequest)
	if err != nil {
		common.Fail(ctx, common.RspOpts{
			Msg: err.Error(),
		})
		return
	}

	common.Ok(ctx, common.RspOpts{
		Msg: "登录成功",
		Data: gin.H{
			"token": token,
			"user":  data,
		},
	})
}

func (u UserHandler) Register(ctx *gin.Context) {
	var registerRequest dto.UserRegisterRequest

	if err := utils.BindRequest(ctx, &registerRequest); err != nil {
		return
	}

	captchaResult := u.publicService.VerifyCaptcha(registerRequest.CaptchaID, registerRequest.CaptchaAnswer)
	if !captchaResult {
		common.Fail(ctx, common.RspOpts{
			Msg: "验证码错误或者过期",
		})
		return
	}

	data, err := u.userService.RegisterKanboard(registerRequest)
	if err != nil {
		common.Fail(ctx, common.RspOpts{
			Msg: err.Error(),
		})
		return
	}

	common.Ok(ctx, common.RspOpts{
		Msg:  "注册成功",
		Data: data,
	})
}

func (u UserHandler) GetUserById(ctx *gin.Context) {
	var userIdRequest dto.UserIDRequest

	if err := utils.BindUri(ctx, &userIdRequest); err != nil {
		return
	}

	data, err := u.userService.GetUserInfo(userIdRequest)
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

func (u UserHandler) UpdateInfo(ctx *gin.Context) {
	var updateRequest dto.UserUpdateRequest
	var idRequest dto.UserIDRequest

	if err := utils.BindUri(ctx, &idRequest); err != nil {
		return
	}
	if err := utils.BindRequest(ctx, &updateRequest); err != nil {
		return
	}
	updateRequest.ID = idRequest.ID

	data, err := u.userService.UpdateInfo(updateRequest)
	if err != nil {
		common.Fail(ctx, common.RspOpts{
			Msg: err.Error(),
		})
		return
	}

	common.Ok(ctx, common.RspOpts{
		Msg:  "更新成功",
		Data: data,
	})
}

func (u UserHandler) UpdatePassword(ctx *gin.Context) {
	var updatePasswordRequest dto.UserUpdatePasswordRequest
	var idRequest dto.UserIDRequest

	if err := utils.BindUri(ctx, &idRequest); err != nil {
		return
	}
	if err := utils.BindRequest(ctx, &updatePasswordRequest); err != nil {
		return
	}
	updatePasswordRequest.ID = idRequest.ID

	if updatePasswordRequest.Current == "" || updatePasswordRequest.New == "" {
		common.Fail(ctx, common.RspOpts{
			Msg: "密码不能为空",
		})
		return
	}

	err := u.userService.UpdatePassword(updatePasswordRequest)
	if err != nil {
		common.Fail(ctx, common.RspOpts{
			Msg: err.Error(),
		})
		return
	}

	common.Ok(ctx, common.RspOpts{
		Msg: "修改成功",
	})
}

func (u UserHandler) GetStatistics(ctx *gin.Context) {
	var userIdRequest dto.UserIDRequest
	if utils.BindUri(ctx, &userIdRequest) != nil {
		return
	}

	data, err := u.userService.GetStatistics(userIdRequest.ID)
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

func (u UserHandler) GetCalendar(ctx *gin.Context) {
	var userIdRequest dto.UserIDRequest
	if utils.BindUri(ctx, &userIdRequest) != nil {
		return
	}

	data, err := u.userService.GetCalendar(userIdRequest.ID)
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
