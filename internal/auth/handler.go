package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	tokenService TokenService
	username     string
	password     string
}

type loginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func NewHandler(tokenService TokenService, username string, password string) Handler {
	return Handler{
		tokenService: tokenService,
		username:     username,
		password:     password,
	}
}

func (h Handler) Login(ctx *gin.Context) {
	var request loginRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "username and password are required"})
		return
	}

	if request.Username != h.username || request.Password != h.password {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	token, err := h.tokenService.Generate(request.Username, []string{"users:read"})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"token_type": "Bearer",
		"token":      token,
	})
}

func (h Handler) Valid(ctx *gin.Context) {
	subject, _ := ctx.Get("subject")
	permissions, _ := ctx.Get("permissions")
	ctx.JSON(http.StatusOK, gin.H{
		"valid":       true,
		"subject":     subject,
		"permissions": permissions,
	})
}

func (h Handler) Read(ctx *gin.Context) {
	subject, _ := ctx.Get("subject")
	permissions, _ := ctx.Get("permissions")
	ctx.JSON(http.StatusOK, gin.H{
		"subject":     subject,
		"permissions": permissions,
	})
}
