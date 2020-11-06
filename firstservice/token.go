package firstservice

import (
	"fmt"
	"encoding/json"
	jwt "github.com/dgrijalva/jwt-go"
	log "github.com/sirupsen/logrus"
)

type Token struct {
	Token string `json:"token"`
}

func NewToken(token string) *Token {
	return &Token{Token: token}
}

func (t *Token) Marahsll() ([]byte, error) {
	return json.Marshal(t)
}

func (t *Token) parseToken(userName string) (int, error) {
	token, err := jwt.Parse(t.Token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			log.Warn("Not okay")
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(userName), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if !ok {
			log.Warn("Not okay, but id is ", claims["id"]," Valid:", token.Valid)
			return -1, err
		}
		return int (claims["id"].(float64)), err
	} else {
		log.Warn(err)
		return -1, err
	}
}
