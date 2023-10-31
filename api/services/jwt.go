package services

import (
	"github.com/dgrijalva/jwt-go"
	"os"
	"strings"
	"time"
)

var jwtKey = []byte(os.Getenv("SECRET_KEY"))
var anotherJwtKey = []byte(os.Getenv("ANOTHER_SECRET_KEY"))
var refreshTokenKey = []byte(os.Getenv("REFRESH_TOKEN_KEY"))

type Claims struct {
	UserID string `json:"email"`
	jwt.StandardClaims
}

func GenerateToken(userID string) (string, string, error) {
	expirationTime := time.Now().Add(4320 * time.Minute)
	refreshExpTime := time.Now().Add(4320 * time.Minute)

	claims := &Claims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	refreshClaims := &Claims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: refreshExpTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)

	tokenString, err := token.SignedString(jwtKey)
	refreshTokenString, err := refreshToken.SignedString(refreshTokenKey)
	return tokenString, refreshTokenString, err
}

func GenerateNonAuthToken(userID string) (string, error) {
	expirationTime := time.Now().Add(4320 * time.Minute)
	claims := &Claims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(anotherJwtKey)

	return tokenString, err
}

func DecodeNonAuthToken(tkStr string) (string, error) {
	claims := &Claims{}

	tkn, err := jwt.ParseWithClaims(tkStr, claims, func(token *jwt.Token) (interface{}, error) {
		return anotherJwtKey, nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return "", err
		}
		return "", err
	}

	if !tkn.Valid {
		return "", err
	}

	return claims.UserID, nil
}

func DecodeToken(tkStr string) (string, error) {
	claims := &Claims{}
	tkStr = strings.Replace(tkStr, "Bearer ", "", 1)
	//print(tkStr)
	tkn, err := jwt.ParseWithClaims(tkStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return "", err
		}
		return "", err
	}

	if !tkn.Valid {
		return "", err
	}

	return claims.UserID, nil
}

func DecodeRefreshToken(tkStr string) (string, error) {
	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(tkStr, claims, func(token *jwt.Token) (interface{}, error) {
		return refreshTokenKey, nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return "", err
		}
		return "", err
	}

	if !tkn.Valid {
		return "", err
	}

	return claims.UserID, nil
}
