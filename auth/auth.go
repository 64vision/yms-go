package auth

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	u "gollux/utils"

	jwt "github.com/dgrijalva/jwt-go"
)

type Token struct {
	AccountID int
	jwt.StandardClaims
}

var JwtAuthentication = func(next http.Handler) http.Handler {
	now := time.Now()
	accessKey := u.Md5hash(now.Format("2006-02-01"))
	fmt.Println(now.Format("2006-02-01"), accessKey)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		notAccessKey := []string{
			"/race/vmix",
			"/race/vmix_odds",
			"/race/vmix_lap",
			"/race/vmix_prizes",
			"/race/vmix_recent",
			"/credits/callback",
			"/credits/maya_callback",
		} //List of endpoints that doesn't require access-key
		notAuth := []string{
			"/account/register",
			"/account/register",
			"/account/verify",
			"/account/resendcode",
			"/account/login_by_password",
			"/account/forgot",
			"/admin/login",
			"/get_access",
			"/race/vmix",
		} //List of endpoints that doesn't require auth 09091666962

		requestPath := r.URL.Path

		for _, value := range notAccessKey {

			if value == requestPath {
				next.ServeHTTP(w, r)
				return
			}
		}

		response := make(map[string]interface{}) //current request path
		// Remove access token
		// access := r.Header.Get("access-token")
		// if !access != accessKey {
		// 	response = u.Message(false, "Access Token Invalid")
		// 	w.WriteHeader(http.StatusForbidden)
		// 	w.Header().Add("Content-Type", "application/json")
		// 	u.Respond(w, response)
		// 	return
		// }

		///check public
		// if strings.Contains(requestPath, "/public/") == true {
		// 	next.ServeHTTP(w, r)
		// 	return
		// }
		//check if request does not need authentication, serve the request if it doesn't need it
		for _, value := range notAuth {
			if value == requestPath {
				next.ServeHTTP(w, r)
				return
			}
		}
		//fmt.Println("JwtAuthentication 3")

		tokenHeader := r.Header.Get("Authorization") //Grab the token from the header
		//fmt.Println(tokenHeader)
		if tokenHeader == "" { //Token is missing, returns with error code 403 Unauthorized
			response = u.Message(false, "Missing Authorization Token")
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			u.Respond(w, response)
			return
		}
		//fmt.Println("JwtAuthentication 4")

		splitted := strings.Split(tokenHeader, " ") //The token normally comes in format `Bearer {token-body}`, we check if the retrieved token matched this requirement
		if len(splitted) != 2 {
			response = u.Message(false, "Invalid/Malformed auth token")
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			u.Respond(w, response)
			return
		}

		tokenPart := splitted[1] //Grab the token part, what we are truly interested in
		//tk := Token{}
		fmt.Println(tokenPart)

		token, err := ParseToken(tokenPart)
		if err != nil {
			fmt.Println(err.Error())
			response = u.Message(false, "Malformed authentication token")
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			u.Respond(w, response)
			return
		} else {
			fmt.Println("Parsed Token:", token)
			fmt.Println("Parsed accountID", token.AccountID)
		}

		//fmt.Printf("%v|USER:%d|ADDRESS:%s|URI:%s\n", time.Now(), tk.UserId, r.RemoteAddr, r.URL.Path)
		ctx := context.WithValue(r.Context(), "user", token.AccountID)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r) //proceed in the middleware chain!
	})
}

func ParseToken(tokenString string) (*Token, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Token{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("token_mydog_app"), nil
	})

	if err != nil {
		return nil, err
	}

	// Extract and validate claims
	if claims, ok := token.Claims.(*Token); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}
