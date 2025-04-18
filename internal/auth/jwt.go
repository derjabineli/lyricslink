package auth

import (
	"errors"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TokenType string

const (
	TokenTypeAccess TokenType = "lyriclink-access"
)

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer: string(TokenTypeAccess),
		IssuedAt: jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
		Subject: userID.String(),
	})
	return token.SignedString([]byte(tokenSecret))
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	claims := jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(
		tokenString, 
		&claims, 
		func(token *jwt.Token) (interface{}, error) {return []byte(tokenSecret), nil},
	)
	if err != nil {
		return uuid.Nil, err
	}

	userIDString, err := token.Claims.GetSubject()
	if err != nil {
		return uuid.Nil, err
	}
	issuer, err := token.Claims.GetIssuer()
	if err != nil {
		return uuid.Nil, err
	}
	if issuer != string(TokenTypeAccess) {
		return uuid.Nil, errors.New("invalid issuer")
	}
	expiration, err := token.Claims.GetExpirationTime()
	if err != nil {
		return uuid.Nil, err
	}
	isExpired := checkTokenExpiration(expiration.Time)
	if isExpired {
		return uuid.Nil, errors.New("token expired")
	}
	
	id, err := uuid.Parse(userIDString)
	if err != nil {
		return uuid.Nil, err
	}
	return id, nil
}

func checkTokenExpiration(expiration time.Time) bool{
	now := time.Now()

	return now.After(expiration)
}

func NewJWTCookie(id uuid.UUID, tokenSecret string, expiresIn time.Duration) (*http.Cookie, error) {
	jwtToken, err := MakeJWT(id, tokenSecret, expiresIn)
	if err != nil {
		return nil, err
	}

	cookie := &http.Cookie{
		Name:     "ll_user",
		Value:    jwtToken,
		HttpOnly: true, // Make the cookie inaccessible to JavaScript
		Secure:   false, // Ensure the cookie is only sent over HTTPS
		SameSite: http.SameSiteLaxMode, // Protect against CSRF attacks
		Expires:  time.Now().Add(24 * time.Hour), // Set cookie expiration
		Path:     "/", // Define cookie scope
	}

	return cookie, nil
}