package auth

import (
	"context"
)

type contextKey string

var claimsKey = contextKey("claims")

// add claims
func AddClaimsToContext(ctx context.Context, claims map[string]interface{}) context.Context {
	return context.WithValue(ctx, claimsKey, claims)
}

// retrieve claims
func GetClaimsFromContext(ctx context.Context) map[string]interface{} {
	val := ctx.Value(claimsKey)
	if claims, ok := val.(map[string]interface{}); ok {
		return claims
	}
	return nil
}
