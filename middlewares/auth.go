package middlewares

import (
	"context"
	"net/http"
	"os"
	"p2p-lending/types"

	"p2p-lending/models"
	"p2p-lending/utils"

	"github.com/dgrijalva/jwt-go"
)

var JwtAuthentication = func(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "text/json")

		notAuth := []string{"/register", "/login"} //List of endpoints that doesn't require auth
		requestPath := r.URL.Path                      //current request path

		//check if request does not need authentication, serve the request if it doesn't need it
		for _, value := range notAuth {
			if value == requestPath {
				next.ServeHTTP(w, r)
				return
			}
		}

		tokenHeader := r.Header.Get("Authorization")

		//Token is missing, returns with error code 403 Unauthorized
		if tokenHeader == "" {
			utils.Response(w, http.StatusUnauthorized, types.Response.Unauthorized)
			return
		}

		tk := &models.Token{}

		// Parse the token
		token, err := jwt.ParseWithClaims(tokenHeader, tk, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("TOKEN_PASSWORD")), nil
		})

		//Malformed token, returns with http code 403 as usual
		if err != nil {
			utils.Response(w, http.StatusUnauthorized, types.Response.Unauthorized)
			return
		}

		// Find that user
		user := models.GetUserById(tk.UserId)

		// Check if the token is valid or if this user exists
		if !token.Valid || user.ID == "" { //Token is invalid, maybe not signed on this server
			utils.Response(w, http.StatusUnauthorized, types.Response.Unauthorized)
			return
		}

		// Pin the user id the the request
		ctx := context.WithValue(r.Context(), "user", tk.UserId)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r) // Call next middleware
	})
}
