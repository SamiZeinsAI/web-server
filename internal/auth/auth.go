package auth

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// ErrNoAuthHeaderIncluded -
var ErrNoAuthHeaderIncluded = errors.New("not auth header included in request")

func AuthenticateUser(headers http.Header, secret string) (*jwt.Token, int, error) {
	tokenString, err := GetBearerToken(headers)
	if err != nil {
		msg := fmt.Sprintf("Error getting bearer token: %s\n", err)
		return nil, 0, errors.New(msg)
	}
	token, err := ParseToken(tokenString, secret)
	if err != nil {
		return nil, 0, err
	}
	issuer, err := token.Claims.GetIssuer()
	if err != nil {
		return nil, 0, err
	}
	if issuer == "chirpy-refresh" {
		return nil, 0, errors.New("Issuer matched refresh token, access token required")
	}
	id, err := token.Claims.GetSubject()
	if err != nil {
		return nil, 0, err
	}
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return nil, 0, err
	}
	return token, idInt, nil

}

func GetBearerToken(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", ErrNoAuthHeaderIncluded
	}
	splitAuth := strings.Split(authHeader, " ")
	if len(splitAuth) < 2 || splitAuth[0] != "Bearer" {
		return "", errors.New("malformed authorization header")
	}

	return splitAuth[1], nil
}

func ParseToken(tokenString, secret string) (*jwt.Token, error) {
	claimsStruct := jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(
		tokenString,
		&claimsStruct,
		func(key *jwt.Token) (interface{}, error) { return []byte(secret), nil },
	)
	if err != nil {
		return nil, err
	}
	return token, nil
}

func MakeToken(id int, issuer, secret string) (string, error) {
	duration := time.Duration(time.Hour)

	if issuer == "chirpy-refresh" {
		duration = time.Duration(time.Hour * 24 * 60)
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer: issuer,
		IssuedAt: &jwt.NumericDate{
			Time: time.Now().UTC(),
		},
		ExpiresAt: &jwt.NumericDate{
			Time: time.Now().UTC().Add(duration),
		},
		Subject: strconv.Itoa(id),
	})
	signedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}
