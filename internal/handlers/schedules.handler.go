package handlers

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"minitask1.go/internal/models"
	"minitask1.go/internal/repositories"
)

type ScheduleHandler struct {
	scheduleRepo *repositories.ScheduleRepository
}

func NewScheduleHandler(scheduleRepo *repositories.ScheduleRepository) *ScheduleHandler {
	return &ScheduleHandler{scheduleRepo: scheduleRepo}
}
func (h *ScheduleHandler) GetScheduleHandler(ctx *gin.Context) {
	id := ctx.Param("id")

	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "movie ID is required",
		})
		return
	}

	idInt, err := strconv.Atoi(id)
	if err != nil || idInt <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid movie ID format",
		})
		return
	}

	fmt.Println("[DEBUG] ID ", idInt)

	result, err := h.scheduleRepo.GetScheduleMovieRepo(ctx, idInt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, gin.H{
				"message": "movie not found",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data":    result,
	})
}

func (h *ScheduleHandler) CreateScheduleHandler(ctx *gin.Context) {
	var newSchedule models.ScheduleRequest

	if err := ctx.ShouldBindJSON(&newSchedule); err != nil {
		log.Printf("Binding Error %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format",
		})
		return
	}

	id, err := h.scheduleRepo.CreateScheduleRepo(ctx, newSchedule)
	if err != nil {
		log.Printf("[ERROR] Database Error %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to add schedule",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message":     "Schedule add successfully",
		"Schedule ID": id,
	})

}
