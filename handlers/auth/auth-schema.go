package handlers

type ReqBodyAuthLogin struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,gte=5,lte=20"`
}

type ResBodyAuthLogin struct {
	User string `json:"user"`
	JWT  string `json:"jwt"`
}
