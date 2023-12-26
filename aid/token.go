package aid

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

func JWTSign(m JSON) (string, error) {
	claims := jwt.MapClaims{}

	for k, v := range m {
		claims[k] = v
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(Config.JWT.Secret))
}

func JWTVerify(tokenString string) (JSON, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(Config.JWT.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid claims")
	}

	json := JSON{}
	for k, v := range claims {
		json[k] = v
	}

	return json, nil
}