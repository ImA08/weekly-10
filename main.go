package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	type userStruct struct {
		Email    string `json:"email" form:"email"`
		Password string `json:"password" form:"password"`
	}

	users := []userStruct{
		{Email: "caelus@trailblazer.hsr", Password: "lordtrashcane"},
		{Email: "march7th@trailblazer.hsr", Password: "caelussimper"},
	}

	// LOGIN

	r.POST("/tickitz/auth/login", func(ctx *gin.Context) {
		// Bind request body ke struct
		var loginData userStruct
		if err := ctx.ShouldBind(&loginData); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "Invalid input format",
			})
			return
		}

		// Validasi email dan password tidak kosong
		if loginData.Email == "" || loginData.Password == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "Email and password are required",
			})
			return
		}

		// Cari user dengan email yang cocok
		var foundUser *userStruct
		for _, user := range users {
			if user.Email == loginData.Email {
				foundUser = &user
				break
			}
		}

		// Jika user tidak ditemukan
		if foundUser == nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"message": "Invalid email or password",
			})
			return
		}

		// Verifikasi password
		if foundUser.Password != loginData.Password {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"message": "Invalid email or password",
			})
			return
		}

		// Jika login berhasil
		ctx.JSON(http.StatusOK, gin.H{
			"message": "Login successful",
			"user": gin.H{
				"email": foundUser.Email,
			},
		})
	})

	// REGISTER

	r.POST("tickitz/auth/signup", func(ctx *gin.Context) {
		newUser := &userStruct{}
		if err := ctx.ShouldBind(newUser); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "Invalid Input",
			})
			return
		}

		users = append(users, *newUser)
		ctx.JSON(http.StatusOK, gin.H{
			"msg":  "Succes",
			"data": users,
		})
	})

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

	mockMovies := []Movie{
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

	// upcoming movies

	r.GET("/tickitz/movies/upcoming", func(ctx *gin.Context) {
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
	})

	// popular movie

	r.GET("/tickitz/movies/popular", func(ctx *gin.Context) {
		var popularMovies []Movie

		for _, movie := range mockMovies {
			if movie.Rating > 3 {
				popularMovies = append(popularMovies, movie)
			}
		}

		if len(popularMovies) == 0 {
			ctx.JSON(http.StatusOK, gin.H{
				"status":  "succes",
				"message": "Movie tidak ditemukan",
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"status": "success",
			"data":   popularMovies,
		})

	})

	r.GET("/tickitz/movies", func(ctx *gin.Context) {
		movieQ := ctx.Query("movie")
		genresQ := []string{ctx.Query("genres")}

		if movieQ == "" {
			ctx.JSON(http.StatusOK, gin.H{
				"msg":  "Succes",
				"data": mockMovies,
			})
			return
		}

		// for _, movie := range mockMovies {
		// 	movie.Genres = strings.ToLower(movie.Genres)
		// 	if strings.ToLower(movie.Title) == strings.ToLower(movieQ) || slices.Equal(movie.Genres, genresQ) {
		// 		ctx.JSON(http.StatusOK, gin.H{
		// 			"status": "success",
		// 			"data":   movie,
		// 		})
		// 	}

		// }

	})

	r.Run()

}
