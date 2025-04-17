package handlers

import (
	"server/internal/app/admin/dto"
	"server/internal/app/admin/services"
	"server/internal/common"
	"server/internal/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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

	data, err := u.userService.Register(registerRequest)
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

func (u UserHandler) UploadAvatar(ctx *gin.Context) {
	file, err := ctx.FormFile("file")
	if err != nil {
		common.Fail(ctx, common.RspOpts{
			Msg: err.Error(),
		})
		return
	}

	md5 := ctx.PostForm("md5")

	data, err := u.userService.UploadAvatar(ctx, md5, file)
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

func (u UserHandler) CreateUser(ctx *gin.Context) {
	var createRequest dto.UserCreateRequest

	if err := utils.BindRequest(ctx, &createRequest); err != nil {
		return
	}

	data, err := u.userService.CreateUser(createRequest)
	if err != nil {
		common.Fail(ctx, common.RspOpts{
			Msg: err.Error(),
		})
		return
	}

	common.Ok(ctx, common.RspOpts{
		Msg:  "创建用户成功",
		Data: data,
	})
}

func (u UserHandler) SearchUser(ctx *gin.Context) {
	var searchRequest dto.UserSearchRequest

	if err := utils.BindRequest(ctx, &searchRequest); err != nil {
		return
	}

	data, err := u.userService.SearchUser(searchRequest)
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

func (u UserHandler) GetUsers(ctx *gin.Context) {
	var request dto.PageRequest

	if err := utils.BindQuery(ctx, &request); err != nil {
		return
	}

	data, err := u.userService.GetUserList(request)
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

func (u UserHandler) GetUser(ctx *gin.Context) {
	var request dto.UserIDRequest

	if err := utils.BindRequest(ctx, &request); err != nil {
		return
	}

	data, err := u.userService.GetUser(request)
	if err == gorm.ErrRecordNotFound {
		common.NotFound(ctx, common.RspOpts{
			Msg: err.Error(),
		})
		return
	}

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

func (u UserHandler) UpdateUser(ctx *gin.Context) {
	var updateUserRequest dto.UserUpdateRequest

	if err := utils.BindRequest(ctx, &updateUserRequest); err != nil {
		return
	}

	data, err := u.userService.UpdateUser(updateUserRequest)
	if err != nil {
		common.Fail(ctx, common.RspOpts{
			Msg: err.Error(),
		})
		return
	}

	common.Ok(ctx, common.RspOpts{
		Msg:  "更新用户成功",
		Data: data,
	})
}

func (u UserHandler) DeleteUser(ctx *gin.Context) {
	var idDto dto.UserIDRequest

	if err := utils.BindRequest(ctx, &idDto); err != nil {
		return
	}

	err := u.userService.DeleteUser(idDto)
	if err != nil {
		common.Fail(ctx, common.RspOpts{
			Msg: err.Error(),
		})
		return
	}

	common.Ok(ctx, common.RspOpts{
		Msg: "删除成功",
	})
}

func (u UserHandler) ChangePassword(ctx *gin.Context) {
	var request dto.UserChangePasswordRequest

	if err := utils.BindRequest(ctx, &request); err != nil {
		return
	}

	err := u.userService.ChangePassword(request)
	if err != nil {
		common.Fail(ctx, common.RspOpts{
			Msg: err.Error(),
		})
		return
	}

	common.Ok(ctx, common.RspOpts{
		Msg: "修改密码成功",
	})
}

func (u UserHandler) GetMembers(ctx *gin.Context) {
	members, err := u.userService.GetMembers()
	if err != nil {
		common.Fail(ctx, common.RspOpts{
			Msg: err.Error(),
		})
		return
	}

	common.Ok(ctx, common.RspOpts{
		Data: members,
	})
}
