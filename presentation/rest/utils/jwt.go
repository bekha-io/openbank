package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)


var (
	jwtSigningKey = []byte("some-random-jwt-signing-key") // TODO: ONLY FOR TESTING PURPOSES
)


func GenerateJwtToken(expiresAt time.Time, values map[string]interface{}) (string, error) {
	claims := jwt.MapClaims{
		"iss": "openbank-api-server",
		"iat": time.Now().UTC().Unix(),
		"exp": expiresAt.Unix(),
	}
	for k, v := range values {
		claims[k] = v
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := t.SignedString(jwtSigningKey)
	if err != nil {
		return "", err
	}

	return signed, nil
}

func ParseToken(token string) (jwt.MapClaims, error) {
	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return jwtSigningKey, nil
	}, jwt.WithExpirationRequired(), jwt.WithIssuedAt())
	if err != nil {
		return nil, err
	}

	if !parsedToken.Valid {
		return nil, errors.New("token is not valid")
	}

	return parsedToken.Claims.(jwt.MapClaims), nil
}