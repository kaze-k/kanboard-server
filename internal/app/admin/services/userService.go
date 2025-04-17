package services

import (
	"errors"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"server/internal/app/admin/dto"
	"server/internal/constant"
	"server/internal/repositories"

	"server/internal/global"
	"server/internal/models"
	md5 "server/pkg/MD5"
	"server/pkg/crypto"

	"github.com/gin-gonic/gin"
)

type UserService struct {
	userRepo     *repositories.UserRepo
	resourceRepo *repositories.ResourceRepo
}

var userService *UserService

func NewUserService() *UserService {
	if userService == nil {
		userService = &UserService{
			userRepo:     repositories.NewUserRepo(),
			resourceRepo: repositories.NewResourceRepo(),
		}
	}
	return userService
}

func (u *UserService) checkUserInfoExist(username, email, mobile string) (bool, error) {
	isUsernameExist := u.userRepo.CheckUserExistByName(username)
	if isUsernameExist {
		return isUsernameExist, errors.New("用户名已存在")
	}

	isEmailExist := u.userRepo.CheckEmailExist(email)
	if isEmailExist {
		return isEmailExist, errors.New("邮箱已存在")
	}

	isMobileExist := u.userRepo.CheckMobileExist(mobile)
	if isMobileExist {
		return isMobileExist, errors.New("手机号已存在")
	}

	return false, nil
}

func GenerateToken(id uint, username string) (string, error) {
	token, err := crypto.GenerateJWTToAdmin(id, username)
	global.Redis.Set(constant.ADMIN_TOKEN, strconv.Itoa(int(id)), token, constant.JWTConfig.AdminTokenExpiration)
	return token, err
}

func (u *UserService) Login(request dto.UserLoginRequest) (*dto.UserResponse, *string, error) {
	user, err := u.userRepo.GetUserByName(request.Username)
	if user.ID == 0 || err != nil {
		return nil, nil, errors.New("用户名或密码错误")
	}

	if user.IsAdmin == constant.NOT_ADMIN {
		return nil, nil, errors.New("权限不足")
	}

	comp := crypto.CheckPasswordHash(user.Password, request.Password)

	if !comp {
		return nil, nil, errors.New("用户名或密码错误")
	}

	token, err := GenerateToken(user.ID, user.Username)
	if err != nil {
		return nil, nil, errors.New("token生成失败")
	}

	var userResponse *dto.UserResponse
	resource, err := u.resourceRepo.GetResourceById(user.Avatar)
	if err != nil {
		return nil, nil, err
	}

	return userResponse.Set(user, resource), &token, nil
}

func (u *UserService) Register(request dto.UserRegisterRequest) (*dto.UserResponse, error) {
	isExist, err := u.checkUserInfoExist(request.Username, request.Email, request.Mobile)
	if isExist {
		return nil, err
	}

	createUser := models.User{
		Username:   request.Username,
		Password:   request.Password,
		Email:      request.Email,
		Mobile:     request.Mobile,
		Gender:     request.Gender,
		CreateFrom: constant.ADMIN,
	}
	user, err := u.userRepo.CreateUser(createUser)

	var userResponse *dto.UserResponse
	resource, err := u.resourceRepo.GetResourceById(user.Avatar)
	if err != nil {
		return nil, err
	}

	return userResponse.Set(user, resource), err
}

