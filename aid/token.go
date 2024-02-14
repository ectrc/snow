package aid

import (
	"fmt"
	"regexp"

	"github.com/golang-jwt/jwt/v5"
)

var JWTCompile = regexp.MustCompile(`eg1~(.*)`)

func JWTSign(m JSON) (string, error) {
	claims := jwt.MapClaims{}

	for k, v := range m {
		claims[k] = v
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token.Header["kid"] = "g__VKjSSmqJ0xZj1RYkLGKQ7dnHiM9MLhFVwKPySDB4"
	return token.SignedString([]byte(Config.JWT.Secret))
}

func JWTVerify(tokenString string) (JSON, error) {
	compiled := JWTCompile.FindStringSubmatch(tokenString)
	if len(compiled) > 0 {
		tokenString = compiled[1]
	}
	
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

func GetSnowFromToken(token string) (string, error) {
	claims, err := JWTVerify(token)
	if err != nil {
		return "", err
	}

	if claims["snow_id"] == nil {
		return "", fmt.Errorf("invalid access token")
	}

	snowId, ok := claims["snow_id"].(string)
	if !ok {
		return "", fmt.Errorf("invalid access token")
	}

	return snowId, nil
}