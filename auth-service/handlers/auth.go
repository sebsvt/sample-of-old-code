package handlers

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/sebsvt/cmu-contest-2024/auth-service/services"
)

var (
	ErrInvalidRequestBody = errors.New("invalid request body")
)

type UserCrendentail struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RefreshTokenBody struct {
	RefreshToken string `json:"refresh_token"`
}

type authHandler struct {
	user_srv services.UserService
	auth_srv services.Authorization
}

func NewAuthHandler(user_srv services.UserService, auth_srv services.Authorization) authHandler {
	return authHandler{
		user_srv: user_srv,
		auth_srv: auth_srv,
	}
}

func (h authHandler) SignUp(c *fiber.Ctx) error {
	var new_user services.UserCreatedModel
	if err := c.BodyParser(&new_user); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": ErrInvalidRequestBody.Error(),
		})
	}

	user_id, err := h.user_srv.CreateNewUser(new_user)
	if err != nil {
		switch err {
		case services.ErrUserEmailAlreadyInUse,
			services.ErrInsecurePassword,
			services.ErrInvalidEmail:
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		default:
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
	}

	return c.JSON(fiber.Map{
		"user_id": user_id,
	})
}

func (h authHandler) SignIn(c *fiber.Ctx) error {
	var user_crendentail UserCrendentail
	if err := c.BodyParser(&user_crendentail); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": ErrInvalidRequestBody.Error(),
		})
	}
	access_token, refresh_token, err := h.auth_srv.SignIn(user_crendentail.Email, user_crendentail.Password)
	if err != nil {
		switch err {
		case services.ErrAuthenticationFailed, services.ErrInvalidEmail:
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})

		default:
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
	}
	return c.JSON(fiber.Map{
		"access_token":  access_token,
		"refresh_token": refresh_token,
		"type":          "Bearer",
	})
}

func (h authHandler) RefreshToken(c *fiber.Ctx) error {
	var refresh_token RefreshTokenBody
	if err := c.BodyParser(&refresh_token); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": ErrInvalidRequestBody.Error(),
		})
	}
	new_access_token, new_refresh_token, err := h.auth_srv.Refresh(refresh_token.RefreshToken)
	if err != nil {
		switch err {
		case services.ErrUserNotFound, services.ErrBadClaim, services.ErrInvalidSignature, services.ErrTokenExpired:
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		default:
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"error": services.ErrUnexpectedError,
			})
		}
	}
	return c.JSON(fiber.Map{
		"access_token":  new_access_token,
		"refresh_token": new_refresh_token,
		"type":          "Bearer",
	})
}

func (h authHandler) Authorize(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")

	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "missing token",
		})
	}

	const bearerPrefix = "Bearer "
	if !strings.HasPrefix(authHeader, bearerPrefix) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "invalid token format",
		})
	}

	// Extract the token without the Bearer prefix
	tokenString := authHeader[len(bearerPrefix):]

	// Use the authorization service to validate the token
	info, err := h.auth_srv.Authorize(tokenString)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "invalid or expired token",
		})
	}

	if info.Type != services.AccessToken {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "invalid token",
		})
	}
	user, err := h.user_srv.FromID(info.Subject)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "user not found",
		})
	}
	return c.JSON(user)
}
