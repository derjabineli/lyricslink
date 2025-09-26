package auth

import (
	"errors"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type CustomClaims struct {
	SessionID string `json:"sid"`
	TokenType TokenType `json:"typ"`
	jwt.RegisteredClaims
}


type TokenType string

const (
	AccessTokenType  TokenType = "lyriclink-access"
	RefreshTokenType TokenType = "lyriclink-refresh"
)

const (
	UserCookieName  string = "ll_user"
	RefreshCookieName string = "ll_refresh"
)


var (
	AccessTokenDuration  = 15 * time.Minute
	RefreshTokenDuration = 30 * 24 * time.Hour
)

func MakeJWT(userID uuid.UUID, sessionID uuid.UUID, tokenSecret string, expiresIn time.Duration, issuer TokenType) (string, error) {
	claims := CustomClaims{
		SessionID: sessionID.String(),
		TokenType: issuer,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    	string(issuer),
			IssuedAt:  	jwt.NewNumericDate(time.Now()),
			ExpiresAt: 	jwt.NewNumericDate(time.Now().Add(expiresIn)),
			Subject:   	userID.String(),
			ID:			uuid.NewString(), // jti
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(tokenSecret))
}

func ValidateJWT(tokenString, tokenSecret string, expectedType TokenType) (uuid.UUID, uuid.UUID, error) {
	claims := CustomClaims{}

	_, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(tokenSecret), nil
	})
	if err != nil {
		return uuid.Nil, uuid.Nil, err
	}

	// Verify issuer
	if claims.Issuer != string(expectedType) || claims.TokenType != expectedType {
		return uuid.Nil, uuid.Nil, errors.New("invalid token type or issuer")
	}

	// Verify expiration
	if claims.ExpiresAt == nil || time.Now().After(claims.ExpiresAt.Time) {
		return uuid.Nil, uuid.Nil, errors.New("token expired")
	}

	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return uuid.Nil, uuid.Nil, err
	}

	sessionID, err := uuid.Parse(claims.SessionID)
	if err != nil {
		return uuid.Nil, uuid.Nil, err
	}

	return userID, sessionID, nil
}

func NewAccessTokenCookie(userID uuid.UUID, sessionID uuid.UUID, tokenSecret string) (*http.Cookie, error) {
	jwtToken, err := MakeJWT(userID, sessionID, tokenSecret, AccessTokenDuration, AccessTokenType)
	if err != nil {
		return nil, err
	}

	cookie := &http.Cookie{
		Name:     UserCookieName,
		Value:    jwtToken,
		HttpOnly: true,                                // Make the cookie inaccessible to JavaScript
		Secure:  true,                                // Ensure the cookie is only sent over HTTPS
		SameSite: http.SameSiteLaxMode,                // Protect against CSRF attacks
		Expires:  time.Now().Add(AccessTokenDuration), // Set cookie expiration
		Path:     "/",                                 // Define cookie scope
	}

	return cookie, nil
}

func NewRefreshTokenCookie(userID uuid.UUID, sessionID uuid.UUID, tokenSecret string) (*http.Cookie, error) {
	jwtToken, err := MakeJWT(userID, sessionID, tokenSecret, RefreshTokenDuration, RefreshTokenType)
	if err != nil {
		return nil, err
	}

	cookie := &http.Cookie{
		Name:     RefreshCookieName,
		Value:    jwtToken,
		HttpOnly: true,                                 // Make the cookie inaccessible to JavaScript
		Secure:   true,                                 // Ensure the cookie is only sent over HTTPS
		SameSite: http.SameSiteLaxMode,              // Protect against CSRF attacks
		Expires:  time.Now().Add(RefreshTokenDuration), // Set cookie expiration
		Path:     "/",                                  // Define cookie scope
	}

	return cookie, nil
}

func ExtractSessionID(tokenString, tokenSecret string) (uuid.UUID, error) {
	claims := CustomClaims{}
	_, _, err := jwt.NewParser().ParseUnverified(tokenString, &claims)
	if err != nil {
		return uuid.Nil, err
	}

	return uuid.Parse(claims.SessionID)
}
