package jwt

import (
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
	v1 "github.com/sonalys/goshare/internal/application/v1"
	"github.com/sonalys/goshare/internal/domain"
)

type (
	Client struct {
		jwtSignKey []byte
	}
)

func NewClient(jwtSignKey []byte) *Client {
	return &Client{
		jwtSignKey: jwtSignKey,
	}
}

func (c *Client) Decode(tokenString string) (*v1.Identity, error) {
	var claims jwt.MapClaims

	keyFunc := func(token *jwt.Token) (any, error) {
		return c.jwtSignKey, nil
	}

	supportedMethods := []string{
		jwt.SigningMethodHS256.Alg(),
	}

	token, err := jwt.ParseWithClaims(
		tokenString,
		&claims,
		keyFunc,
		jwt.WithValidMethods(supportedMethods),
	)
	if err != nil {
		return nil, fmt.Errorf("parsing token: %w", err)
	}

	if !token.Valid {
		return nil, v1.ErrAuthenticationExpired
	}

	email, ok := claims["email"].(string)
	if !ok {
		return nil, errors.New("missing email claim")
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		return nil, errors.New("missing user_id claim")
	}

	userUUID, err := domain.ParseID(userID)
	if err != nil {
		return nil, fmt.Errorf("parsing user_id: %w", err)
	}

	exp, ok := claims["exp"].(float64)
	if !ok {
		return nil, errors.New("missing exp claim")
	}

	identity := &v1.Identity{
		Email:  email,
		UserID: userUUID,
		Exp:    int64(exp),
	}

	return identity, nil
}

func (c *Client) Encode(identity *v1.Identity) (string, error) {
	claims := jwt.MapClaims{
		"email":   identity.Email,
		"user_id": identity.UserID.String(),
		"exp":     identity.Exp,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(c.jwtSignKey)
	if err != nil {
		return "", fmt.Errorf("signing token: %w", err)
	}

	return tokenString, nil
}
