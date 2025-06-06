package repository

import (
	"context"
	"errors"
	"github.com/Khangvn20/FlyJourney_Backend/internal/core/dto"
	coreRepo "github.com/Khangvn20/FlyJourney_Backend/internal/core/port/repository"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"time"
)

type userRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db coreRepo.Database) coreRepo.UserRepository {
	return &userRepository{db: db.GetPool()}
}

func (r *userRepository) FindByEmail(email string) (*dto.User, error) {

	query := `
		SELECT user_id, email, password, name, phone, role, created_at, updated_at, last_login 
		FROM users 
		WHERE email = $1 
		LIMIT 1
	`

	var user dto.User
	var lastLogin *time.Time

	// Thêm log để debug
	log.Printf("Executing FindByEmail query for email: %s", email)

	err := r.db.QueryRow(context.Background(), query, email).Scan(
		&user.UserID,
		&user.Email,
		&user.Password,
		&user.Name,
		&user.Phone,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
		&lastLogin,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		log.Printf("No user found with email: %s", email)
		return nil, nil
	}

	if err != nil {
		log.Printf("Error in FindByEmail: %v", err)
		return nil, err
	}

	user.LastLogin = lastLogin
	log.Printf("User found with email: %s", email)
	return &user, nil
}

func (r *userRepository) Create(user *dto.User) (*dto.User, error) {
	// Sửa lại RETURNING để lấy chính xác user_id
	query := `
		INSERT INTO users (email, password, name, phone, role, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, $5, $6, $7) 
		RETURNING user_id
	`

	// Thêm log để debug
	log.Printf("Creating user with email: %s", user.Email)

	err := r.db.QueryRow(context.Background(), query,
		user.Email,
		user.Password,
		user.Name,
		user.Phone,
		user.Role,
		user.CreatedAt,
		user.UpdatedAt,
	).Scan(&user.UserID)

	if err != nil {
		log.Printf("Error creating user: %v", err)
		return nil, err
	}

	log.Printf("User created successfully with ID: %d", user.UserID)
	return user, nil
}

func (r *userRepository) UpdateLastLogin(userID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `UPDATE users SET last_login = $1 WHERE user_id = $2`
	_, err := r.db.Exec(ctx, query, time.Now(), userID)
	if err != nil {
		log.Printf("Error updating last login time: %v", err)
		return err
	}
	return nil
}
