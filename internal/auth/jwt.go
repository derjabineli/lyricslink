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
	AccessTokenType TokenType = "lyriclink-access"
	RefreshTokenType TokenType = "lyriclink-refresh"
)

var (
	AccessTokenDuration  = 15 * time.Minute
	RefreshTokenDuration = 30 * 24 * time.Hour
)

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration, issuer TokenType) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer: string(issuer),
		IssuedAt: jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
		Subject: userID.String(),
	})
	return token.SignedString([]byte(tokenSecret))
}

func ValidateJWT(tokenString, tokenSecret string, issuerType TokenType) (uuid.UUID, error) {
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
	if issuer != string(issuerType) {
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

func NewAccessTokenCookie(id uuid.UUID, tokenSecret string) (*http.Cookie, error) {
	jwtToken, err := MakeJWT(id, tokenSecret, AccessTokenDuration, AccessTokenType)
	if err != nil {
		return nil, err
	}

	cookie := &http.Cookie{
		Name:     "ll_user",
		Value:    jwtToken,
		HttpOnly: true, // Make the cookie inaccessible to JavaScript
		Secure:   true, // Ensure the cookie is only sent over HTTPS
		SameSite: http.SameSiteLaxMode, // Protect against CSRF attacks
		Expires:  time.Now().Add(AccessTokenDuration), // Set cookie expiration
		Path:     "/", // Define cookie scope
	}

	return cookie, nil
}

func NewRefreshTokenCookie(id uuid.UUID, tokenSecret string) (*http.Cookie, error) {
	jwtToken, err := MakeJWT(id, tokenSecret, RefreshTokenDuration, RefreshTokenType)
	if err != nil {
		return nil, err
	}

	cookie := &http.Cookie{
		Name:     "ll_refresh",
		Value:    jwtToken,
		HttpOnly: true, // Make the cookie inaccessible to JavaScript
		Secure:   true, // Ensure the cookie is only sent over HTTPS
		SameSite: http.SameSiteStrictMode, // Protect against CSRF attacks
		Expires:  time.Now().Add(RefreshTokenDuration), // Set cookie expiration
		Path:     "/", // Define cookie scope
	}

	return cookie, nil
}