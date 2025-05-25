package models

import "time"

type OrderStruct struct {
	Id                int       `json:"id"`
	User_id           int       `json:"user_id"`
	Schedule_id       int       `json:"schedule_id" form:"schedule_id"`
	Seats             []string  `json:"seats" form:"seats"`
	Total             int       `json:"total"`
	Fullname          string    `json:"fullname" form:"fullname"`
	Email             string    `json:"email" form:"email"`
	Status_paid       bool      `json:"status_paid"`
	Payment_method_id int       `json:"payment_method_id" form:"payment_method_id"`
	Transaction_date  time.Time `json:"transaction_date" `
	Updated_at        time.Time `json:"updated_at"`
	Invoice_number    string    `json:"invoice_number"`
	PhoneNumber       string    `json:"phone_number" form:"phone_number"`
}
