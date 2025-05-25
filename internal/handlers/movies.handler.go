package handlers

import (
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"minitask1.go/internal/models"
	"minitask1.go/internal/repositories"
)

type MovieHandler struct {
	movieRepo *repositories.MoviesRepository
}

func NewMovieHandler(movieRepo *repositories.MoviesRepository) *MovieHandler {
	return &MovieHandler{movieRepo: movieRepo}
}

func (h *MovieHandler) GetMovies(ctx *gin.Context) {
	pageStr := ctx.DefaultQuery("page", "1")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page number"})
		return
	}

	var movies models.ShowMovie
	result, err := h.movieRepo.GetAllMovies(ctx.Request.Context(), movies, page)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch movies",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, models.MoviesResponse{
		Page:   page,
		Movies: result,
	})
}

func (h *MovieHandler) GetUpcomingMovies(ctx *gin.Context) {
	movies, err := h.movieRepo.UpcomingMovies(ctx.Request.Context())
	if err != nil {
		log.Printf("Failed to fetch upcoming movies: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch upcoming movies",
			"details": err.Error(),
		})
		return
	}

	if len(movies) == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "No upcoming movies found",
			"movies":  []models.UpcomingMovie{},
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"movies": movies,
	})
}

func (h *MovieHandler) AddMovieHandler(ctx *gin.Context) {
	// 1. Get user ID from JWT claims

	// pageStr := ctx.DefaultQuery("page", "1")
	// page, err := strconv.Atoi(pageStr)
	// if err != nil || page < 1 {
	// 	ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page number"})
	// 	return
	// }

	var formBody models.MovieForm
	if err := ctx.ShouldBind(&formBody); err != nil {
		log.Printf("Binding Error %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format",
			// "details": err.Error(),
		})
		return
	}

	var posterURL string
	var BackdropURL string
	if formBody.Poster != nil {
		posterFilename, posterFilePath, err := h.handleFileUpload(ctx, formBody.Poster)
		if err != nil {
			log.Printf("File upload error: %v", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to upload profile picture",
			})
			return
		}
		log.Println(posterFilename)
		posterURL = posterFilePath
	}

	if formBody.Backdrop != nil {
		backdropFilename, backdropFilePath, err := h.handleFileUpload(ctx, formBody.Poster)
		if err != nil {
			log.Printf("File upload error: %v", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to upload profile picture",
			})
			return
		}
		log.Println(backdropFilename)
		posterURL = backdropFilePath
	}

	err := h.movieRepo.AddMovies(ctx.Request.Context(), formBody, posterURL, BackdropURL)
	if err != nil {
		log.Printf("Database error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to update profile",
			"details": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Movie add successfully",
	})
}

func (h *MovieHandler) handleFileUpload(ctx *gin.Context, file *multipart.FileHeader) (filename, filePath string, err error) {
	ext := filepath.Ext(file.Filename)
	filename = fmt.Sprintf("%d_movie%s", time.Now().UnixNano(), ext)
	filePath = filepath.Join("public", "img", filename)

	if err := ctx.SaveUploadedFile(file, filePath); err != nil {
		return "", "", fmt.Errorf("failed to save file: %w", err)
	}

	// Return URL path instead of filesystem path
	return filename, "/img/" + filename, nil
}