func (u *UserService) UploadAvatar(ctx *gin.Context, MD5 string, file *multipart.FileHeader) (uint, error) {
	fileMD5, err := md5.GetFileMD5(file)
	if err != nil {
		return 0, err
	}
	if MD5 != *fileMD5 {
		return 0, errors.New("md5校验失败")
	}

	contentType := file.Header.Get("Content-Type")
	filetype := strings.Split(contentType, "/")[1]

	resource, err := u.resourceRepo.GetResourceByMd5(MD5)
	if err != nil {
		return 0, nil
	}
	if resource.ID != 0 {
		return resource.ID, nil
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

	resource, err = u.resourceRepo.AddResource(MD5, contentType, filePath, staticPath)
	if err != nil {
		return 0, err
	}

	ctx.SaveUploadedFile(file, filePath)

	return resource.ID, nil
}

func (u *UserService) CreateUser(request dto.UserCreateRequest) (*dto.UserResponse, error) {
	isUsernameExist := u.userRepo.CheckUserExistByName(request.Username)
	if isUsernameExist {
		return nil, errors.New("用户名已存在")
	}

	if request.Email != nil {
		email := *request.Email
		if isExist := u.userRepo.CheckEmailExist(email); isExist {
			return nil, errors.New("邮箱已存在")
		}
	}

	if request.Mobile != nil {
		mobile := *request.Mobile
		if isMobileExist := u.userRepo.CheckMobileExist(mobile); isMobileExist {
			return nil, errors.New("手机号已存在")
		}
	}

	var createUser models.User
	createUser.CreateFrom = constant.ADMIN
	createUser.Username = request.Username
	createUser.Password = request.Password
	if request.Email != nil {
		createUser.Email = *request.Email
	}
	if request.Mobile != nil {
		createUser.Mobile = *request.Mobile
	}
	createUser.Gender = request.Gender
	createUser.Loginable = request.Loginable
	createUser.IsAdmin = request.IsAdmin
	if request.Position != nil {
		createUser.Position = *request.Position
	}
	if request.Avatar != nil {
		createUser.Avatar = *request.Avatar
	}
	user, err := u.userRepo.CreateUser(createUser)

	var userResponse *dto.UserResponse
	resource, err := u.resourceRepo.GetResourceById(user.Avatar)
	if err != nil {
		return nil, err
	}

	return userResponse.Set(user, resource), err
}

func (u *UserService) UpdateUser(request dto.UserUpdateRequest) (*dto.UserResponse, error) {
	if request.Email != nil {
		email := *request.Email
		if isExist := u.userRepo.CheckEmailExist(email); isExist {
			return nil, errors.New("邮箱已存在")
		}
	}

	if request.Mobile != nil {
		mobile := *request.Mobile
		if isMobileExist := u.userRepo.CheckMobileExist(mobile); isMobileExist {
			return nil, errors.New("手机号已存在")
		}
	}

	updateData := make(map[string]any)
	if request.Avatar != nil {
		updateData["avatar"] = request.Avatar
	}
	if request.Email != nil {
		updateData["email"] = request.Email
	}
	if request.Position != nil {
		updateData["position"] = request.Position
	}
	if request.Mobile != nil {
		updateData["mobile"] = request.Mobile
	}
	if request.Gender != nil {
		updateData["gender"] = request.Gender
	}
	if request.IsAdmin != nil {
		updateData["is_admin"] = request.IsAdmin
	}
	if request.Loginable != nil {
		updateData["loginable"] = request.Loginable
	}

	user, err := u.userRepo.UpdateUserById(updateData, request.ID)

	if user.Loginable == constant.NOT_LOGINABLE {
		global.Redis.Delete(constant.ADMIN_TOKEN, strconv.Itoa(int(user.ID)))
		global.Redis.Delete(constant.KANBOARD_TOKEN, strconv.Itoa(int(user.ID)))
	}

	var userResponse *dto.UserResponse
	resource, err := u.resourceRepo.GetResourceById(user.Avatar)
	if err != nil {
		return nil, err
	}

	return userResponse.Set(user, resource), err
}

func (u *UserService) GetUserList(request dto.PageRequest) (*dto.UserPageResponse, error) {
	total := u.userRepo.GetUserCount()
	users, err := u.userRepo.GetUserList(request.Page, request.PageSize)
	if err != nil {
		return nil, err
	}

	var pageResponse *dto.UserPageResponse

	return pageResponse.Set(total, request.Page, request.PageSize, users), nil
}

func (u *UserService) GetUser(request dto.UserIDRequest) (*dto.UserResponse, error) {
	user, err := u.userRepo.GetUserById(request.ID)
	if err != nil {
		return nil, err
	}

	var userResponse *dto.UserResponse
	resource, err := u.resourceRepo.GetResourceById(user.Avatar)
	if err != nil {
		return nil, err
	}

	return userResponse.Set(user, resource), nil
}

func (u *UserService) SearchUser(request dto.UserSearchRequest) (*dto.UserPageResponse, error) {
	query := make(map[string]any)

	if request.ID != nil {
		query["id"] = *request.ID
	}
	if request.Username != nil {
		query["username"] = *request.Username
	}
	if request.Loginable != nil {
		query["loginable"] = *request.Loginable
	}
	if request.IsAdmin != nil {
		query["is_admin"] = *request.IsAdmin
	}
	if request.CreateFrom != nil {
		query["create_from"] = *request.CreateFrom
	}
	if request.Position != nil {
		query["position"] = *request.Position
	}
	if request.Gender != nil {
		query["gender"] = *request.Gender
	}

	users, err := u.userRepo.GetUsers(query)
	if err != nil {
		return nil, err
	}
	total := int64(len(*users))

	var pageResponse *dto.UserPageResponse

	return pageResponse.Set(total, request.Page, request.PageSize, users), nil
}

func (u *UserService) DeleteUser(request dto.UserIDRequest) error {
	global.Redis.Delete(constant.ADMIN_TOKEN, strconv.Itoa(int(request.ID)))
	global.Redis.Delete(constant.KANBOARD_TOKEN, strconv.Itoa(int(request.ID)))
	return u.userRepo.DeleteUserById(request.ID)
}

func (u *UserService) ChangePassword(request dto.UserChangePasswordRequest) error {
	err := u.userRepo.UpdatePasswordById(request.NewPassword, request.ID)
	if err != nil {
		return err
	}
	global.Redis.Delete(constant.ADMIN_TOKEN, strconv.Itoa(int(request.ID)))
	global.Redis.Delete(constant.KANBOARD_TOKEN, strconv.Itoa(int(request.ID)))

	return err
}

func (u *UserService) GetMembers() ([]models.Member, error) {
	users, err := u.userRepo.GetAllUsers()
	if err != nil {
		return nil, err
	}

	var member models.Member
	var members []models.Member

	for _, user := range *users {
		member = member.Set(user)
		members = append(members, member)
	}

	return members, nil
}
