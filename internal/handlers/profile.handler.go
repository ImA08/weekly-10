package handlers

import (
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"path/filepath"

	"time"

	"github.com/gin-gonic/gin"
	"minitask1.go/internal/models"
	"minitask1.go/internal/repositories"
	"minitask1.go/pkg"
)

type ProfileHandler struct {
	profileRepo *repositories.ProfileRepository
}

func NewProfileHandler(profileRepo *repositories.ProfileRepository) *ProfileHandler {
	return &ProfileHandler{profileRepo: profileRepo}
}

func (h *ProfileHandler) GetProfile(ctx *gin.Context) {
	claims, _ := ctx.Get("Payload")

	userClaims := claims.(*pkg.Claims)
	log.Println("[DEBUG USERCLAIMS PROFILE]", userClaims)
	// userID := userClaims.Id

	user, err := h.profileRepo.GetUserProfile(ctx.Request.Context(), userClaims.Id)
	if err != nil {
		log.Println("[DEBUG USER]", user, err)
		ctx.JSON(http.StatusNotFound, gin.H{
			"msg": "User Not found",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "Succes",
		"User": user,
	})
}

func (h *ProfileHandler) EditProfile(ctx *gin.Context) {
	// 1. Get user ID from JWT claims
	claims, exists := ctx.Get("Payload")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userClaims := claims.(*pkg.Claims)
	userID := userClaims.Id

	// 2. Bind form data (including multipart form)
	var formBody models.ProfileUserForm
	if err := ctx.ShouldBind(&formBody); err != nil {
		log.Printf("Binding error: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format",
			// "details": err.Error(),
		})
		return
	}

	// 3. Handle profile picture upload if exists
	var profilePictureURL string
	if formBody.ProfilePicture != nil {
		filename, filePath, err := h.handleFileUpload(ctx, formBody.ProfilePicture, userID)
		if err != nil {
			log.Printf("File upload error: %v", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to upload profile picture",
			})
			return
		}
		log.Println(filename)
		profilePictureURL = filePath
	}

	// 5. Update profile in database
	result, err := h.profileRepo.UpdateUserProfile(ctx.Request.Context(), userID, formBody, profilePictureURL)
	if err != nil {
		log.Printf("Database error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to update profile",
			"details": err.Error(),
		})
		return
	}

	// 6. Return updated profile
	// updatedProfile, err := h.profileRepo.GetUserProfile(ctx.Request.Context(), userID)
	// if err != nil {
	// 	log.Printf("Fetch error: %v", err)
	// 	ctx.JSON(http.StatusInternalServerError, gin.H{
	// 		"error": "Profile updated but failed to fetch updated data",
	// 	})
	// 	return
	// }

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Profile updated successfully",
		"data":    result,
	})
}

func (h *ProfileHandler) handleFileUpload(ctx *gin.Context, file *multipart.FileHeader, userID int) (filename, filePath string, err error) {
	ext := filepath.Ext(file.Filename)
	filename = fmt.Sprintf("%d_%d_profile%s", time.Now().UnixNano(), userID, ext)
	filePath = filepath.Join("public", "img", filename)

	if err := ctx.SaveUploadedFile(file, filePath); err != nil {
		return "", "", fmt.Errorf("failed to save file: %w", err)
	}

	// Return URL path instead of filesystem path
	return filename, "/img/" + filename, nil
}
