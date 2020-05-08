package auth

import (
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"time"
)

//a middleware returns a func that works like a "filter"
//if the func does not call next.ServeHTTP the chain will be interrupted
//this filter is applied at every incoming HTTP request
func CheckUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Our middleware logic goes here...

		// We can obtain the session token from the requests cookies, which come with every request
		c, err := r.Cookie("token")
		if err != nil {
			if err == http.ErrNoCookie {
				// If the cookie is not set, return an unauthorized status
				http.Error(w, "cookie not present", http.StatusUnauthorized)
				return
			}
			// For any other type of error, return a bad request status
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Get the JWT string from the cookie
		tknStr := c.Value

		// Initialize a new instance of `Claims`
		claims := &Claims{}

		// Parse the JWT string and store the result in `claims`.
		// Note that we are passing the key in this method as well. This method will return an error
		// if the token is invalid (if it has expired according to the expiry time we set on sign in),
		// or if the signature does not match
		tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if !tkn.Valid {
			http.Error(w, "token not valid", http.StatusUnauthorized)
			return
		}

		// We ensure that a new token is not issued until enough time has elapsed
		// In this case, a new token will only be issued if the old token is within
		// 30 seconds of expiry.
		if time.Unix(claims.ExpiresAt, 0).Sub(time.Now()) < 30*time.Second {
			// Now, create a new token for the current use, with a renewed expiration time
			expirationTime := time.Now().Add(5 * time.Minute)
			claims.ExpiresAt = expirationTime.Unix()
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
			tokenString, err := token.SignedString(jwtKey)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			// Set the new token as the users `token` cookie
			http.SetCookie(w, &http.Cookie{
				Name:    "token",
				Value:   tokenString,
				Expires: expirationTime,
			})
		}

		//ok proceed to next middleware or final handler code
		next.ServeHTTP(w, r)
	})
}
