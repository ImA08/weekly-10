package models

import "mime/multipart"

type ProfileUser struct {
	FirstName   string `json:"firstname" form:"firstname"`
	LastName    string `json:"lastname" form:"lastname"`
	Status      string `json:"status" form:"status"`
	Points      int    `json:"points" form:"points"`
	PhoneNumber string `json:"phonenumber" form:"phonenumber"`
}

type ProfileUserForm struct {
	ProfileUser
	ProfilePicture *multipart.FileHeader `form:"profile_picture"`
}

type Profile struct {
	ProfileUser
	ProfilePicture string `json:"profile_picture"`
}
