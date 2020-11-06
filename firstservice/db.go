package firstservice

import (
	"io/ioutil"
	"os"
	_ "github.com/mattn/go-sqlite3"
	"database/sql"
	log "github.com/sirupsen/logrus"
)

var pathToDB string
// DB is a pointer to DB
// var DB *sql.DB

func init() {
	
	data, err := ioutil.ReadFile("./bd/firstservice.sql")
	if err != nil {
		log.Fatal(err)
	}
	

	var ( 
		DBdir = "bd"
		DBfile = "firstservice.db"
		sqlScript = string(data)
	)

	pathToDB = "./"+DBdir+"/"+DBfile
	
	if CheckFile(".", DBdir, true) { // if bd dir is not exsist
		os.Mkdir(DBdir, 0750)
		if CheckFile(DBdir, DBfile, false) {
			// create bd file
			file, err := os.Create(pathToDB)
			if err != nil {
				log.Warn(err)
			}
			file.Close()

			// open DB file
			db, err := sql.Open("sqlite3", pathToDB)
			if err != nil {
				log.Warn(err)
			}

			// Exec command to create table
			_, err = db.Exec(sqlScript)
			if err != nil {
				log.Warn(err)
			}
			db.Close()
		}
	}

	db, err := sql.Open("sqlite3", pathToDB)
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()
	
	_, err = db.Exec(sqlScript)
	if err != nil {
		log.Fatalln(err)
	}
	
}

// ConnectToDB is func that connect to database and return DB pointer and err
func ConnectToDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", pathToDB)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return db, nil
}
// CheckFile check is a current file is exsists
func CheckFile(rootDir ,fileName string, isDir bool) bool {
	files, err := ioutil.ReadDir(rootDir)
	if err != nil {
		log.Warn(err)
	}

	for _, file := range files {
		if file.Name() == fileName && file.IsDir() == isDir {
			return false // yes its exsist 
		}
	}

	log.Warn(fileName, " is not exsist")
	return true // db is not exsist
}
