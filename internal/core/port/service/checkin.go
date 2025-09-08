package service

import (
	"github.com/Khangvn20/FlyJourney_Backend/internal/core/model/response"
	"github.com/Khangvn20/FlyJourney_Backend/internal/core/model/request"
)

type CheckinService interface {
	ValidateCheckin(req *request.ValidateCheckin) *response.Response
	GetSeatMap(flightID int64) *response.Response
	ProcessOnlineCheckin(req *request.OnlineCheckinRequest) *response.Response
}
