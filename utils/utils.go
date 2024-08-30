package utils

import (
	"SiteForDsBot/conf"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// кодирует переданный пароль в sha256
func EncodePassword(password string) string {
	hasher := sha256.New()
	data := password + conf.Salt
	hasher.Write([]byte(data))
	hash := hasher.Sum(nil)
	return hex.EncodeToString(hash)
}


func GenerateJWT(uuid string) (string, error) {
    claims := jwt.MapClaims{
        "uuid": uuid,
        "exp":  time.Now().Add(time.Hour * 100).Unix(),
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(conf.Jwt_secret))
}


func GetUserUuidFromJWT(tokenString string) (string, error) {
  token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
  if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
    return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
  }
    return []byte(conf.Jwt_secret), nil
  })
  if err != nil {
    return "", err
  }

  if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
    uuid, ok := claims["uuid"].(string)
    if !ok {
      return "", fmt.Errorf("uuid not found or invalid type")
    }
    return uuid, nil
  }
  return "", fmt.Errorf("invalid token")
}

