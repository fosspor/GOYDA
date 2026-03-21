package handlers

import (
	"strings"

	"github.com/fosspor/GOYDA/internal/middleware"
	"github.com/fosspor/GOYDA/internal/store"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

type registerBody struct {
	Email       string   `json:"email"`
	Password    string   `json:"password"`
	DisplayName string   `json:"display_name"`
	Interests   []string `json:"interests"`
}

type loginBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (a *API) Register(c *fiber.Ctx) error {
	var b registerBody
	if err := c.BodyParser(&b); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid json")
	}
	b.Email = strings.TrimSpace(b.Email)
	if b.Email == "" || len(b.Password) < 6 {
		return fiber.NewError(fiber.StatusBadRequest, "email and password (min 6) required")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(b.Password), bcrypt.DefaultCost)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "hash failed")
	}
	u, err := a.Store.CreateUser(c.UserContext(), b.Email, string(hash), strings.TrimSpace(b.DisplayName), b.Interests)
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "unique") || strings.Contains(strings.ToLower(err.Error()), "duplicate") {
			return fiber.NewError(fiber.StatusConflict, "email already registered")
		}
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	tok, err := middleware.SignJWT(a.JWTKey, u.ID, a.JWTTTL)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "token error")
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"token": tok,
		"user": fiber.Map{
			"id":           u.ID,
			"email":        u.Email,
			"display_name": u.DisplayName,
			"interests":    u.Interests,
		},
	})
}

func (a *API) Login(c *fiber.Ctx) error {
	var b loginBody
	if err := c.BodyParser(&b); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid json")
	}
	b.Email = strings.TrimSpace(b.Email)
	u, err := a.Store.UserByEmail(c.UserContext(), b.Email)
	if err != nil {
		if err == store.ErrNotFound {
			return fiber.NewError(fiber.StatusUnauthorized, "invalid credentials")
		}
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	if bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(b.Password)) != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "invalid credentials")
	}
	tok, err := middleware.SignJWT(a.JWTKey, u.ID, a.JWTTTL)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "token error")
	}
	return c.JSON(fiber.Map{
		"token": tok,
		"user": fiber.Map{
			"id":           u.ID,
			"email":        u.Email,
			"display_name": u.DisplayName,
			"interests":    u.Interests,
		},
	})
}

func (a *API) Me(c *fiber.Ctx) error {
	uid, ok := middleware.UserID(c)
	if !ok {
		return fiber.ErrUnauthorized
	}
	u, err := a.Store.UserByID(c.UserContext(), uid)
	if err != nil {
		if err == store.ErrNotFound {
			return fiber.ErrUnauthorized
		}
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.JSON(fiber.Map{
		"id":           u.ID,
		"email":        u.Email,
		"display_name": u.DisplayName,
		"interests":    u.Interests,
	})
}

type patchMeBody struct {
	Interests []string `json:"interests"`
}

func (a *API) PatchMe(c *fiber.Ctx) error {
	uid, ok := middleware.UserID(c)
	if !ok {
		return fiber.ErrUnauthorized
	}
	var b patchMeBody
	if err := c.BodyParser(&b); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid json")
	}
	if err := a.Store.UpdateUserInterests(c.UserContext(), uid, b.Interests); err != nil {
		if err == store.ErrNotFound {
			return fiber.ErrUnauthorized
		}
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return a.Me(c)
}
