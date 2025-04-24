package router

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Movie struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	ReleaseDate time.Time `json:"release_date"`
	PosterURL   string    `json:"poster_url"`
	Duration    int       `json:"duration"` // dalam menit
	Genres      []string  `json:"genres"`
	Rating      float32   `json:"rating"`
}

var mockMovies = []Movie{
	{
		ID:          1,
		Title:       "Avengers: Secret Wars",
		Description: "lorem ipsum",
		ReleaseDate: time.Now().AddDate(0, 3, 0), // 3 bulan dari sekarang
		PosterURL:   "https://example.com/posters/avengers.jpg",
		Duration:    150,
		Genres:      []string{"Action", "Drama"},
		Rating:      4.5,
	},
	{
		ID:          2,
		Title:       "Dune: Part Two",
		Description: "lorem ipsum",
		ReleaseDate: time.Now().AddDate(0, 2, 15), // 2 bulan 15 hari dari sekarang
		PosterURL:   "https://example.com/posters/dune2.jpg",
		Duration:    165,
		Genres:      []string{"Sci-Fi", "Action"},
		Rating:      2.9,
	},
	{
		ID:          3,
		Title:       "The Batman Part II",
		Description: "lorem ipsum",
		ReleaseDate: time.Now().AddDate(0, 5, 0), // 5 bulan dari sekarang
		PosterURL:   "https://example.com/posters/batman2.jpg",
		Duration:    180,
		Genres:      []string{"Action", "Drama"},
		Rating:      4.8,
	},
}

func getPopularMovies(ctx *gin.Context) {
	// Filter hanya film yang release date-nya di masa depan
	var upcomingMovies []Movie
	currentDate := time.Now()

	for _, movie := range mockMovies {
		if movie.ReleaseDate.After(currentDate) {
			upcomingMovies = append(upcomingMovies, movie)
		}
	}

	// Jika tidak ada upcoming movies
	if len(upcomingMovies) == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": "Movie tidak ditemukan",
			"data":    []Movie{},
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   upcomingMovies,
	})
}
