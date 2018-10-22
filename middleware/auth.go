package middleware

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"learnabout-filemanager/config"
	"net/http"
	"strings"
)

func Authentication(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		authenticated, err := authenticate(r, config.AppSecret)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "Internal Server Error")
			return
		}

		if authenticated == false {
			w.WriteHeader(http.StatusForbidden)
			fmt.Fprint(w, "Not Authorized")
			return
		}

		handler.ServeHTTP(w, r)
	})
}


func authenticate(r *http.Request, appSecret []byte) (bool, error) {
	authorization := r.Header.Get("Authorization")
	if authorization == "" {
		return false, nil
	}

	authorization = strings.Replace(authorization, "Bearer ", "", 1)

	parser := jwt.Parser{
		SkipClaimsValidation: true,
	}

	token, err := parser.Parse(authorization, func(t *jwt.Token) (interface{}, error){

		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", t.Header["alg"])
		}

		return appSecret, nil
	})

	if err != nil {
		fmt.Printf("error in parsing token: %v", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Println(claims["userId"])
		return true, nil
	} else {
		fmt.Println(err)
		return false, nil
	}
}
