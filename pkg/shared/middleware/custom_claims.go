package middleware

import (
	"context"
	"errors"
)

// CustomClaims contains custom data we want from the token.
type CustomClaims struct {
	Permissions  []string `json:"permissions,omitempty"`
	ShouldReject bool     `json:"shouldReject,omitempty"`
}

// Validate errors out if `ShouldReject` is true.
func (c *CustomClaims) Validate(ctx context.Context) error {
	if c.ShouldReject {
		return errors.New("should reject was set to true")
	}
	return nil
}
