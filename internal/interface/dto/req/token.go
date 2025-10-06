package req

type UpdateToken struct {
	Token string `json:"token" binding:"required"`
}
