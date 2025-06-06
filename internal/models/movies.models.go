package models

import (
	"mime/multipart"
	"time"
)

type MovieStruct struct {
	ID          int       `json:"id" DB:"id"`
	Title       string    `json:"title" DB:"title"`
	Synopsis    string    `json:"synopsis" DB:"synopsis"`
	Duration    int       `json:"duration" DB:"duration"`
	Genre       string    `json:"genre" DB:"genre"`
	Image       []string  `json:"image" DB:"image"`
	ReleaseDate time.Time `json:"release_date" DB:"release_date"`
}

type UpcomingMovie struct {
	ID       int       `json:"id"`
	Title    string    `json:"title"`
	Synopsis string    `json:"synopsis"`
	Duration int       `json:"duration"`
	Cinemas  string    `json:"cinemas"` // Confirm this is string
	Location string    `json:"location"`
	City     string    `json:"city"`
	Schedule time.Time `json:"schedule"`
	Date     string    `json:"date"`
	Genres   []string  `json:"genres"`
}

type DetailMovieStruct struct {
	ID          int       `json:"id" DB:"id"`
	Title       string    `json:"title,binding:required" DB:"title" form:"title"`
	Genres      []string  `json:"genres" DB:"genres" form:"genres"`
	Synopsis    string    `json:"synopsis" DB:"synopsis" form:"synopsis"`
	Duration    int       `json:"duration" DB:"duration" form:"duration"`
	Casts       []string  `json:"casts" DB:"casts" form:"casts"`
	Directors   []string  `json:"directors" DB:"directors" form:"directors"`
	ReleaseDate time.Time `json:"release_date" DB:"release_date" form:"release"`
	CreatedAt   time.Time `json:"created_at" DB:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" DB:"updated_at"`
}

type MovieForm struct {
	DetailMovieStruct
	Poster   *multipart.FileHeader `form:"poster"`
	Backdrop *multipart.FileHeader `form:"backdrop"`
	Page     int                   `form:"page"`
}

type Movie struct {
	DetailMovieStruct
	Poster   *string `json:"poster"`
	Backdrop *string `json:"backdrop"`
}

type MoviesResponse struct {
	Page   int         `json:"page"`
	Movies []ShowMovie `json:"movies"`
}

type ShowMovie struct {
	Id     int      `json:"id"`
	Title  string   `json:"title"`
	Genres []string `json:"genres"`
	Image  *string  `json:"image"`
}

type Schedules struct {
	MovieName string     `json:"movie_name" form:"movie_name"`
	Cinema    string     `json:"cinema" DB:"cinema" form:"cinema"`
	Location  string     `json:"location"`
	City      string     `json:"city"`
	ShowTime  *time.Time `json:"show_time" DB:"show_time" form:"show_time"`
	Date      *time.Time `json:"date" DB:"date" form:"date"`
	Price     float64    `json:"price" DB:"price" form:"price"`
}
