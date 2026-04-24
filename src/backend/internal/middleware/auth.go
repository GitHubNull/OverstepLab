package middleware

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/oversteplab/oversteplab/internal/common"
	"github.com/oversteplab/oversteplab/internal/model"
	"github.com/oversteplab/oversteplab/internal/repository"
	"github.com/oversteplab/oversteplab/internal/vuln"
)

const (
	CtxCurrentUser = "currentUser"
)

type JWTClaims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	UserType string `json:"user_type"`
	Role     string `json:"role"`
	CompanyID *uint `json:"company_id,omitempty"`
	jwt.RegisteredClaims
}

func AuthMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			common.Unauthorized(c, "Missing authorization header")
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			common.Unauthorized(c, "Invalid authorization format")
			c.Abort()
			return
		}

		tokenString := parts[1]

		// Check if using API Key authentication
		if strings.HasPrefix(tokenString, "sk_") {
			if handleAPIKeyAuth(c, tokenString, jwtSecret) {
				return
			}
			common.Unauthorized(c, "Invalid API key")
			c.Abort()
			return
		}

		claims, err := parseJWT(tokenString, jwtSecret)
		if err != nil {
			common.Unauthorized(c, "Invalid or expired token")
			c.Abort()
			return
		}

		userRepo := repository.GetUserRepo()
		user, err := userRepo.FindByID(claims.UserID)
		if err != nil {
			common.Unauthorized(c, "User not found")
			c.Abort()
			return
		}

		if user.Status != "active" {
			common.Forbidden(c, "Account is disabled")
			c.Abort()
			return
		}

		c.Set(CtxCurrentUser, user)
		c.Next()
	}
}

func handleAPIKeyAuth(c *gin.Context, tokenString, jwtSecret string) bool {
	apiKeyRepo := repository.GetAPIKeyRepo()
	userRepo := repository.GetUserRepo()

	apiKey, err := apiKeyRepo.FindByValue(tokenString)
	if err != nil {
		return false
	}

	// C-03 Vulnerability: In vulnerable mode, skip status and expiration checks
	if !vuln.IsSecureMode() {
		// Vulnerable mode: only check if key exists
	} else {
		// Secure mode: enforce status and expiration
		if apiKey.Status != "active" {
			return false
		}
		if apiKey.ExpireAt != nil && apiKey.ExpireAt.Before(time.Now()) {
			return false
		}
	}

	user, err := userRepo.FindByID(apiKey.UserID)
	if err != nil {
		return false
	}

	if user.Status != "active" {
		return false
	}

	// Update last used
	now := time.Now()
	apiKey.LastUsedAt = &now
	apiKeyRepo.Update(apiKey)

	c.Set(CtxCurrentUser, user)
	return true
}

func ParseJWTForRefresh(tokenString, secret string) (*JWTClaims, error) {
	return parseJWT(tokenString, secret)
}

func parseJWT(tokenString, secret string) (*JWTClaims, error) {
	claims := &JWTClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil || !token.Valid {
		return nil, err
	}
	return claims, nil
}

func GenerateJWT(user *model.User, secret string, expiry time.Duration) (string, error) {
	claims := &JWTClaims{
		UserID:    user.ID,
		Username:  user.Username,
		UserType:  user.UserType,
		Role:      user.Role,
		CompanyID: user.CompanyID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func GetCurrentUser(c *gin.Context) *model.User {
	val, exists := c.Get(CtxCurrentUser)
	if !exists {
		return nil
	}
	user, ok := val.(*model.User)
	if !ok {
		return nil
	}
	return user
}

func OptionalAuth(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.Next()
			return
		}

		claims, err := parseJWT(parts[1], jwtSecret)
		if err != nil {
			c.Next()
			return
		}

		userRepo := repository.GetUserRepo()
		user, err := userRepo.FindByID(claims.UserID)
		if err != nil {
			c.Next()
			return
		}

		c.Set(CtxCurrentUser, user)
		c.Next()
	}
}
