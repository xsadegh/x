package auth

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

var (
	ErrForbidden     = errors.New("access denied")
	ErrExpiredToken  = errors.New("token is expired")
	ErrInvalidToken  = errors.New("token is invalid")
	ErrNotAuthorized = errors.New("not authorized")
)

const ContextKey = "AUTH_CONTEXT"

type Context struct {
	Secret []byte
	Header string
	UserID uuid.UUID
}

func parseHeader(secret []byte, raw string) (*jwt.Token, error) {
	token, err := jwt.Parse(raw, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return secret, nil
	})

	if token != nil {
		if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			return token, nil
		}
	}

	return nil, err
}

func checkError(ctx context.Context, next graphql.Resolver, err error) (any, error) {
	errs := graphql.GetErrors(ctx)
	if len(errs) == 0 {
		return nil, err
	}

	return next(ctx)
}

func GenerateToken(ctx context.Context, exp time.Duration, id string) (string, error) {
	token := &jwt.StandardClaims{ExpiresAt: time.Now().Add(exp).Unix(), Id: id}
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, token)
	secret := ctx.Value(ContextKey).(Context).Secret
	raw, err := claims.SignedString(secret)
	if err != nil {
		return "", fmt.Errorf("unexpected jwt signing error: %v", err)
	}

	return raw, nil
}

func Authenticate(ctx context.Context, _ any, next graphql.Resolver) (any, error) {
	secret := ctx.Value(ContextKey).(Context).Secret
	header := ctx.Value(ContextKey).(Context).Header
	if header == "" {
		return checkError(ctx, next, ErrNotAuthorized)
	}

	header = strings.Replace(header, "Bearer ", "", 1)
	token, err := parseHeader(secret, header)
	if err != nil {
		switch err {
		case ErrExpiredToken:
			return checkError(ctx, next, ErrExpiredToken)
		default:
			return checkError(ctx, next, ErrInvalidToken)
		}
	}

	if token.Valid {
		c := ctx.Value(ContextKey).(Context)
		c.UserID, _ = uuid.Parse(token.Claims.(jwt.MapClaims)["jti"].(string))
		return next(context.WithValue(ctx, ContextKey, c))
	}

	return checkError(ctx, next, ErrForbidden)
}

func GetUserID(ctx context.Context) uuid.UUID {
	return ctx.Value(ContextKey).(Context).UserID
}
