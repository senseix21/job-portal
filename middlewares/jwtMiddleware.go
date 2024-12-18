package middlewares

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

// JWTSecret is the secret key used for signing the JWT. Replace it with your own secret.
var JWTSecret = []byte("secret_key")

// Claims represents the custom claims in the JWT
type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// JWTMiddleware authenticates a JWT token and checks if the user's role matches any of the allowed roles.
func JWTMiddleware(allowedRoles ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Get the token from the Authorization header
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "Authorization header is required")
			}

			// Validate the Authorization header format
			// if !strings.HasPrefix(authHeader, "Bearer ") {
			// 	return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token format, must start with 'Bearer '")
			// }

			tokenString := authHeader

			// Parse and validate the token
			claims := &Claims{}
			token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
				// Ensure the signing method is HMAC
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return JWTSecret, nil
			})

			if err != nil {
				// Handle token parsing errors
				if errors.Is(err, jwt.ErrTokenExpired) {
					return echo.NewHTTPError(http.StatusUnauthorized, "Token has expired")
				}
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token")
			}

			// Ensure the token is valid
			if !token.Valid {
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token")
			}

			// Check if the user's role matches any of the allowed roles
			roleAllowed := false
			for _, role := range allowedRoles {
				if claims.Role == role {
					roleAllowed = true
					break
				}
			}

			if !roleAllowed {
				return echo.NewHTTPError(http.StatusForbidden, "Access denied for this role")
			}

			// Set user information in the context
			c.Set("userID", claims.UserID)
			c.Set("email", claims.Email)
			c.Set("role", claims.Role)

			// Proceed to the next handler
			return next(c)
		}
	}
}
