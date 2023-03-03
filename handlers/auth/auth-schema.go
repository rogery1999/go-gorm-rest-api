package handlers

type ReqBodyAuthLogin struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type ResBodyAuthLogin struct {
	Username string `json:"username"`
	JWT      string `json:"jwt"`
}
