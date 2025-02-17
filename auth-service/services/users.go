package services

import (
	"errors"
	"time"
)

var (
	ErrUserNotFound            = errors.New("user not found")
	ErrUserEmailAlreadyInUse   = errors.New("email already in use")
	ErrInvalidEmail            = errors.New("user's email is invalid")
	ErrInsecurePassword        = errors.New("user's password is insecure")
	ErrUnexpectedError         = errors.New("unexpected error")
	ErrCanNotChangeToSameEmail = errors.New("can not change user's email to same email")
)

type UserModel struct {
	UserID    string    `json:"user_id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserCreatedModel struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserService interface {
	CreateNewUser(new_user UserCreatedModel) (string, error)
	FromID(user_id string) (*UserModel, error)
	UpdateEmailFromUserID(user_id string, new_email string) error
	UpdatePasswordFromUserID(user_id string, new_passworde string) error
	DeleteFromUserID(user_id string) error
}
