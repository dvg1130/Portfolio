package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/dvg1130/Portfolio/secure-backend/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
)

func CreateRefreshToken(username string, UUID string, role string) (string, time.Time, error) {
	var exp = RefreshExpiry()
	claims := jwt.MapClaims{
		"sub":    username,
		"uuid":   UUID,
		"role":   role,
		"issued": jwt.NewNumericDate(time.Now()),
		"exp":    jwt.NewNumericDate(exp), // 7 days

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
	// Verify JWT first
	claims, err := VerifyToken(oldToken)
	if err != nil || claims["typ"] != "refresh_token" {
		fmt.Println("Token validation failed:", err)
		return "", "", "", errors.New("invalid token")
	}

	username, _ := claims["sub"].(string)
	uuid, _ := claims["uuid"].(string)
	role, _ := claims["role"].(string)

	// Get session JSON from Redis
	key := fmt.Sprintf("refresh:%s", username)
	sessionJSON, err := redis.Get(ctx, key).Result()
	if err != nil {
		return "", "", "", errors.New("session not found")
	}

	var session models.RefreshSession
	if err := json.Unmarshal([]byte(sessionJSON), &session); err != nil {
		return "", "", "", errors.New("invalid session format")
	}

	// Compare tokens
	if session.RefreshToken != oldToken {
		return "", "", "", errors.New("token mismatch")
	}

	// Delete old session
	redis.Del(ctx, key)

	// Generate new tokens
	newAccessToken, deviceID, err := CreateAccessToken(username, uuid, role)
	if err != nil {
		return "", "", "", err
	}

	newRefreshToken, exp, err := CreateRefreshToken(username, uuid, role)
	if err != nil {
		return "", "", "", err
	}

	// Save new session
	newSession := models.RefreshSession{
		RefreshToken: newRefreshToken,
		DeviceID:     deviceID,
		ExpiresAt:    exp,
	}
	newSessionJSON, _ := json.Marshal(newSession)
	ttl := time.Until(exp)

	err = redis.Set(ctx, key, newSessionJSON, ttl).Err()
	if err != nil {
		return "", "", "", err
	}

	return newAccessToken, newRefreshToken, deviceID, nil
}

func RefreshExpiry() time.Time {
	return time.Now().Add(7 * 24 * time.Hour)
}
