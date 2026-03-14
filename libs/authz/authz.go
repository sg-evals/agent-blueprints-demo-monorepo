package authz

import (
	"context"
	"fmt"
	"net/http"
	"strings"
)

type contextKey string

const userContextKey contextKey = "user"

// Role represents a user role.
type Role string

const (
	RoleAdmin  Role = "admin"
	RoleUser   Role = "user"
	RoleViewer Role = "viewer"
)

// User represents an authenticated user.
type User struct {
	ID    string
	Email string
	Roles []Role
}

// HasRole checks if the user has the given role.
func (u *User) HasRole(role Role) bool {
	for _, r := range u.Roles {
		if r == role {
			return true
		}
	}
	return false
}

// WithUser adds a user to the context.
func WithUser(ctx context.Context, u *User) context.Context {
	return context.WithValue(ctx, userContextKey, u)
}

// UserFromContext extracts the user from context.
func UserFromContext(ctx context.Context) (*User, bool) {
	u, ok := ctx.Value(userContextKey).(*User)
	return u, ok
}

// RequireRole returns middleware that requires the given role.
func RequireRole(role Role, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u, ok := UserFromContext(r.Context())
		if !ok {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		if !u.HasRole(role) {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// ParseAuthHeader extracts a bearer token from the Authorization header.
func ParseAuthHeader(header string) (string, error) {
	parts := strings.SplitN(header, " ", 2)
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return "", fmt.Errorf("invalid authorization header")
	}
	return parts[1], nil
}
