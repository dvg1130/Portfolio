package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
)

func CreateRefreshToken(username string, UUID string, role string) (string, time.Time, error) {
	var exp = RefreshExpiry()
	claims := jwt.MapClaims{
		"sub":    username,
		"uuid":   UUID,
		"role":   role,
		"issued": time.Now().Unix(),
		"exp":    exp, // 7 days

		"typ": "refresh_token",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(secret)
	if err != nil {
		return "", time.Time{}, err
	}
	return signedToken, exp, err
}

func RotateRefreshToken(ctx context.Context, redis *redis.Client, oldToken string) (string, string, string, error) {
	//log for testing

	//verify old refreshtoken via jwt
	claims, err := VerifyToken(oldToken)
	if err != nil || claims["typ"] != "refresh_token" {
		fmt.Println("Token Validation failed: ", err)
		return "", "", "", errors.New("invalid token")
	}

	username := claims["sub"].(string)
	uuid := claims["uuid"].(string)
	role := claims["role"].(string)

	// check Redis to confirm token is still valid
	storedToken, err := redis.Get(ctx, username).Result()
	if err != nil || storedToken != oldToken {
		return "", "", "", errors.New("token mismatch or expired")
	}

	/// delete old token
	redis.Del(ctx, username)

	// generate new tokens
	newAccessToken, device_id, err := CreateAccessToken(username, uuid, role)
	if err != nil {
		return "", "", "", err
	}

	newRefreshToken, exp, err := CreateRefreshToken(username, uuid, role)
	if err != nil {
		return "", "", "", err
	}

	// store new refresh token
	ttl := time.Until(exp) // time.Duration until expiration
	err = redis.Set(ctx, username, newRefreshToken, ttl).Err()
	if err != nil {
		return "", "", "", err
	}

	return newAccessToken, newRefreshToken, device_id, nil

}

func RefreshExpiry() time.Time {
	return time.Now().Add(7 * 24 * time.Hour)
}
