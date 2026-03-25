package middleware

import (
	"ai-support-bot-service/config"
	"ai-support-bot-service/utils"
	"context"
	"net/http"
	"os"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
)

func UserAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("access_token")
		mySecret := os.Getenv("JWT_SECRET")
		if err != nil {
			// Handle no cookie error (e.g., return unauthorized)
			utils.HandleErrorResponse("No access token found", http.StatusUnauthorized, w)
			return
		}

		token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
			// Validate signing method and return the secret key
			return []byte(mySecret), nil
		})
		if err != nil || !token.Valid {
			// Handle invalid token error
			utils.HandleErrorResponse("Invalid access token", http.StatusUnauthorized, w)
			return
		}

		// If the token is valid, you can pass the user data to the next handler
		next.ServeHTTP(w, r)
	})
}

func APIKeyAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("X-API-Key")
		if apiKey == "" {
			utils.HandleErrorResponse("No API key found", http.StatusUnauthorized, w)
			return
		}
		// Validate API key
		exists, err := config.APIKeyCollection.CountDocuments(context.Background(), bson.D{{Key: "key", Value: apiKey}})
		if err != nil {
			utils.HandleErrorResponse("Invalid API key", http.StatusUnauthorized, w)
			return
		}
		if exists == 0 {
			utils.HandleErrorResponse("Invalid API key", http.StatusUnauthorized, w)
			return
		}
		next.ServeHTTP(w, r)
	})
}
