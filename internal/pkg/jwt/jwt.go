package jwt

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
	v1 "github.com/sonalys/goshare/internal/pkg/v1"
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
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		return c.jwtSignKey, nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, v1.ErrAuthenticationExpired
	}

	email, ok := claims["email"].(string)
	if !ok {
		return nil, fmt.Errorf("missing email claim")
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		return nil, fmt.Errorf("missing user_id claim")
	}

	userUUID, err := v1.ParseID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to parse user_id: %v", err)
	}

	exp, ok := claims["exp"].(float64)
	if !ok {
		return nil, fmt.Errorf("missing exp claim")
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
		return "", fmt.Errorf("failed to sign token: %v", err)
	}

	return tokenString, nil
}
