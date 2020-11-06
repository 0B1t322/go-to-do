package firstservice

import (
	"errors"
	"database/sql"
	log "github.com/sirupsen/logrus"
	"net/http"
	"encoding/json"
)

var ( 
	// ErrNameTaken means thats entered name is exsits in bd
	ErrNameTaken = errors.New("Name is already taken")

	// ErrNameFieldEmpty means that you dont entered name
	ErrNameFieldEmpty = errors.New("A name field is empty")
	// ErrPassFieldEmpty meant that you dont entered password
	ErrPassFieldEmpty = errors.New("A password field is empty")
)

// User is a struct of user Model
type User struct {
	ID int `json:"id"`
	Name string `json:"name"`
	Password string `json:"password"`
	Token interface{} `json:"token"`
}

// NewUser create a new User
func NewUser(name, password string) *User {
	return &User{Name: name, Password: password}
}

// UpdateUserToken put new token
func UpdateUserToken(id int, token string) {

}

// GetUserByName is return a User model and error of Scanning row
func GetUserByName(name string) (*User, error) {
	DB, err := ConnectToDB()
	if err != nil {
		log.Warnln(err)
	}
	defer DB.Close()

	row := DB.QueryRow("select * from users where name = $1", name)
	getUser := &User{}
	err = row.Scan(&getUser.ID, &getUser.Name, &getUser.Password, &getUser.Token)
	return getUser, err
}

func GetUserByToken(token string) (*User, error) {
	DB, err := ConnectToDB()
	if err != nil {
		log.Warn(err)
		return nil, err
	}
	defer DB.Close()
	
	row := DB.QueryRow("select * from users where token = $1", token)
	getUser := &User{}

	err = row.Scan(&getUser.ID, &getUser.Name, &getUser.Password, &getUser.Token)
	return getUser, err
}

func ParseJSONUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	decoder := json.NewDecoder(r.Body) // создаем новый декодер
	var user User
	err := decoder.Decode(&user) // декодируем тело запроса в объет user

	if user.Name == "" {
		b, _ := NewError(ErrNameFieldEmpty).Marshall()
		w.Write(b)
		return
	}

	if user.Password == "" {
		b, _ := NewError(ErrPassFieldEmpty).Marshall()
		w.Write(b)
		return
	}

	if err != nil {
		log.Warn(err)
	}

	DB, err := ConnectToDB()
	if err != nil {
		log.Warnln(err)
	}
	defer DB.Close()

	getUser, err := GetUserByName(user.Name)
	if err == sql.ErrNoRows { // если поле не было найдено создаем его
		DB.Exec("insert into users (name, password) values ($1, $2)", user.Name, user.Password)
		b, _ := NewError(errors.New("")).Marshall()
		w.Write(b)

	} else if err != nil { 
		log.Warnln(err)
	}

	if getUser.Name == user.Name { // если такое имя уже существует
		b, err := NewError(ErrNameTaken).Marshall()
		if err != nil {
			log.Warn()
		}
		log.Info(ErrNameTaken)
		w.Write(b)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}


