package handlers

type updateUserRequestBody struct {
	Name string `json:"name" validate:"required"`
	Age  uint   `json:"age" validate:"gte=1,lte=130"`
}

type createUserRequestBody struct {
	Name string `json:"name" validate:"required,alpha"`
	Age  uint   `json:"age" validate:"required,gte=1,lte=130"`
}
