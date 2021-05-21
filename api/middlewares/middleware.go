package middlewares

import (
	"context"
	"net/http"
	"os"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"

	"github.com/SarthakJain26/GoAuth/api/responses"
)

// SetContentMiddleware sets content type to json
func SetContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

// The AuthJWTVerify verifies user authentication tokens which are acquired at login.
// After successful verification, the user's ID is added to the request context.
// This middleware will only allow authenticated users access to protected API routes such as Venue creation.

// AuthJwtVerify verify token and add userID to the request context
func AuthJwtVerify(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var resp = map[string]interface{}{"status": "failed", "message": "Missing authorization token"}
		var header = r.Header.Get("Authorization")
		header = strings.TrimSpace(header)
		if header == "" {
			responses.JSON(w, http.StatusForbidden, resp)
			return
		}

		token, err := jwt.Parse(header, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("SECRET")), nil
		})
		if err != nil {
			resp["status"] = "failed"
			resp["message"] = "Invalid token, please login"
			responses.JSON(w, http.StatusForbidden, resp)
			return
		}
		claims, _ := token.Claims.(jwt.MapClaims)

		ctx := context.WithValue(r.Context(), "userID", claims["userID"]) // adding the user ID to the context
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
