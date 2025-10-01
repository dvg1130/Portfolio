package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/dvg1130/Portfolio/secure-backend/config"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var secret = []byte(config.AuthConfig.JWT_SECRET_KEY)

// create jwt access token
func CreateAccessToken(username string, UUID string, role string) (string, string, error) {
	//create device_id
	var device_id = uuid.New().String()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"uuid":     UUID,
		"role":     role,
		"issued":   jwt.NewNumericDate(time.Now()),
		"exp":      jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
	})

	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", "", err
	}
	return tokenString, device_id, nil
}

// verify access token
func VerifyToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		//validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secret, nil
	})

	if err != nil {
		return nil, err
	}

	//return claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}
