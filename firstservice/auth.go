package firstservice

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"encoding/json"
	"net/http"
	jwt "github.com/dgrijalva/jwt-go"
)

var (
	ErrInvalidNameOrPassword = errors.New("Invalid name or password")
)

func GetTokenHandler(w http.ResponseWriter, r *http.Request) {
	
	decoder := json.NewDecoder(r.Body) // создаем новый декодер
	var user User
	err := decoder.Decode(&user) // декодируем тело запроса в объет user
	if err != nil {
		log.Warn(err)
		writeError(err,w)
	}
	
	getUser, err := GetUserByName(user.Name)
	if err !=  nil {
		log.Warn(err)
		writeError(err,w)
		return
	}

	if getUser.Name == user.Name && getUser.Password == user.Password {
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"id": getUser.ID,
		})

		claim, _ := token.Claims.(jwt.MapClaims)
		log.Info(claim["id"])

		tokenString, err := token.SignedString( []byte( getUser.Name  ) ) // для простоты сделаем так
		if err != nil{
			log.Warn(err)
			writeError(err,w)
			return
		}

		DB, err := ConnectToDB()
		if err != nil {
			log.Warn(err)
			writeError(err,w)
			return
		}
		defer DB.Close()

		_, err = DB.Exec("update users set token = $1 WHERE _id = $2", tokenString, getUser.ID)
		if err != nil {
			writeError(err,w)
			return
		}

		b, _ := NewToken(tokenString).Marahsll()
		w.Write(b)
		return

	} else {
		b, _ := NewError(ErrInvalidNameOrPassword).Marshall()
		w.Write(b)
		return
	}

}
