package riAuth

import (
	"errors"
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

type IRiAdmin interface {
	SignToken() string
}

type riAuth struct {
	mapClaims *riMapClaims
	cfg       config.IJwtConfig
}

type riAdmin struct {
	*riAuth	
}

type riMapClaims struct{
	Claims *users.UserClaims `json:"claims"`
	jwt.RegisteredClaims 
}

func (a *riAuth) SignToken() string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, a.mapClaims)
	ss, _ := token.SignedString(a.cfg.SecretKey())
	return ss
}

func (a *riAdmin) SignToken() string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, a.mapClaims)
	ss, _ := token.SignedString(a.cfg.AdminKey())
	return ss
}



func jwtTimeDuration(t int) *jwt.NumericDate {
	return jwt.NewNumericDate(time.Now().Add(time.Duration(int64(t) * int64(math.Pow10(9)))))
}

func jwtTimeRepeatAdapter(t int64) *jwt.NumericDate {
	return jwt.NewNumericDate(time.Unix(t, 0))
}

func ParseToken(cfg config.IJwtConfig, tokenString string) (*riMapClaims, error){
	token, err := jwt.ParseWithClaims(tokenString, &riMapClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("signing method is invalid")
		}
		return cfg.SecretKey(), nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenMalformed){
			return nil, fmt.Errorf("token format is invalid")
		} else if errors.Is(err, jwt.ErrTokenExpired){
			return nil, fmt.Errorf("token is expired")
		} else {
			return nil, fmt.Errorf("parse token  failed : %v", err)
		}
	}

	if claims, ok := token.Claims.(*riMapClaims); ok {
		return claims, nil
	} else {
		return nil, fmt.Errorf("claims type is invalid")
	}

}

func ParseAdminToken(cfg config.IJwtConfig, tokenString string) (*riMapClaims, error){
	token, err := jwt.ParseWithClaims(tokenString, &riMapClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("signing method is invalid")
		}
		return cfg.AdminKey(), nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenMalformed){
			return nil, fmt.Errorf("token format is invalid")
		} else if errors.Is(err, jwt.ErrTokenExpired){
			return nil, fmt.Errorf("token is expired")
		} else {
			return nil, fmt.Errorf("parse token  failed : %v", err)
		}
	}

	if claims, ok := token.Claims.(*riMapClaims); ok {
		return claims, nil
	} else {
		return nil, fmt.Errorf("claims type is invalid")
	}

}


func RepeatToken(cfg config.IJwtConfig, claims *users.UserClaims, exp int64) string{
	obj := &riAuth{
		cfg: cfg,
		mapClaims: &riMapClaims{
			Claims: claims,
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    "rishop-api",
				Subject:  "refresh-token",
				Audience: []string{"customer", "admin"},
				ExpiresAt: jwtTimeRepeatAdapter(exp),
				NotBefore: jwt.NewNumericDate(time.Now()),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
			},
		},
	}
	return obj.SignToken()

}


func NewRiAuth(tokenType TokenType,cfg config.IJwtConfig, claims *users.UserClaims) (IRiAuth, error) {
	switch tokenType {
	case Access:
		return newAccessToken(cfg , claims), nil
	case Refresh:
		return newRefreshToken(cfg , claims), nil
	case Admin:
		return newAdminToken(cfg), nil
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

func newAdminToken(cfg config.IJwtConfig) IRiAuth {
	return &riAdmin{
		&riAuth{
			cfg: cfg,
			mapClaims: &riMapClaims{
				Claims: nil,
				RegisteredClaims: jwt.RegisteredClaims{
					Issuer:    "rishop-api",
					Subject:  "admin-token",
					Audience: []string{"admin"},
					ExpiresAt: jwtTimeDuration(300),
					NotBefore: jwt.NewNumericDate(time.Now()),
					IssuedAt:  jwt.NewNumericDate(time.Now()),
				},
			},
		},

	}
}




