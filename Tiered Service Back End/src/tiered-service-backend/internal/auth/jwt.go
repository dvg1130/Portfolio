package auth

import (
	"context"
	"errors"
	"fmt"
	"tiered-service-backend/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

var secret = []byte(config.Config.JWT_SECRET_KEY)

// create jwt
func CreateToken(username string, role string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"role":     role,
		"exp":      jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
		"issued":   jwt.NewNumericDate(time.Now()),
	})
	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

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

// refresh token
func CreateRefreshToken(username string, role string) (string, error) {
	claims := jwt.MapClaims{
		"sub":  username,
		"role": role,
		"iat":  time.Now().Unix(),
		"exp":  time.Now().Add(7 * 24 * time.Hour).Unix(), // 7 days
		"typ":  "refresh_token",
		"jti":  uuid.New().String(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

// rotate refresh token
func RotateRefresh(ctx context.Context, redis *redis.Client, oldToken string) (string, string, error) {
	fmt.Println("Raw token:", oldToken)

	// verify old refresh token
	claims, err := VerifyToken(oldToken)
	if err != nil || claims["typ"] != "refresh_token" {
		fmt.Println("Token verification failed:", err)
		return "", "", errors.New("invalid token")
	}

	username := claims["sub"].(string)
	fmt.Println("Token claims:", claims)
	fmt.Println("Claim type:", claims["typ"])

	// check Redis to confirm token is still valid
	storedToken, err := redis.Get(ctx, username).Result()
	if err != nil || storedToken != oldToken {
		return "", "", errors.New("token mismatch or expired")
	}

	// delete old token
	redis.Del(ctx, username)

	// generate new tokens
	newAccessToken, err := CreateToken(username, claims["role"].(string))
	if err != nil {
		return "", "", err
	}

	newRefreshToken, err := CreateRefreshToken(username, claims["role"].(string))
	if err != nil {
		return "", "", err
	}

	// store new refresh token
	err = redis.Set(ctx, username, newRefreshToken, 7*24*time.Hour).Err()
	if err != nil {
		return "", "", err
	}

	return newAccessToken, newRefreshToken, nil
}
