package repositories

import (
	"time"

	"server/internal/constant"
	"server/internal/global"
	"server/internal/models"
	"server/internal/utils"

	"gorm.io/gorm"
)

type UserRepo struct {
	db *gorm.DB
}

var userRepo *UserRepo

func NewUserRepo() *UserRepo {
	if userRepo == nil {
		userRepo = &UserRepo{
			db: global.DB,
		}
	}
	return userRepo
}

func (u *UserRepo) CreateUser(user models.User) (*models.User, error) {
	err := u.db.Create(&user).Error
	return utils.HandleError(&user, err)
}

func (u *UserRepo) CheckUserExistByName(username string) bool {
	var user models.User
	count := u.db.Find(&user, "username = ?", username).RowsAffected
	return count > 0
}

func (u *UserRepo) CheckEmailExist(email string) bool {
	var user models.User
	count := u.db.Find(&user, "email = ?", email).RowsAffected
	return count > 0
}

func (u *UserRepo) CheckMobileExist(mobile string) bool {
	var user models.User
	count := u.db.Find(&user, "mobile = ?", mobile).RowsAffected
	return count > 0
}

func (u *UserRepo) GetUserById(id uint) (*models.User, error) {
	var user models.User
	err := u.db.Find(&user, "id = ?", id).First(&user).Error
	return utils.HandleError(&user, err)
}

func (u *UserRepo) GetUserByName(username string) (*models.User, error) {
	var user models.User
	err := u.db.Find(&user, "username = ? and loginable = ?", username, true).Error
	return utils.HandleError(&user, err)
}

func (u *UserRepo) GetUserCount() int64 {
	var count int64
	u.db.Model(&models.User{}).Count(&count)
	return count
}

func (u *UserRepo) GetUserList(page, pageSize int) (*[]models.User, error) {
	var users []models.User
	err := u.db.Limit(pageSize).Offset((page - 1) * pageSize).Find(&users).Error
	return utils.HandleError(&users, err)
}

func (u *UserRepo) GetUsers(query map[string]any) (*[]models.User, error) {
	var users []models.User
	var user models.User
	ctx := u.db.Model(&user)
	if id, ok := query["id"].(uint); ok {
		ctx.Where("id = ?", id)
	}
	if username, ok := query["username"].(string); ok {
		ctx.Where("username LIKE ?", "%"+username+"%")
	}
	if loginable, ok := query["loginable"].(bool); ok {
		ctx.Where("loginable = ?", loginable)
	}
	if is_admin, ok := query["is_admin"].(bool); ok {
		ctx.Where("is_admin = ?", is_admin)
	}
	if create_from, ok := query["create_from"].(constant.From); ok {
		ctx.Where("create_from = ?", create_from)
	}
	if position, ok := query["position"].(string); ok {
		ctx.Where("position = ?", position)
	}
	if gender, ok := query["gender"].(constant.Gender); ok {
		ctx.Where("gender = ?", gender)
	}
	err := ctx.Find(&users).Error
	return utils.HandleError(&users, err)
}

func (u *UserRepo) GetAllUsers() (*[]models.User, error) {
	var users []models.User
	err := u.db.Find(&users, "loginable = ?", true).Error
	return utils.HandleError(&users, err)
}

func (u *UserRepo) UpdateUserById(values map[string]any, id uint) (*models.User, error) {
	var user models.User
	if err := u.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	err := u.db.Model(&user).Where("id = ?", id).Updates(values).First(&user).Error
	return utils.HandleError(&user, err)
}

func (u *UserRepo) DeleteUserById(id uint) error {
	var user models.User
	if err := u.db.First(&user, id).Error; err != nil {
		return err
	}
	err := u.db.Delete(&user, "id = ?", id).Error
	return err
}

func (u *UserRepo) CheckIsAdminById(id uint) bool {
	var user models.User
	count := u.db.Find(&user, "id = ?", id, "is_admin = ?", constant.IS_ADMIN).RowsAffected
	return count > 0
}

func (u *UserRepo) UpdatePasswordById(password string, id uint) error {
	var user models.User
	if err := u.db.First(&user, id).Error; err != nil {
		return err
	}
	err := u.db.Find(&user, "id = ?", id).Update("password", password).Error
	return err
}

func (u *UserRepo) GetMaleCount() int64 {
	var count int64
	u.db.Model(&models.User{}).Where("gender = ?", constant.MALE).Count(&count)
	return count
}

func (u *UserRepo) GetFemaleCount() int64 {
	var count int64
	u.db.Model(&models.User{}).Where("gender = ?", constant.FEMALE).Count(&count)
	return count
}

func (u *UserRepo) GetCreateFromKanboardCount() int64 {
	var count int64
	u.db.Model(&models.User{}).Where("create_from = ?", constant.KANBOARD).Count(&count)
	return count
}

func (u *UserRepo) GetCreateFromAdminCount() int64 {
	var count int64
	u.db.Model(&models.User{}).Where("create_from = ?", constant.ADMIN).Count(&count)
	return count
}

func (u *UserRepo) GetAdminCount() int64 {
	var count int64
	u.db.Model(&models.User{}).Where("is_admin = ?", constant.IS_ADMIN).Count(&count)
	return count
}

func (u *UserRepo) GetLoginableCount() int64 {
	var count int64
	u.db.Model(&models.User{}).Where("loginable = ?", true).Count(&count)
	return count
}

func (u *UserRepo) GetUnLoginableCount() int64 {
	var count int64
	u.db.Model(&models.User{}).Where("loginable = ?", false).Count(&count)
	return count
}

func (u *UserRepo) GetAllUserCreateAt() (*[]time.Time, error) {
	createdAtData := []time.Time{}
	var user models.User
	err := u.db.Model(&user).Pluck("created_at", &createdAtData).Error
	return &createdAtData, err
}

func (u *UserRepo) GetAllAdminID() (*[]uint, error) {
	userIDData := []uint{}
	var user models.User
	err := u.db.Model(&user).Where("is_admin = ?", constant.IS_ADMIN).Pluck("id", &userIDData).Error
	return &userIDData, err
}
