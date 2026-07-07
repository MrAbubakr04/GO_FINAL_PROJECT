package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/MrAbubakr04/GO_FINAL_PROJECT/internal/domain"
	"github.com/golang-jwt/jwt/v5"
)

type JWTService struct{}

func NewJWTService() *JWTService {
	return &JWTService{}
}

func (s *JWTService) GeneratePair(ctx context.Context, phoneNum string) (*domain.TokenPair, error) {
	_ = ctx
	claims := jwt.MapClaims{
		"sub": phoneNum,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte("dev-secret"))
	if err != nil {
		return nil, fmt.Errorf("sign access token: %w", err)
	}

	refreshClaims := jwt.MapClaims{
		"sub": phoneNum,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(7 * 24 * time.Hour).Unix(),
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	signedRefresh, err := refreshToken.SignedString([]byte("dev-secret-refresh"))
	if err != nil {
		return nil, fmt.Errorf("sign refresh token: %w", err)
	}

	return &domain.TokenPair{AccessToken: signed, RefreshToken: signedRefresh, ExpiresIn: int64(24 * time.Hour / time.Second)}, nil
}
