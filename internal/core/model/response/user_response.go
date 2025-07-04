package response

import "github.com/Khangvn20/FlyJourney_Backend/internal/core/entity/error_code"


type Response struct {
	Data         interface{}          `json:"data"`
	Status       bool                 `json:"status"`
	ErrorCode    error_code.ErrorCode `json:"errorCode"`
	ErrorMessage string               `json:"errorMessage"`
}

