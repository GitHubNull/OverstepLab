package middleware

import (
	"crypto/sha256"
	"encoding/hex"
	"strconv"
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
	Type     string `json:"type"`
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
				c.Next()
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

		if claims.Type != "" && claims.Type != "access" {
			common.Unauthorized(c, "Invalid token type")
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

	// Hash the raw API key before querying (database stores hashed values)
	h := sha256.Sum256([]byte(tokenString))
	hashedKey := hex.EncodeToString(h[:])

	apiKey, err := apiKeyRepo.FindByValue(hashedKey)
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

func GenerateJWT(user *model.User, secret string, expiry time.Duration, tokenType string) (string, error) {
	claims := &JWTClaims{
		UserID:    user.ID,
		Username:  user.Username,
		UserType:  user.UserType,
		Role:      user.Role,
		CompanyID: user.CompanyID,
		Type:      tokenType,
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

// ---- Encoding Challenge Helpers ----
// NOTE: Since EncodingMiddleware now handles all decoding transparently,
// these helpers simply return raw values from the request. Decoding is
// already done by the middleware before the handler runs.

// GetEncodingType extracts the encoding type from the X-Encoding-Type header.
func GetEncodingType(c *gin.Context) string {
	encType := c.GetHeader("X-Encoding-Type")
	if encType == "" {
		return "none"
	}
	return encType
}

// DecodeQueryParam returns the raw query parameter value.
// Decoding is handled by EncodingMiddleware, so this just reads from the request.
func DecodeQueryParam(c *gin.Context, key string) string {
	return c.Query(key)
}

// DecodeBodyField is a no-op in the global encoding architecture.
// Body fields are already decoded by EncodingMiddleware.
func DecodeBodyField(data map[string]interface{}, key string, encType string) {
	// No-op: EncodingMiddleware handles all decoding
}

// DecodeUintParam returns the query parameter parsed as uint.
// Decoding is handled by EncodingMiddleware.
func DecodeUintParam(c *gin.Context, key string) uint {
	raw := c.Query(key)
	if raw == "" {
		return 0
	}
	id, _ := strconv.ParseUint(raw, 10, 64)
	return uint(id)
}

// DecodeUintBodyField returns a body field parsed as uint.
// Decoding is handled by EncodingMiddleware.
func DecodeUintBodyField(data map[string]interface{}, key string, encType string) uint {
	if val, ok := data[key]; ok {
		switch v := val.(type) {
		case float64:
			return uint(v)
		case string:
			id, _ := strconv.ParseUint(v, 10, 64)
			return uint(id)
		case int:
			return uint(v)
		case uint:
			return v
		}
	}
	return 0
}
