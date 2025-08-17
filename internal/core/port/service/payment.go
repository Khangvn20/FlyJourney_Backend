package service

import (
    "github.com/Khangvn20/FlyJourney_Backend/internal/core/model/request"
	"github.com/Khangvn20/FlyJourney_Backend/internal/core/model/response"
	"github.com/gin-gonic/gin"
)

type PaymentService interface {
    CreateMomoPayment(req *request.MomoRequest) response.Response
	HandleMomoSuccess(ctx *gin.Context) response.Response
	HandleMomoCallback(req *request.MomoCallbackRequest) response.Response
	
}