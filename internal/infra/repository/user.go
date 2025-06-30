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
	"fmt"
)

type userRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db coreRepo.Database) coreRepo.UserRepository {
	return &userRepository{db: db.GetPool()}
}

func (r *userRepository) FindByEmail(email string) (*dto.User, error) {
      query := `
        SELECT 
            user_id, 
            email, 
            password, 
            name, 
            roles,  
            created_at, 
            updated_at, 
            last_login 
        FROM users 
        WHERE email = $1 
        LIMIT 1
    `

    var user dto.User
    var lastLogin *time.Time
    var roles string 

    err := r.db.QueryRow(context.Background(), query, email).Scan(
        &user.UserID,
        &user.Email,
        &user.Password,
        &user.Name,
        &roles,   
        &user.CreatedAt,
        &user.UpdatedAt,
        &lastLogin,
    )

    if err == pgx.ErrNoRows {
        return nil, nil
    }
    if err != nil {
        log.Printf("Error scanning user data: %v", err)
        return nil, fmt.Errorf("error finding user: %w", err)
    }

    user.Roles = dto.UserRole(roles)
    user.LastLogin = lastLogin

    return &user, nil
}

func (r *userRepository) Create(user *dto.User) (*dto.User, error) {
    if user.Roles == "" {
        user.Roles = dto.RoleUser
    }
      query := `
        INSERT INTO users (email, password, name, roles, created_at, updated_at) 
        VALUES ($1, $2, $3, ARRAY[$4]::user_role[], $5, $6) 
        RETURNING user_id
    `


    log.Printf("Creating user with email: %s and role: %s", user.Email, user.Roles)

    err := r.db.QueryRow(context.Background(), query,
        user.Email,
        user.Password,
        user.Name,
        string(user.Roles), 
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
func (r *userRepository) UpdatePassword(userID int, newPassword string) error {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    query := `UPDATE users SET password = $1, updated_at = $2 WHERE user_id = $3`
    _, err := r.db.Exec(ctx, query, newPassword, time.Now(), userID)
    if err != nil {
        return err
    }
    return nil
}
func (r *userRepository) GetUserByID(userID int) (*dto.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		SELECT user_id, email, name, phone, roles, created_at, updated_at, last_login 
		FROM users 
		WHERE user_id = $1 
		LIMIT 1
	`

	var user dto.User
	var lastLogin *time.Time
    var roles string 

	err := r.db.QueryRow(ctx, query, userID).Scan(
		&user.UserID,
		&user.Email,
		&user.Name,
		&user.Phone,
		&roles,  
		&user.CreatedAt,
		&user.UpdatedAt,
		&lastLogin,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		log.Printf("No user found with ID: %d", userID)
		return nil, nil
	}

	if err != nil {
		log.Printf("Error in GetUserByID: %v", err)
		return nil, err
	}

	user.LastLogin = lastLogin
	return &user, nil
}
func (r *userRepository) UpdateProfile(userID int, user *dto.User) (*dto.User, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    // Add logging
    log.Printf("Executing DB query to update profile for user ID: %d with name: %s, phone: %s", 
        userID, user.Name, user.Phone)
    
    query := `UPDATE users SET name = $1, phone = $2, updated_at = $3 WHERE user_id = $4`
    result, err := r.db.Exec(ctx, query, user.Name, user.Phone, time.Now(), userID)
    if err != nil {
        log.Printf("Database error updating user profile: %v", err)
        return nil, err
    }
   rowsAffected:= result.RowsAffected()
    if rowsAffected == 0 {
        log.Printf("No rows were updated for user ID: %d", userID)
        return nil, fmt.Errorf("no rows updated")
    }   
    log.Printf("Successfully updated %d rows for user ID: %d", rowsAffected, userID)
    return r.GetUserByID(userID)
}