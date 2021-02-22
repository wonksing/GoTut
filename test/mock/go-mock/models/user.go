package models

type User struct {
	Name   string `json:"name"`
	Age    int16  `json:"age"`
	Gender string `json:"gender"`
}
