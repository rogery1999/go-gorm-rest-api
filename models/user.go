package models

type User struct {
	Id   uint64 `json:"id"`
	Name string `json:"name"`
	Age  uint   `json:"age"`
}
