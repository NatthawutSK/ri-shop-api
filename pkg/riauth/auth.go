package riAuth

import (
	"fmt"
	"math"
	"time"

	"github.com/NatthawutSK/ri-shop/config"
	"github.com/NatthawutSK/ri-shop/modules/users"
	"github.com/golang-jwt/jwt/v5"
)

type TokenType string

const (
	Access  TokenType = "access"
	Refresh TokenType = "refresh"
	Admin   TokenType = "admin"
	ApiKey  TokenType = "api_key"
)

type IRiAuth interface {
	SignToken() string
}

type riAuth struct {
	mapClaims *riMapClaims
	cfg       config.IJwtConfig
}

type riMapClaims struct{
	Claims *users.UserClaims `json:"claims"`
	jwt.RegisteredClaims 
}

func jwtTimeDuration(t int) *jwt.NumericDate {
	return jwt.NewNumericDate(time.Now().Add(time.Duration(int64(t) * int64(math.Pow10(9)))))
}

func jwtTimeRepeatAdapter(t int64) *jwt.NumericDate {
	return jwt.NewNumericDate(time.Unix(t, 0))
}


func NewRiAuth(tokenType TokenType,cfg config.IJwtConfig, claims *users.UserClaims) (IRiAuth, error) {
	switch tokenType {
	case Access:
		return newAccessToken(cfg , claims), nil
	case Refresh:
		return newRefreshToken(cfg , claims), nil
	default:
		return nil, fmt.Errorf("unknown token type")


	// return nil, nil
}
}


func newAccessToken(cfg config.IJwtConfig, claims *users.UserClaims) IRiAuth {
	return &riAuth{
		cfg: cfg,
		mapClaims: &riMapClaims{
			Claims: claims,
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    "rishop-api",
				Subject:  "access-token",
				Audience: []string{"customer", "admin"},
				ExpiresAt: jwtTimeDuration(cfg.AccessExpiresAt()),
				NotBefore: jwt.NewNumericDate(time.Now()),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
			},
		},
	}
}

func newRefreshToken(cfg config.IJwtConfig, claims *users.UserClaims) IRiAuth {
	return &riAuth{
		cfg: cfg,
		mapClaims: &riMapClaims{
			Claims: claims,
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    "rishop-api",
				Subject:  "refresh-token",
				Audience: []string{"customer", "admin"},
				ExpiresAt: jwtTimeDuration(cfg.RefreshExpiresAt()),
				NotBefore: jwt.NewNumericDate(time.Now()),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
			},
		},
	}
}




func (a *riAuth) SignToken() string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, a.mapClaims)
	ss, _ := token.SignedString(a.cfg.SecretKey())
	return ss
}