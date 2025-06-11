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
	emailOTPService service.EmailOTPService
	tokenService service.TokenService
}

func NewUserService(userRepo repository.UserRepository , emailOTPService service.EmailOTPService, tokenService service.TokenService ) service.UserService {
	return &userService{
		userRepo: userRepo,
		emailOTPService:  emailOTPService,
		tokenService: tokenService,
	
	}
}

func validateName(name string) bool {
	re := regexp.MustCompile(`^\d+$`)
	return name != "" && !re.MatchString(name)
}

func validateEmail(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}

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
     otpResult := s.emailOTPService.SendOTPEmail(req.Email)
    if !otpResult.Status {
        return otpResult
    }

    return &response.Response{
        Status:       true,
        ErrorCode:    error_code.Success,
        ErrorMessage: "Đã gửi OTP xác thực email, vui lòng kiểm tra email.",
    }
	
}
func (s *userService) Login(req *request.LoginRequest) *response.Response {
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
            ErrorCode:    error_code.InternalError,
            ErrorMessage: "Email hoặc mật khẩu không đúng",
        }
    }
	  err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
    if err != nil {
        return &response.Response{
            Status:       false,
            ErrorCode:    error_code.InternalError,
            ErrorMessage: "Email hoặc mật khẩu không đúng",
        }
    }

    tokenDuration := time.Hour * 24
    token, err := s.tokenService.GenerateToken(user.UserID, tokenDuration)
    if err != nil {
        return &response.Response{
            Status:       false,
            ErrorCode:    error_code.InternalError,
            ErrorMessage: error_code.InternalErrMsg,
        }
    }

    // 5. Trả về response
    return &response.Response{
        Status:    true,
        ErrorCode: error_code.Success,
        Data: map[string]interface{}{
            "user_id":   user.UserID,
            "email":     user.Email,
            "name":      user.Name,
            "role":      user.Role,
            "token":     token,
           
        },
        ErrorMessage: error_code.SuccessErrMsg,
    }
}

func (s *userService) ConfirmRegister(req *request.ConfirmRegisterRequest) *response.Response {
    otpResult := s.emailOTPService.VerifyEmail(req.Email, req.OTP)
    if !otpResult.Status {
        return otpResult
    }
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

    // 3. Hash password
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

    return &response.Response{
        Data: map[string]interface{}{
            "user_id": newUser.UserID,
            "email":   newUser.Email,
            "name":    newUser.Name,
        },
        Status:       true,
        ErrorCode:    error_code.Success,
        ErrorMessage: error_code.SuccessErrMsg,
    }
}
func (s *userService) ResetPassword(req *request.ResetPasswordRequest) *response.Response {
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
            ErrorCode:    error_code.InternalError,
            ErrorMessage: "Email không tồn tại",
        }
    }
    // 3. Gửi OTP đến email
    otpResult := s.emailOTPService.SendOTPEmail(req.Email)
    if !otpResult.Status {
        return otpResult
    }
    return &response.Response{
        Status:       true,
        ErrorCode:    error_code.Success,
        ErrorMessage: "Đã gửi OTP xác thực đến email, vui lòng kiểm tra email.",
    }
}


func (s *userService) ConfirmResetPassword(req *request.ConfirmResetPasswordRequest) *response.Response {
    otpResult := s.emailOTPService.VerifyEmail(req.Email, req.OTP)
    if !otpResult.Status {
        return otpResult
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
            ErrorCode:    error_code.InternalError,
            ErrorMessage: "Email không tồn tại",
        }
    }

    hashed, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
    if err != nil {
        return &response.Response{
            Status:       false,
            ErrorCode:    error_code.InternalError,
            ErrorMessage: error_code.InternalErrMsg,
        }
    }
    user.Password = string(hashed)
    user.UpdatedAt = time.Now()
    err = s.userRepo.UpdatePassword(user.UserID, string(hashed))
    if err != nil {
        return &response.Response{
            Status:       false,
            ErrorCode:    error_code.InternalError,
            ErrorMessage: error_code.InternalErrMsg,
        }
    }
    return &response.Response{
        Status:       true,
        ErrorCode:    error_code.Success,
        ErrorMessage: "Đặt lại mật khẩu thành công",
    }
}