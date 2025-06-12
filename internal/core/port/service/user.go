package service

import (
	"github.com/Khangvn20/FlyJourney_Backend/internal/core/model/request"
	"github.com/Khangvn20/FlyJourney_Backend/internal/core/model/response"
)

type UserService interface {
	Login(request *request.LoginRequest) *response.Response
	Register(request *request.RegisterRequest) *response.Response
	ConfirmRegister(request *request.ConfirmRegisterRequest) *response.Response
	ConfirmResetPassword(request *request.ConfirmResetPasswordRequest) *response.Response
	ResetPassword(request *request.ResetPasswordRequest) *response.Response
	Logout(token string) *response.Response
	GetUserInfo(userID int) *response.Response
}
