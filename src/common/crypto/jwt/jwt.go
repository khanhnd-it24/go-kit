package jwtprovider

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type Token struct {
	Token     string `json:"token"`
	ExpiresIn int64  `json:"expires_in"`
}

type Payload struct {
	Id       string `json:"id"`
	Username string `json:"username"`
}

type CanBeEmpty interface {
	IsEmpty() bool
}

type myClaims struct {
	Id       string `json:"id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func Generate(secret string, payload Payload, expiry int64) (*Token, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, myClaims{
		payload.Id,
		payload.Username,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Second * time.Duration(expiry))),
		},
	})

	myToken, err := t.SignedString([]byte(secret))
	if err != nil {
		return nil, err
	}

	return &Token{
		Token:     myToken,
		ExpiresIn: expiry,
	}, nil
}

var ErrTokenExpired = errors.New("token expired")

func Verify(secret string, myToken string) (*Payload, error) {
	res, err := jwt.ParseWithClaims(myToken, &myClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if errors.Is(err, jwt.ErrTokenExpired) {
		return nil, ErrTokenExpired
	}

	if err != nil {
		return nil, err
	}

	if !res.Valid {
		return nil, fmt.Errorf("invalid token %s", myToken)
	}
	claims, ok := res.Claims.(*myClaims)
	if !ok {
		return nil, fmt.Errorf("invalid claims with token %s", myToken)
	}

	return &Payload{
		Id:       claims.Id,
		Username: claims.Username,
	}, nil
}
