package handlers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/oauth2"
	"minitask1.go/internal/models"
	"minitask1.go/internal/repositories"
	"minitask1.go/pkg"
)

type AuthHandler struct {
	userRepo *repositories.UserRepository
}

func NewAuthHandler(userRepo *repositories.UserRepository) *AuthHandler {
	return &AuthHandler{userRepo: userRepo}
}

type UserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type UserResponse struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

// Register
// @summary			Register User
// @Router			/auth/register [post]
// @accept 		 	json
// @param			body body models.AuthForm true "register information"
// @produce			json
// @failure			500 {object} models.ErrorResponse
// @failure			400 {object} models.ErrorResponse
// @succes			201 {object} gin.H
func (h *AuthHandler) Register(ctx *gin.Context) {

	var body models.UserStruct
	if err := ctx.ShouldBind(&body); err != nil {
		log.Println(err.Error())
		if strings.Contains(err.Error(), "Field validation") {
			if strings.Contains(err.Error(), "min") {
				ctx.JSON(http.StatusBadRequest, gin.H{
					"message": "Password have to greater than 8 characters",
				})
				return
			}
			if strings.Contains(err.Error(), "max") {
				ctx.JSON(http.StatusBadRequest, gin.H{
					"message": "Password have to less than 10 characters",
				})
				return
			}
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "Name and password required!",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal server Error",
		})
		return
	}

	hash := pkg.InitHashConfig()
	hash.UseDefaultConfig()
	hashedPass, err := hash.GenPasswordHash(body.Password)
	if err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Hash Failed",
		})
		return
	}

	cmd, err := h.userRepo.CreateUser(ctx.Request.Context(), body.Email, hashedPass)
	if err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal server Error",
		})
		return
	}
	if cmd.RowsAffected() == 0 {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal server Error",
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": fmt.Sprintf("Selamat datang %s silahkan login", body.Email),
	})

	// var req UserRequest

	// if err := ctx.ShouldBindJSON(&req); err != nil {
	// 	ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
	// 	return
	// }

	// hashedPassword, err := bcrypt.GenerateFromPassword(
	// 	[]byte(req.Password),
	// 	bcrypt.DefaultCost,
	// )
	// if err != nil {
	// 	ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to secure password"})
	// 	return
	// }

	// user, err := h.userRepo.CreateUser(ctx.Request.Context(), req.Email, string(hashedPassword))
	// if err != nil {
	// 	ctx.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
	// 	return
	// }

	// // Return
	// ctx.JSON(http.StatusOK, gin.H{
	// 	"ID":    user.ID,
	// 	"Email": user.Email,
	// 	"msg":   "SignUp Succes",
	// })
}

// Login
// @summary			Login User
// @Router			/auth [post]
// @param			body body models.AuthForm true "login information"
// @accept 			json
// @produce			json
// @failure			500 {object} models.ErrorResponse
// @succes			201 {object} models.Response
func (h *AuthHandler) Login(ctx *gin.Context) {
	// var req UserRequest
	var body models.UserStruct

	// 1. Input validation
	if err := ctx.ShouldBindJSON(&body); err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	// 2. Get user from repository
	user, err := h.userRepo.LogInUserRepo(ctx.Request.Context(), body.Email)

	if err != nil {

		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error":  "Authentication failed",
			"detail": "Email not found or invalid",
		})
		return
	}

	hash := pkg.InitHashConfig()
	log.Println("[DEBUG]", user.Password, body.Password)

	valid, err := hash.CompareHashAndPassword(user.Password, body.Password)
	if err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal server Error",
		})
		return
	}

	log.Println("[DEBUG] FIRST TEST")
	if !valid {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": "invalid username/password",
		})
		return
	}
	log.Println("[DEBUG] SECOND TEST")
	claims := pkg.NewClaims(user.ID, user.Role)
	log.Println("[DEBUG] USER CHECK", claims)
	token, err := claims.GenerateToken()
	if err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal server error",
		})
		return
	}

	decoded, err := jwt.ParseWithClaims(token, &pkg.Claims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		log.Printf("[DEBUG] Token verification failed: %v", err)
	} else if claims, ok := decoded.Claims.(*pkg.Claims); ok {
		log.Printf("[DEBUG] Generated token contains: ID=%d, Role=%s", claims.Id, claims.Role)
	}
	// 4. Successful response
	ctx.JSON(http.StatusOK, gin.H{
		"ID":    user.ID,
		"Email": user.Email,
		"token": token,
		"msg":   "LogIn Succes",
	})
}

// OAuth2

var (
	googleOauthConfig *oauth2.Config
	oauthStateString  string = "iisme-cutecumber"
	key                      = "kodakofi_auth_key"
	maxAge                   = 60 * 60 * 24
	isProd                   = false
)

