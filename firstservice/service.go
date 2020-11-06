package firstservice

import (
	"errors"
	"net/http"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

var port = ":8080"


// var activeTokens mapPointer

func init() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetReportCaller(true)
}

func writeError(err error, w http.ResponseWriter) {
	b, _ := NewError(err).Marshall()
	w.Write(b)
}

// Start is a staring first service
func Start() {
	r := mux.NewRouter()
	r.HandleFunc("/{name}",homeHandler)
	r.HandleFunc("/api/create-user", ParseJSONUser).Methods("POST")
	r.HandleFunc("/api/create-task", POSTTaskHandler).Methods("POST")
	r.HandleFunc("/api/get-token", GetTokenHandler).Methods("GET")
	r.HandleFunc("/api/update-task", PUTTaskHandler).Methods("PUT")
	r.HandleFunc("/api/delete-task", DELETETask).Methods("DELETE")
	r.HandleFunc("/api/get-tasks", GETTasks).Methods("GET")


	log.Infoln("Starting server on", port)
	if err := http.ListenAndServe(port,r); err != nil {
		log.Fatalln(err)
	}
	
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	writeError(errors.New( name + " is uknown, use /api"), w)
}


