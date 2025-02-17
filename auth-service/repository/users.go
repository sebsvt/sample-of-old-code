package repository

import (
	"time"
)

type User struct {
	UserID             string    `db:"user_id"`
	Email              string    `db:"email"`
	HashedPassword     string    `db:"hashed_password"`
	Salt               string    `db:"salt"`
	RefreshToken       string    `db:"refresh_token"`
	RefreshTokenExpiry time.Time `db:"refresh_token_expiry"`
	CreatedAt          time.Time `db:"created_at"`
	UpdatedAt          time.Time `db:"updated_at"`
}

type UserRepository interface {
	Create(entity User) error
	FromEmail(email string) (*User, error)
	FromUserID(user_id string) (*User, error)
	Update(entity User) error
	UpdateRefreshToken(user_id string, refresh_token string, expriy time.Time) error
	DeleteByUserID(user_id string) error
}