// func (h *AuthHandler) Init(ctx *gin.Context) {
// 	err := godotenv.Load()

// 	if err != nil {
// 		log.Printf("ERROR loading .env file, assuming variabel set %v", err)
// 	}

// 	googleOauthConfig = &oauth2.Config{
// 		RedirectURL:  "http://localhost:8080/auth/google/callback",
// 		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
// 		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
// 		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
// 		Endpoint:     google.Endpoint,
// 	}

// 	if googleOauthConfig.ClientID == "" || googleOauthConfig.ClientSecret == "" {
// 		log.Fatal("GOOGLE_CLIENT_ID and GOOGLE_CLIENT_SECRET must be set in .env or as environment variables")
// 	}

// 	// Clear any existing providers to avoid conflicts
// 	goth.ClearProviders()

// 	// Setup cookie store for sessions
// 	store := sessions.NewCookieStore([]byte(key))
// 	store.MaxAge(maxAge)
// 	store.Options = &sessions.Options{
// 		Path:     "/",
// 		HttpOnly: true,
// 		Secure:   isProd,
// 	}

// 	// Set the session store for gothic
// 	gothic.Store = store

// 	// Register the Google provider
// 	provider := google.New(
// 		googleOauthConfig.ClientID,
// 		googleOauthConfig.ClientSecret,
// 		googleOauthConfig.RedirectURL,
// 		"email", "profile",
// 	)
// 	goth.UseProviders(provider)

// 	log.Printf("Google OAuth provider registered successfully")
// }

// // GoogleLogin initiates OAuth flow with Google
// func (a *AuthHandler) GoogleLogin(ctx *gin.Context) {
// 	// Clear session to prevent stale data
// 	session, _ := gothic.Store.Get(ctx.Request, gothic.SessionName)
// 	session.Values = map[any]any{}
// 	session.Save(ctx.Request, ctx.Writer)

// 	// Set the provider explicitly
// 	ctx.Request.URL.RawQuery = "provider=google"

// 	log.Printf("Starting Google OAuth flow with URL: %s", ctx.Request.URL.String())

// 	// Begin the OAuth flow
// 	gothic.BeginAuthHandler(ctx.Writer, ctx.Request)
// }

// // GoogleCallback handles the callback from Google OAuth
// func (a *AuthHandler) GoogleCallback(ctx *gin.Context) {
// 	response := models.NewResponse(ctx)

// 	// Ensure provider is set in query parameters
// 	q := ctx.Request.URL.Query()
// 	q.Set("provider", "google")
// 	ctx.Request.URL.RawQuery = q.Encode()

// 	log.Printf("Processing OAuth callback for provider: google")

// 	// Complete the OAuth flow
// 	user, err := gothic.CompleteUserAuth(ctx.Writer, ctx.Request)
// 	if err != nil {
// 		log.Printf("Error during OAuth callback: %v", err)

// 		// Return the error as JSON response instead of redirecting
// 		response.BadRequest("Authentication failed", err.Error())
// 		return
// 	}

// 	log.Printf("Google auth successful for email: %s", user.Email)

// 	// Check if user exists in our database
// 	result, err := a.userRepo.LogInUserRepo(ctx.Request.Context(), user.Email)
// 	if err != nil {
// 		// If there's an error but it's not because the user doesn't exist
// 		if err.Error() != "no rows in result set" {
// 			response.InternalServerError("Failed to check user", err.Error())
// 			return
// 		}
// 	}

// 	// User does not exist, register them
// 	if result.Email == "" {
// 		log.Printf("Registering new user from Google OAuth: %s", user.Email)

// 		// Generate a random password for OAuth users
// 		hash := pkg.InitHashConfig()
// 		hash.UseDefaultConfig()
// 		randomPass := utils.GenerateRandomString(12)

// 		hashedPass, err := hash.GenPasswordHash(randomPass)
// 		if err != nil {
// 			response.InternalServerError("Failed to hash password", err.Error())
// 			return
// 		}

// 		// Create new user request with name from profile
// 		userReq := models.AuthForm{
// 			Email:    user.Email,
// 			Password: randomPass,
// 		}

// 		// Register the user
// 		newUser, err := a.userRepo.CreateUser(ctx.Request.Context(), userReq.Email, hashedPass)
// 		if err != nil {
// 			response.InternalServerError("Failed to register user", err.Error())
// 			return
// 		}

// 		// Since this is OAuth, we can mark them as verified directly
// 		// if err := a.repo.UpdateUserVerificationStatus(ctx.Request.Context(), newUser.AuthID); err != nil {
// 		// 	response.InternalServerError("Failed to verify user", err.Error())
// 		// 	return
// 		// }

