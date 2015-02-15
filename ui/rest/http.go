package rest

import (
	//log "github.com/golang/glog"
	"fmt"
	"github.com/gorilla/mux"
	//"github.com/silenteh/gantryos/models"
	"net/http"
)

var router *mux.Router

func InitWebInterface(port string) {

	router = mux.NewRouter()
	router.StrictSlash(true)
	router.HandleFunc("/task", taskHandler)
	router.HandleFunc("/kill/{taskid}", killHandler)
	router.HandleFunc("/undeploy/{appname}/{appver}", undeployHandler)
	router.HandleFunc("/scale/{appname}/{appver}/{instances}", scaleHandler)
	router.HandleFunc("/deploy", deployHandler)
	router.HandleFunc("/ping", pingHandler)

	http.Handle("/", router)
	http.ListenAndServe(":"+port, nil)
}

func taskHandler(w http.ResponseWriter, r *http.Request) {

}

func killHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Kill task handler")
}

func deployHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "App deployment handler")
}

func scaleHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "App scaling handler")
}

func undeployHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "App undeployment handler")
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Connection", "close")
	fmt.Fprintf(w, "pong")
}
