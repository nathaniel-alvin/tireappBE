package auth

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	tireappbe "github.com/nathaniel-alvin/tireappBE"
	"github.com/nathaniel-alvin/tireappBE/config"
	"github.com/nathaniel-alvin/tireappBE/types"
	"github.com/nathaniel-alvin/tireappBE/utils"
)

var secret = []byte(config.Envs.JWTSecret)

type Claims struct {
	UserID int `json:"userId"`
	jwt.RegisteredClaims
}

func WithJWTAuth(handlerFunc http.HandlerFunc, store types.UserRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("token")
		if err != nil {
			if err == http.ErrNoCookie {
				permissionDenied(w)
				return
			}
			utils.WriteError(w, http.StatusBadRequest, err)
			return
		}
		tokenString := c.Value

		token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			return secret, nil
		})
		if err != nil {
			log.Printf("failed to validate token: %v", err)
			permissionDenied(w)
			return
		}

		if !token.Valid {
			log.Printf("invalid token")
			permissionDenied(w)
			return
		}

		var userID int
		claims, ok := token.Claims.(*Claims)
		if ok && token.Valid {
			userID = claims.UserID
		} else {
			log.Printf("invalid token for claims")
			permissionDenied(w)
			return
		}

		u, err := store.GetUserById(r.Context(), userID)
		if err != nil {
			log.Printf("failed to get user by id: %v", err)
			permissionDenied(w)
			return
		}

		// add user to context
		ctx := r.Context()
		ctx = tireappbe.NewContextWithUser(ctx, u)
		r = r.WithContext(ctx)

		// call the function if the token is valid
		handlerFunc(w, r)
	}
}

func CreateTokensAndSetCookies(w http.ResponseWriter, userID int, expireDuration int64) (string, error) {
	accessToken, exp, err := createJWT(userID, expireDuration)
	if err != nil {
		return "", err
	}

	setCookie(w, "token", accessToken, exp)

	return accessToken, nil
}

func TokenRefresher(w http.ResponseWriter, c *http.Cookie, expireDuration int64) error {
	tknStr := c.Value
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			// utils.WriteError(w, http.StatusUnauthorized, err)
			return err
		}
		// utils.WriteError(w, http.StatusBadRequest, err)
		return err
	}
	if !token.Valid {
		// utils.WriteError(w, http.StatusUnauthorized, err)
		return err
	}

	// cannot refresh token when expire time > 30sec
	if time.Until(claims.ExpiresAt.Time) > 30*time.Second {
		return err
	}

	var userID int
	if claims == nil {
		return fmt.Errorf("empty claim")
	}
	userID = claims.UserID
	// if ok && token.Valid && claims != nil {
	// 	userID = claims.UserID
	// } else {
	// 	log.Printf("invalid token for claims")
	// 	permissionDenied(w)
	// 	return err
	// }

	newToken, _, err := createJWT(userID, expireDuration)
	if err != nil {
		// utils.WriteError(w, http.StatusInternalServerError, err)
		return err
	}

	expirationDuration := time.Second * time.Duration(expireDuration)
	expirationTime := time.Now().Add(expirationDuration)
	setCookie(w, "token", newToken, expirationTime)

	return nil
}

func createJWT(userID int, expireDuration int64) (string, time.Time, error) {
	expirationDuration := time.Second * time.Duration(expireDuration)
	expirationTime := time.Now().Add(expirationDuration)
	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", time.Now(), err
	}
	return tokenString, expirationTime, nil
}

func setCookie(w http.ResponseWriter, name, token string, expiration time.Time) {
	cookie := new(http.Cookie)
	cookie.Name = name
	cookie.Value = token
	cookie.Expires = expiration
	cookie.Path = "/"
	cookie.HttpOnly = true

	http.SetCookie(w, cookie)
}

func permissionDenied(w http.ResponseWriter) {
	utils.WriteError(w, http.StatusForbidden, fmt.Errorf("permission denied"))
}

func ValidateJWT(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secret, nil
	})
}
