package dto

type CaptchaResponse struct {
	Captcha string `form:"captcha" json:"captcha"`
	Id      string `form:"id" json:"id"`
}

type UploadResponse struct {
	ID  uint   `json:"id"`
	URL string `json:"URL"`
}
