package service

import (
	"errors"
	"time"

	"github.com/oversteplab/oversteplab/internal/config"
	"github.com/oversteplab/oversteplab/internal/middleware"
	"github.com/oversteplab/oversteplab/internal/model"
	"github.com/oversteplab/oversteplab/internal/repository"
	"github.com/oversteplab/oversteplab/internal/vuln"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserExists         = errors.New("username already exists")
	ErrAccountDisabled    = errors.New("account is disabled")
)

type AuthService struct {
	userRepo    *repository.UserRepository
	companyRepo *repository.CompanyRepository
	cfg         *config.Config
}

func NewAuthService(cfg *config.Config) *AuthService {
	return &AuthService{
		userRepo:    repository.GetUserRepo(),
		companyRepo: repository.GetCompanyRepo(),
		cfg:         cfg,
	}
}

type RegisterInput struct {
	Username     string `json:"username"`
	Password     string `json:"password"`
	Email        string `json:"email"`
	Phone        string `json:"phone"`
	UserType     string `json:"user_type"`
	CompanyName  string `json:"company_name"`
}

func (s *AuthService) Register(input *RegisterInput) (*model.User, error) {
	// Check if user exists
	existing, err := s.userRepo.FindByUsername(input.Username)
	if err == nil && existing.ID > 0 {
		return nil, ErrUserExists
	}

	hash, err := hashPassword(input.Password)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Username:     input.Username,
		PasswordHash: hash,
		Email:        input.Email,
		Phone:        input.Phone,
		UserType:     input.UserType,
		Status:       "active",
	}

	// Set default role based on user type
	if input.UserType == "individual" {
		user.Role = "individual"
	}

	// Create company for company-type registration
	if input.UserType == "company" {
		if input.CompanyName == "" {
			return nil, errors.New("company name is required for company registration")
		}
		company := &model.Company{
			Name:   input.CompanyName,
			Status: "active",
		}
		if err := s.companyRepo.Create(company); err != nil {
			return nil, err
		}
		// Capture company ID by value to ensure it's persisted correctly
		companyID := company.ID
		user.CompanyID = &companyID
		user.Role = "admin"
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

type LoginInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginOutput struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
	User         *model.User `json:"user"`
}

func (s *AuthService) Login(input *LoginInput) (*LoginOutput, error) {
	user, err := s.userRepo.FindByUsername(input.Username)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	if user.Status != "active" {
		return nil, ErrAccountDisabled
	}

	if !checkPassword(input.Password, user.PasswordHash) {
		return nil, ErrInvalidCredentials
	}

	token, err := middleware.GenerateJWT(user, s.cfg.JWTSecret, 24*time.Hour)
	if err != nil {
		return nil, err
	}

	refreshToken, err := middleware.GenerateJWT(user, s.cfg.JWTSecret, 7*24*time.Hour)
	if err != nil {
		return nil, err
	}

	return &LoginOutput{
		Token:        token,
		RefreshToken: refreshToken,
		User:         user,
	}, nil
}

func (s *AuthService) Refresh(token string) (*LoginOutput, error) {
	claims, err := middleware.ParseJWTForRefresh(token, s.cfg.JWTSecret)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	user, err := s.userRepo.FindByID(claims.UserID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if user.Status != "active" {
		return nil, ErrAccountDisabled
	}

	newToken, err := middleware.GenerateJWT(user, s.cfg.JWTSecret, 24*time.Hour)
	if err != nil {
		return nil, err
	}

	newRefreshToken, err := middleware.GenerateJWT(user, s.cfg.JWTSecret, 7*24*time.Hour)
	if err != nil {
		return nil, err
	}

	return &LoginOutput{
		Token:        newToken,
		RefreshToken: newRefreshToken,
		User:         user,
	}, nil
}

func (s *AuthService) GetUserByID(id uint) (*model.User, error) {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// H-02 Vulnerability: In vulnerable mode, allow viewing other users' profiles
	if !vuln.IsSecureMode() {
		return user, nil
	}

	// Secure mode: can only view own profile (unless platform admin)
	return user, nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func checkPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
