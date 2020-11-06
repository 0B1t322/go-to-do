package firstservice

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

var (
	ErrDescriptionFieldEmpty = errors.New("Description field empty")
	ErrNoTasks = errors.New("You dont have tasks")
)
// Task is a struct of Model task
type Task struct {
	USERID	int			`json:"-"`
	Name string			`json:"name"`
	Description string	`json:"description"`
	Date time.Time		 `json:"-"`
	DateStr string 		`json:"date"`
	Done bool			`json:"done"`
}

type taskWithToken struct {
	Task
	Token
}

// Tasks for all task for current user
type Tasks struct {
	Tasks []Task
}

func (t *Task) strsDate() {
	t.DateStr = t.Date.Format(time.Stamp)
}

// Marshall for taking a json
func (t *Tasks) Marshall() ([]byte, error){
	return json.Marshal(t)
}

// NewTask create a new task and return pointer
func NewTask(Name, Decription string) *Task {
	dateNow := time.Now()
	formatDate, _ := time.Parse("01-02-2006 15:04:05 Mon", dateNow.Format("01-02-2006 15:04:05 Mon"))
	return &Task{Name: Name, Description: Decription, Date: formatDate, Done: false}
}

// GetTaskByName return a task model and error of Scanning row and task id
func GetTaskByName(name string) (*Task, error, int) {
	// TODO поменять чтобы не было ошибки
	DB, err := ConnectToDB()
	if err != nil {
		log.Fatal(err)
	}
	defer DB.Close()

	row := DB.QueryRow("select * from tasks where name = $1", name)
	getTask := &Task{}
	var id int
	var dateStr string
	err = row.Scan(&id, &getTask.USERID ,&getTask.Name, &getTask.Description, &dateStr, &getTask.Done)
	date, _ := time.Parse(time.Stamp, dateStr)
	getTask.Date = date
	return getTask, err, id
}


// POSTTaskHandler Post TASK to current user
func POSTTaskHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	decoder := json.NewDecoder(r.Body) // новый декодер
	var taskWT taskWithToken
	err := decoder.Decode(&taskWT)
	if err != nil {
		log.Warnln(err)
		b, _ := NewError(err).Marshall()
		w.Write(b)
		return
	}

	token 	:= &taskWT.Token
	task 	:= &taskWT.Task

	task.Done = false
	task.Date = time.Now()
	task.strsDate()

	user, err := GetUserByToken(token.Token) 
	if err == sql.ErrNoRows {
		b, _ := NewError(errors.New("Unexpected token")).Marshall()
		w.Write(b)
		return
	} else if err != nil {
		log.Warn(err)
		b, _ := NewError(err).Marshall()
		w.Write(b)
		return
	}
	
	id, err := token.parseToken(user.Name)
	if err != nil {
		log.Warn(err)
		b, _ := NewError(err).Marshall()
		w.Write(b)
		return
	}
	
	task.USERID = id

	if task.Name == "" {
		b, _ := NewError(ErrNameFieldEmpty).Marshall()
		w.Write(b)
		return
	}

	if task.Description == "" {
		b, _ := NewError(ErrDescriptionFieldEmpty).Marshall()
		w.Write(b)
		return
	}

	DB, err := ConnectToDB()
	if err != nil {
		log.Warnln(err)
	}
	defer DB.Close()

	getTask, err, _ := GetTaskByName(task.Name)
	if err == sql.ErrNoRows {
		_, err := DB.Exec(
					"insert into tasks (user_id, name, decription, DATE, DONE) values ($1, $2, $3, $4, $5)", 
					task.USERID, task.Name, task.Description, task.Date.Format(time.Stamp),
					task.Done, 
		)
		if err != nil {
			b, _ := NewError(err).Marshall()
			w.Write(b)
			return
		}
		b, _ := NewError(errors.New("")).Marshall()
		w.Write(b)

	} else if err != nil {
		log.Warn(err)
	}

	if getTask.Name == task.Name {
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

// PUTTaskHandler update tasks
func PUTTaskHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")


	decoder := json.NewDecoder(r.Body) // новый декодер
	var taskWT taskWithToken
	err := decoder.Decode(&taskWT)
	if err != nil {
		log.Warnln(err)
		writeError(err, w)
		return
	}

	token 	:= &taskWT.Token
	task 	:= &taskWT.Task

	user, err := GetUserByToken(token.Token) 
	if err == sql.ErrNoRows {
		writeError(errors.New("Unexpected token"), w)
		return
	} else if err != nil {
		log.Warn(err)
		writeError(err, w)
		return
	}

	id, err := token.parseToken(user.Name)
	if err != nil {
		log.Warn(err)
		writeError(err, w)
		return
	}
	
	task.USERID = id

	if task.Name == "" {
		writeError(ErrNameFieldEmpty, w)
		return
	}

	if task.Description == "" {
		writeError(ErrDescriptionFieldEmpty, w)
		return
	}


	DB, err := ConnectToDB()
	if err != nil {
		log.Warnln(err)
	}
	defer DB.Close()


	_, err, taskID := GetTaskByName(task.Name)
	if err != nil {
		writeError(err,w)
		return
	}
	log.Info(task.Done)
	
	quary := `update tasks set name = $1, decription = $2, DONE = $3 
	where _id = $4`
	_, err = DB.Exec(quary, task.Name, task.Description, task.Done, taskID)
	if err != nil {
		log.Warn(err)
		writeError(err, w)
		return
	}

	writeError(errors.New(""), w)

}

