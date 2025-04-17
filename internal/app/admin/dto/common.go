package dto

type PageRequest struct {
	Page     int `json:"page" form:"page" binding:"required"`
	PageSize int `json:"page_size" form:"page_size" binding:"required"`
}

type CaptchaResponse struct {
	Captcha string `form:"captcha" json:"captcha"`
	Id      string `form:"id" json:"id"`
}
