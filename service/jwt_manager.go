package service

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"
)

type JwtManager struct {
	secretKey     string
	tokenDuration time.Duration
}

type UserClaims struct {
	jwt.StandardClaims
	Username string `json:"username"`
	Role     string `json:"role"`
}

func NewJwtManager(secretKey string, tokenDuration time.Duration) *JwtManager {
	return &JwtManager{
		secretKey:     secretKey,
		tokenDuration: tokenDuration,
	}
}

func (manager *JwtManager) Generate(user *User) (string, error) {
	claims := UserClaims{
		Username:       user.Username,
		Role:           user.Role,
		StandardClaims: jwt.StandardClaims{ExpiresAt: time.Now().Add(manager.tokenDuration).Unix()},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(manager.secretKey))
}

func (manager *JwtManager) Parse(accessToken string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(
		accessToken,
		&UserClaims{},
		func(token *jwt.Token) (interface{}, error) {
			// 校验签名方法
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, fmt.Errorf("unexpected token signing method")
			}
			return []byte(manager.secretKey), nil
		})
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}
	// 校验加密claims
	claims, ok := token.Claims.(*UserClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}
	return claims, nil
}
