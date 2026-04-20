package middlewares

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

type ContextKey string

const (
	RoleKey     ContextKey = "role"
	UsernameKey ContextKey = "username"
	IDKey       ContextKey = "id"
)

func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 1. Extract the token from the "Bearer" cookie
		cookie, err := r.Cookie("Bearer")
		if err != nil {
			// If there's no cookie, the user is not authenticated
			http.Error(w, "Unauthorized: Missing Token", http.StatusUnauthorized)
			return
		}

		tokenString := cookie.Value
		jwtSecret := os.Getenv("JWT_SECRET")

		// 2. Parse and validate the JWT
		parsedToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Verify that the signing method used is what we expect (HMAC for HS256)
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			// Provide the secret key to the parser so it can verify the signature
			return []byte(jwtSecret), nil
		})

		// 3. Check if there was an error during parsing, or if the token is invalid/expired
		if err != nil || !parsedToken.Valid {
			http.Error(w, "Unauthorized: Invalid or Expired Token", http.StatusUnauthorized)
			return
		}

		claims, ok := parsedToken.Claims.(jwt.MapClaims)

		if ok {
			fmt.Println(claims["uid"], claims["exp"], claims["role"])
		} else {
			http.Error(w, "Invalid login Token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), RoleKey, claims["role"])
		ctx = context.WithValue(ctx, UsernameKey, claims["username"])
		ctx = context.WithValue(ctx, IDKey, claims["uid"])
		// 4. If everything is valid, proceed to the actual API handler
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
