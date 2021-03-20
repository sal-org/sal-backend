package util

import (
	"fmt"
	CONFIG "salbackend/config"
	CONSTANT "salbackend/constant"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// CheckIfUserLogin - check if token is valid, not expired and token user id is same as user id given
func CheckIfUserLogin(auth string, userID string) bool {
	id, ok, access := parseJWTAccessToken(auth)
	if !ok || !access || !strings.EqualFold(id, userID) {
		return false
	}
	return true
}

// GetUserIDFromJWTToken - get user id from token
func GetUserIDFromJWTToken(token string) string {
	id, _, _ := parseJWTAccessToken(token)
	return id
}

// CreateAccessToken - jwt token for accessing api
func CreateAccessToken(id string) (string, bool) {
	return createJWTToken(id, CONSTANT.JWTAccessExpiry)
}

// CreateRefreshToken - jwt token for getting access token, if expired
func CreateRefreshToken(id string) (string, bool) {
	return createJWTToken(id, CONSTANT.JWTRefreshExpiry, true)
}

func createJWTToken(id string, expiry int, refreshToken ...bool) (string, bool) {
	claims := jwt.MapClaims{}
	claims["id"] = id
	claims["exp"] = strconv.FormatInt(time.Now().Add(time.Minute*time.Duration(expiry)).Unix(), 10)
	if len(refreshToken) > 0 {
		claims["refresh"] = "1"
	}
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := at.SignedString(CONFIG.JWTSecret)
	if err != nil {
		return "", false
	}
	return token, true
}

func extractJWTToken(authorization string) string {
	bearToken := authorization
	strArr := strings.Split(bearToken, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}
	return ""
}

func verifyJWTToken(authorization string) (*jwt.Token, bool) {
	token, err := jwt.Parse(extractJWTToken(authorization), func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return CONFIG.JWTSecret, nil
	})
	if err != nil {
		return nil, false
	}
	// check if token valid
	if _, ok := token.Claims.(jwt.Claims); !ok {
		return nil, false
	}
	if !token.Valid {
		return nil, false
	}
	return token, true
}

func parseJWTAccessToken(authorization string) (string, bool, bool) {
	token, ok := verifyJWTToken(authorization)
	if !ok {
		return "", false, false
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", false, false
	}

	// extract user id
	id, ok := claims["id"].(string)
	if !ok {
		return "", false, false
	}
	// extract expiry
	exp, ok := claims["exp"].(string)
	if !ok {
		return "", false, false
	}
	expiry, err := strconv.ParseInt(exp, 10, 64)
	if err != nil {
		return "", false, false
	}
	// check if token expired
	if expiry < time.Now().Unix() { // expired if less than current time
		return "", false, false
	}
	if claims["refresh"] != nil && strings.EqualFold(claims["refresh"].(string), "1") {
		return id, true, false
	}
	return id, true, true
}
