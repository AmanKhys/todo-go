package main

import (
	"encoding/json"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"strconv"
)

const (
	dbDriver = "mysql"
	dbUser   = "root"
	dbPass   = "pass"
	dbName   = "todo_app"
)

var db, _ = gorm.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName+"?charset=utf8&parseTime=True&loc=Local")

type TodoItemModel struct {
	Id          int `gorm:"primary key"`
	Description string
	Completed   bool
}

func Healthz(w http.ResponseWriter, r *http.Request) {
	log.Info("API health is okay")
	w.Header().Set("Content-Type", "application/json")
	io.WriteString(w, `{ "alive": true }`)
}

func init() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetReportCaller(true)
}

func CreateItem(w http.ResponseWriter, r *http.Request) {
	description := r.FormValue("description")
	log.WithFields(log.Fields{"description": description}).Info("Add new Todo Item. Saving to the database. ")
	todo := &TodoItemModel{Description: description, Completed: false}
	db.Create(&todo)
	result := db.Last(&todo)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&result.Value)
}

func main() {
	defer db.Close()

	db.Debug().DropTableIfExists(&TodoItemModel{})
	db.Debug().AutoMigrate(&TodoItemModel{})

	log.Info("Starting todo list API")
	router := mux.NewRouter()
	router.HandleFunc("/Healthz", Healthz).Methods("GET")
	router.HandleFunc("/todo", CreateItem).Methods("POST")
	http.ListenAndServe(":8000", router)
}
