package service

import (
	"github.com/Khangvn20/FlyJourney_Backend/internal/core/model/request"
	"github.com/Khangvn20/FlyJourney_Backend/internal/core/model/response"
)

type UserService interface {
	Login(request *request.LoginRequest) *response.Response
	Register(request *request.RegisterRequest) *response.Response
}
