package main

import (
	// "encoding/json" //encodes json
	// "fmt"
	"flag"
	"log"
	"net/http" //for using http

	//local packages
	"go-restful/app/auth"
	"go-restful/app/websockets"

	//external libraries must be donwloaded using go get command first
	"github.com/gorilla/mux" //this external library allows you to specify http method
	// "github.com/gorilla/websocket" //this external library used for websockets
)

var addr = flag.String("addr", ":8080", "http service address")

//callback function for home route
func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.URL.Path != "/" { //checks the path
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != "GET" { //checks the http method
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "../public/index.html")
}

//callback function for home route
func serveChat(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.URL.Path != "/chat" { //checks the path
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != "GET" { //checks the http method
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "../public/chat.html")
}

func main() {
	flag.Parse()

	myRouter := mux.NewRouter().StrictSlash(true)

	myRouter.HandleFunc("/", serveHome)
	myRouter.HandleFunc("/chat", serveChat)
	myRouter.HandleFunc("/ws/{id}", func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)                  //get the id parameter
		websockets.ServeWs(params["id"], w, r) //defined in client.go
	})

	err := http.ListenAndServe(*addr, myRouter)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

	mongoErr := auth.InitMongoClient()
	if mongoErr != nil {
		log.Fatal("Mongo client error: ", err)
	}
}
