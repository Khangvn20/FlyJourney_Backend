package repository

import "github.com/Khangvn20/FlyJourney_Backend/internal/core/dto"

type UserRepository interface {
	FindByEmail(email string) (*dto.User, error)
	UpdateLastLogin(userID int) error
	Create(user *dto.User) (*dto.User, error)
	UpdatePassword(userID int, newPassword string) error
	GetUserByID(userID int) (*dto.User, error)
	UpdateProfile(userID int, user *dto.User) (*dto.User, error)
	FindByPhone(phone string) (*dto.User, error)
}
