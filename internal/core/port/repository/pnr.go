package repository

import "github.com/Khangvn20/FlyJourney_Backend/internal/core/dto"

type PnrRepository interface {
    CreatePnr(pnr *dto.PNR) (*dto.PNR, error)
    CheckPnrExists(code string) (bool, error)
    GetPnrByBookingID(bookingID int64) (*dto.PNR, error)
    GetBookingIDByPnrCode(pnrCode string) (int64, error)
    UpdatePnr(pnr *dto.PNR) (*dto.PNR, error)
}