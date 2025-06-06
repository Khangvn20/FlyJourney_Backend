package service

import (
	"github.com/Khangvn20/FlyJourney_Backend/internal/core/dto"
	"github.com/Khangvn20/FlyJourney_Backend/internal/core/entity/error_code"
	"github.com/Khangvn20/FlyJourney_Backend/internal/core/model/request"
	"github.com/Khangvn20/FlyJourney_Backend/internal/core/model/response"
	"github.com/Khangvn20/FlyJourney_Backend/internal/core/port/repository"
	"github.com/Khangvn20/FlyJourney_Backend/internal/core/port/service"
	"golang.org/x/crypto/bcrypt"
	"regexp"
	"time"
)

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) service.UserService {
	return &userService{
		userRepo: userRepo,
	}
}

// validateName returns false if the name is only digits or empty, else true
func validateName(name string) bool {
	re := regexp.MustCompile(`^\d+$`)
	return name != "" && !re.MatchString(name)
}

// validateEmail returns true if email is valid
func validateEmail(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}

// validatePhone returns true if phone is valid (optional field)
func validatePhone(phone string) bool {
	if phone == "" {
		return true
	}
	re := regexp.MustCompile(`^\d{10,15}$`)
	return re.MatchString(phone)
}

// handleValidationErrors helps to reduce code duplication for responses
func handleValidationErrors(message string) *response.Response {
	return &response.Response{
		Status:       false,
		ErrorCode:    error_code.InvalidRequest,
		ErrorMessage: message,
	}
}

func (s *userService) Register(req *request.RegisterRequest) *response.Response {
	// Validate email
	if !validateEmail(req.Email) {
		return handleValidationErrors("Email is invalid")
	}
	// Validate password
	if len(req.Password) < 6 {
		return handleValidationErrors("Password must be at least 6 characters")
	}
	// Validate name
	if !validateName(req.Name) {
		return handleValidationErrors("Name is required and must not be a number")
	}
	// Validate phone if provided
	if !validatePhone(req.Phone) {
		return handleValidationErrors("Phone is invalid")
	}

	// Check if email exists
	userExist, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		return &response.Response{
			Status:       false,
			ErrorCode:    error_code.InternalError,
			ErrorMessage: error_code.InternalErrMsg,
		}
	}
	if userExist != nil {
		return &response.Response{
			Status:       false,
			ErrorCode:    error_code.DuplicateUser,
			ErrorMessage: "Email đã tồn tại",
		}
	}

	// Hash password
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return &response.Response{
			Status:       false,
			ErrorCode:    error_code.InternalError,
			ErrorMessage: error_code.InternalErrMsg,
		}
	}

	now := time.Now()
	user := &dto.User{
		Email:     req.Email,
		Password:  string(hashed),
		Name:      req.Name,
		Phone:     req.Phone,
		Role:      "user",
		CreatedAt: now,
		UpdatedAt: now,
	}

	newUser, err := s.userRepo.Create(user)
	if err != nil {
		return &response.Response{
			Status:       false,
			ErrorCode:    error_code.InternalError,
			ErrorMessage: error_code.InternalErrMsg,
		}
	}

	data := response.RegisterResponse{
		UserID: newUser.UserID,
		Email:  newUser.Email,
		Name:   newUser.Name,
	}
	return &response.Response{
		Data:         data,
		Status:       true,
		ErrorCode:    error_code.Success,
		ErrorMessage: error_code.SuccessErrMsg,
	}
}

func (s *userService) Login(req *request.LoginRequest) *response.Response {
	// Validate email format
	if !validateEmail(req.Email) {
		return handleValidationErrors("Email is invalid")
	}
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		return &response.Response{
			Status:       false,
			ErrorCode:    error_code.InternalError,
			ErrorMessage: error_code.InternalErrMsg,
		}
	}

	if user == nil {
		return &response.Response{
			Status:       false,
			ErrorCode:    error_code.InvalidRequest,
			ErrorMessage: "Email hoặc mật khẩu không đúng",
		}
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return &response.Response{
			Status:       false,
			ErrorCode:    error_code.InvalidRequest,
			ErrorMessage: "Email hoặc mật khẩu không đúng",
		}
	}

	// Tạo JWT token
	tokenService := NewTokenService()
	token, expireAt, err := tokenService.GenerateToken(user)
	if err != nil {
		return &response.Response{
			Status:       false,
			ErrorCode:    error_code.InternalError,
			ErrorMessage: error_code.InternalErrMsg,
		}
	}
	_ = s.userRepo.UpdateLastLogin(user.UserID)

	// Tạo response
	loginResponse := &response.LoginResponse{
		UserID:   user.UserID,
		Email:    user.Email,
		Name:     user.Name,
		Role:     user.Role,
		Token:    token,
		ExpireAt: expireAt,
	}

	return &response.Response{
		Data:         loginResponse,
		Status:       true,
		ErrorCode:    error_code.Success,
		ErrorMessage: error_code.SuccessErrMsg,
	}
}
