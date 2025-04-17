package models

type Member struct {
	UserID   uint   `json:"user_id" form:"user_id" required:"true"`
	Username string `json:"username" form:"username" required:"true"`
	Avatar   string `json:"avatar"`
}

func (m *Member) Set(user User) Member {
	return Member{
		UserID:   user.ID,
		Username: user.Username,
	}
}