// 		log.Printf("Successfully registered user from Google OAuth: %s", user.Email)

// 		// Update result to use the newly created user
// 		result = newUser
// 	} else {
// 		log.Printf("Existing user logged in with Google OAuth: %s", user.Email)
// 	}

// 	// Generate JWT token for the user
// 	payload := pkg.NewJWT(result.AuthID, result.Email, result.Role)
// 	token, err := payload.GenerateToken()
// 	if err != nil {
// 		response.InternalServerError("Failed to generate token", err.Error())
// 		return
// 	}

// 	// Check if this is API call or a redirect
// 	if ctx.Request.Header.Get("X-Requested-With") == "XMLHttpRequest" {
// 		// API call, return JSON with id field
// 		response.Success("Google authentication successful", map[string]string{
// 			"token": token,
// 			"email": user.Email,
// 			"name":  user.Name,
// 			"id":    result.AuthID, // Include the real user ID
// 			"role":  result.Role,
// 		})
// 	} else {
// 		// Browser flow, redirect to frontend with token, user ID, and role
// 		frontendURL := os.Getenv("FRONTEND_URL")
// 		if frontendURL == "" {
// 			frontendURL = "http://localhost:5173"
// 		}

// 		redirectURL := fmt.Sprintf("%s/auth/callback?token=%s&email=%s&name=%s&id=%s&role=%s",
// 			frontendURL,
// 			url.QueryEscape(token),
// 			url.QueryEscape(user.Email),
// 			url.QueryEscape(user.Name),
// 			url.QueryEscape(result.AuthID),
// 			url.QueryEscape(result.Role))

// 		log.Printf("Redirecting to: %s", redirectURL)
// 		ctx.Redirect(http.StatusFound, redirectURL)
// 	}
// }

// // TokenLogin handles login with a Google ID token
// func (a *AuthHandler) TokenLogin(ctx *gin.Context) {
// 	var req struct {
// 		IDToken string `json:"id_token" binding:"required"`
// 	}
// 	response := models.NewResponse(ctx)

// 	if err := ctx.ShouldBindJSON(&req); err != nil {
// 		response.BadRequest("Invalid input", err.Error())
// 		return
// 	}

// 	// Validate the token
// 	userInfo, err := utils.ValidateGoogleIDToken(req.IDToken)
// 	if err != nil {
// 		response.BadRequest("Invalid token", err.Error())
// 		return
// 	}

// 	// Check if user exists in our database
// 	result, err := a.repo.Login(ctx.Request.Context(), userInfo.Email)
// 	if err != nil {
// 		// If there's an error but it's not because the user doesn't exist
// 		if err.Error() != "no rows in result set" {
// 			response.InternalServerError("Failed to check user", err.Error())
// 			return
// 		}
// 	}

// 	// User does not exist, register them
// 	if result.Email == "" {
// 		// Generate a random password for OAuth users
// 		hash := pkg.InitHashConfig()
// 		hash.UseDefaultConfig()
// 		randomPass := utils.GenerateRandomString(12)

// 		hashedPass, err := hash.GenHashedPassword(randomPass)
// 		if err != nil {
// 			response.InternalServerError("Failed to hash password", err.Error())
// 			return
// 		}

// 		// Create new user request
// 		userReq := models.UserReq{
// 			Email:    userInfo.Email,
// 			Password: randomPass,
// 			Fullname: userInfo.Name, // Use name from Google token
// 		}

// 		// Register the user
// 		newUser, _, err := a.repo.Register(ctx.Request.Context(), userReq, hashedPass)
// 		if err != nil {
// 			response.InternalServerError("Failed to register user", err.Error())
// 			return
// 		}

// 		// Since this is OAuth, we can mark them as verified directly
// 		if err := a.repo.UpdateUserVerificationStatus(ctx.Request.Context(), newUser.AuthID); err != nil {
// 			response.InternalServerError("Failed to verify user", err.Error())
// 			return
// 		}

// 		// Update result to use the newly created user
// 		result = newUser
// 	}

// 	// Generate JWT token for the user
// 	payload := pkg.NewJWT(result.AuthID, result.Email, result.Role)
// 	token, err := payload.GenerateToken()
// 	if err != nil {
// 		response.InternalServerError("Failed to generate token", err.Error())
// 		return
// 	}

// 	response.Success("Google authentication successful", map[string]any{
// 		"token": token,
// 		"user": map[string]any{
// 			"email":   userInfo.Email,
// 			"name":    userInfo.Name,
// 			"picture": userInfo.Picture,
// 		},
// 	})
// }