// DELETETask delete a task
func DELETETask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	decoder := json.NewDecoder(r.Body) // новый декодер
	var taskWT taskWithToken
	err := decoder.Decode(&taskWT)
	if err != nil {
		log.Warnln(err)
		b, _ := NewError(err).Marshall()
		w.Write(b)
		return
	}

	token 	:= &taskWT.Token
	task 	:= &taskWT.Task

	_ , err = GetUserByToken(token.Token) 
	if err == sql.ErrNoRows {
		b, _ := NewError(errors.New("Unexpected token")).Marshall()
		w.Write(b)
		return
	} else if err != nil {
		log.Warn(err)
		b, _ := NewError(err).Marshall()
		w.Write(b)
		return
	}

	_ , err, id := GetTaskByName(task.Name)
	if err != nil {
		log.Warn(err)
		writeError(err, w)
		return
	}

	DB, err := ConnectToDB()
	if err != nil {
		log.Warnln(err)
	}
	defer DB.Close()

	_, err = DB.Exec(
		"delete from tasks where _id = $1 and name = $2",
		id, task.Name,
		)
	if err != nil {
		writeError(err, w)
		log.Warn(err)
		return
	}

	writeError(errors.New(""), w) 

}

// GETTasks is a handler for get  all task by a user token
// for input expect a token json
func GETTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var (
		token Token
		userID int
		tasks Tasks
	)

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&token); err != nil {
		log.Warn(err)
		writeError(err,w)
		return
	}

	if getUser, err := GetUserByToken(token.Token); err != nil {
		log.Warn(err)
		writeError(err,w)
		return

	} else {

		userID, err = token.parseToken(getUser.Name)
		if err != nil {
			log.Warn(err)
			writeError(err, w)
			return
		}

	}

	DB, err := ConnectToDB()
	if err != nil {
		log.Warn(err)
		writeError(err, w)
		return
	}
	defer DB.Close()

	rows, err := DB.Query("select * from tasks where user_id = $1", userID)
	if err == sql.ErrNoRows {
		log.Warn(err)
		writeError(ErrNoTasks, w)
	} else if err != nil {
		log.Warn(err)
		writeError(err, w)
	}
	defer rows.Close()

	for rows.Next() {
		var task Task
		var id int
		var dateStr string
		if err := rows.Scan(
			&id, &task.USERID, &task.Name, &task.Description, &dateStr, &task.Done,
			); err != nil {
				log.Warn(err)
				writeError(err, w)
		}
		task.Date, _ = time.Parse(time.Stamp, dateStr)
		task.strsDate()
		tasks.Tasks = append(tasks.Tasks, task)
	}

	if b, err := tasks.Marshall(); err != nil {
		log.Warn(err)
		writeError(err, w)

	} else {
		w.Write(b)
	}




}