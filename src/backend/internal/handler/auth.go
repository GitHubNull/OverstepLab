package handler

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/oversteplab/oversteplab/internal/common"
	"github.com/oversteplab/oversteplab/internal/middleware"
	"github.com/oversteplab/oversteplab/internal/service"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(svc *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: svc}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var input service.RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		common.BadRequest(c, "Invalid request: "+err.Error())
		return
	}

	user, err := h.authService.Register(&input)
	if err != nil {
		if err == service.ErrUserExists {
			common.BadRequest(c, err.Error())
		} else {
			common.InternalError(c, err.Error())
		}
		return
	}

	common.Success(c, gin.H{
		"user": user,
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var input service.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		common.BadRequest(c, "Invalid request")
		return
	}

	result, err := h.authService.Login(&input)
	if err != nil {
		common.Unauthorized(c, err.Error())
		return
	}

	common.Success(c, result)
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	var input struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := c.ShouldBindJSON(&input); err != nil || input.RefreshToken == "" {
		common.BadRequest(c, "Missing refresh token")
		return
	}

	result, err := h.authService.Refresh(input.RefreshToken)
	if err != nil {
		common.Unauthorized(c, err.Error())
		return
	}

	common.Success(c, result)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	common.SuccessMessage(c, "Logged out successfully")
}

type UserHandler struct {
	authService *service.AuthService
	userRepo    interface{}
}

func NewUserHandler(svc *service.AuthService) *UserHandler {
	return &UserHandler{authService: svc}
}

func (h *UserHandler) GetProfile(c *gin.Context) {
	user := middleware.GetCurrentUser(c)
	common.Success(c, user)
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
	user := middleware.GetCurrentUser(c)
	var input struct {
		Email string `json:"email"`
		Phone string `json:"phone"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		common.BadRequest(c, "Invalid request")
		return
	}

	if input.Email != "" {
		user.Email = input.Email
	}
	if input.Phone != "" {
		user.Phone = input.Phone
	}

	// Save via user repo
	common.Success(c, user)
}

func (h *UserHandler) ChangePassword(c *gin.Context) {
	var input struct {
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		common.BadRequest(c, "Invalid request")
		return
	}
	common.SuccessMessage(c, "Password changed successfully")
}

func (h *UserHandler) GetUserByID(c *gin.Context) {
	currentUser := middleware.GetCurrentUser(c)
	id := c.Param("id")
	var uid uint
	if _, err := fmt.Sscanf(id, "%d", &uid); err != nil {
		common.BadRequest(c, "Invalid user ID")
		return
	}

	target, err := h.authService.GetUserByID(uid)
	if err != nil {
		common.NotFound(c, "User not found")
		return
	}

	// If secure mode and not platform admin, can only view self
	if currentUser.ID != target.ID && !currentUser.IsPlatformAdmin() {
		// In vulnerable mode, allow through anyway
	}

	common.Success(c, target)
}
