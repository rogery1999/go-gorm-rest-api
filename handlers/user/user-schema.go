package handlers

type updateUserRequestBody struct {
	Name string `json:"name" validate:"required"`
	Age  uint   `json:"age" validate:"gte=1,lte=130"`
}

type createUserRequestBody struct {
	FirstName  string `json:"firstName" validate:"required,gte=2"`
	MiddleName string `json:"middleName" validate:"required,gte=2"`
	LastName   string `json:"lastName" validate:"required,gte=2"`
	Email      string `json:"email" validate:"required,email"`
	Password   string `json:"password" validate:"required,gte=5,lte=20"`
	Birthday   string `json:"birthday" validate:"required,datetime=2006-01-02"`
}
