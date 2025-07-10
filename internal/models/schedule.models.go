package models

import "time"

type Schedules struct {
	Cinema   string     `json:"cinema" DB:"cinema" form:"cinema"`
	ShowTime *time.Time `json:"show_time" DB:"show_time" form:"show_time"`
	Date     *time.Time `json:"date" DB:"date" form:"date"`
	Location string     `json:"location"`
	City     string     `json:"city"`
	Price    float64    `json:"price" DB:"price" form:"price"`
}

type ScheduleRequest struct {
	MovieId  int     `json:"movie_id" form:"movie_id"`
	CinemaId int     `json:"cinema_id" form:"cinema_id"`
	ShowTime string  `json:"show_time" DB:"show_time" form:"show_time"`
	Date     string  `json:"date" DB:"date" form:"date"`
	Price    float64 `json:"price" DB:"price" form:"price"`
}
