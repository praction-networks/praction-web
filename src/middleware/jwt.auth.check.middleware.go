package middlewares

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/praction-networks/quantum-ISP365/webapp/src/models"
	"github.com/praction-networks/quantum-ISP365/webapp/src/response"
	"github.com/praction-networks/quantum-ISP365/webapp/src/service"
)

// Define a custom type for the context key
type contextKey string

const userContextKey contextKey = "user"

// JWTAuthMiddleware is a middleware that checks the validity of the JWT token
func JWTAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check for both Authorization and Strict-Auth-Key headers
		authHeader := r.Header.Get("Authorization")
		strictAuthKey := r.Header.Get("Strict-Auth-Key")

		// Ensure that only one of the two headers is present
		if (authHeader == "" && strictAuthKey == "") || (authHeader != "" && strictAuthKey != "") {
			response.SendError(w, []response.ErrorDetail{
				{Field: "Authorization/Strict-Auth-Key", Message: "Either Authorization or Strict-Auth-Key header must be provided, not both."},
			}, http.StatusBadRequest)
			return
		}

		// If Authorization header is present, validate the JWT token
		if authHeader != "" {
			// The Authorization header should have the format: "Bearer <token>"
			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			if tokenString == authHeader {
				response.SendBadRequestError(w, "Authorization header format is incorrect")
				return
			}

			// Parse and validate the JWT token
			secret := service.GetJWTSECRET()
			if secret == "" {
				response.SendError(w, []response.ErrorDetail{
					{Field: "JWT Secret", Message: "JWT secret key is missing"},
				}, http.StatusInternalServerError)
				return
			}

			// Parse the JWT token
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				// Ensure that the signing method is HMAC with SHA256
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, errors.New("unexpected signing method")
				}
				return []byte(secret), nil
			})
			if err != nil {
				response.SendError(w, []response.ErrorDetail{
					{Field: "token", Message: "Invalid or expired token"},
				}, http.StatusUnauthorized)
				return
			}

			// Extract claims from the JWT token
			if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
				// Store user information in the request context for downstream handlers
				user := &models.User{
					Username: claims["username"].(string),
					Email:    claims["email"].(string),
					Role:     claims["role"].(string),
				}

				// Store the user in context
				ctx := context.WithValue(r.Context(), userContextKey, user)
				r = r.WithContext(ctx)
				// Proceed to the next handler
				next.ServeHTTP(w, r)
			} else {
				response.SendError(w, []response.ErrorDetail{
					{Field: "token", Message: "Invalid or expired token"},
				}, http.StatusUnauthorized)
			}
		} else if strictAuthKey != "" {

			strictKey := service.GetStrictAuthKey()
			// Validate the Strict-Auth-Key (simple check, modify according to your needs)
			if strictAuthKey != strictKey {
				response.SendError(w, []response.ErrorDetail{
					{Field: "Strict-Auth-Key", Message: "Invalid Strict-Auth-Key"},
				}, http.StatusForbidden)
				return
			}
			// You could potentially set a dummy user or perform any other logic here
			// Proceed to the next handler
			next.ServeHTTP(w, r)
		}
	})
}

// CORS Middleware to handle cross-origin requests
func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Allow all origins for now (remove specific origins check)
		allowedMethods := []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
		allowedHeaders := []string{"Authorization", "Strict-Auth-Key", "Content-Type"}

		// Handle preflight OPTIONS request
		if r.Method == "OPTIONS" {
			// Respond with allowed methods and headers for preflight request
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", strings.Join(allowedMethods, ","))
			w.Header().Set("Access-Control-Allow-Headers", strings.Join(allowedHeaders, ","))
			w.Header().Set("Access-Control-Max-Age", "3600") // Cache preflight request for 1 hour
			w.WriteHeader(http.StatusNoContent)
			return
		}

		// For all other requests, check if the origin is allowed
		origin := r.Header.Get("Origin")
		if origin != "" {
			// Allow all origins (or specific ones if you change the allowedOrigins slice)
			w.Header().Set("Access-Control-Allow-Origin", "*")
		}

		// Set allowed methods and headers
		w.Header().Set("Access-Control-Allow-Methods", strings.Join(allowedMethods, ","))
		w.Header().Set("Access-Control-Allow-Headers", strings.Join(allowedHeaders, ","))
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		// Pass control to the next handler
		next.ServeHTTP(w, r)
	})
}

// MaxBodySizeMiddleware limits the size of request body.
func MaxBodySizeMiddleware(limit int64) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r.Body = http.MaxBytesReader(w, r.Body, limit)
			next.ServeHTTP(w, r)
		})
	}
}
