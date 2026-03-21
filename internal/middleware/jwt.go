package middleware

import (
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const ctxUserIDKey = "user_id"

func JWT(secret []byte) fiber.Handler {
	return func(c *fiber.Ctx) error {
		h := c.Get("Authorization")
		if h == "" || !strings.HasPrefix(strings.ToLower(h), "bearer ") {
			return fiber.ErrUnauthorized
		}
		raw := strings.TrimSpace(h[7:])
		tok, err := jwt.Parse(raw, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fiber.ErrUnauthorized
			}
			return secret, nil
		})
		if err != nil || !tok.Valid {
			return fiber.ErrUnauthorized
		}
		claims, ok := tok.Claims.(jwt.MapClaims)
		if !ok {
			return fiber.ErrUnauthorized
		}
		sub, _ := claims["sub"].(string)
		uid, err := uuid.Parse(sub)
		if err != nil {
			return fiber.ErrUnauthorized
		}
		c.Locals(ctxUserIDKey, uid)
		return c.Next()
	}
}

func UserID(c *fiber.Ctx) (uuid.UUID, bool) {
	v := c.Locals(ctxUserIDKey)
	if v == nil {
		return uuid.Nil, false
	}
	u, ok := v.(uuid.UUID)
	return u, ok
}

func SignJWT(secret []byte, userID uuid.UUID, ttl time.Duration) (string, error) {
	now := time.Now()
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID.String(),
		"iat": now.Unix(),
		"exp": now.Add(ttl).Unix(),
	})
	return t.SignedString(secret)
}
