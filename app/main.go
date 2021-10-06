package main

import (
	// "encoding/json" //encodes json
	// "fmt"
	"flag"
	"log"
	"net/http" //for using http

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

func main() {
	flag.Parse()

	hub := newHub() //creates a new hub struct from hub.go file

	go hub.run() //calls run method defined in hub.go file

	myRouter := mux.NewRouter().StrictSlash(true)

	myRouter.HandleFunc("/", serveHome)
	myRouter.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r) //defined in client.go
	})
	err := http.ListenAndServe(*addr, myRouter)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

// type Article struct {
// 	Title   string `json:"Title"`
// 	Desc    string `json:"desc"`
// 	Content string `json:"content"`
// }

// type Articles []Article

// var upgrader = websocket.Upgrader{
// 	ReadBufferSize:  1024,
// 	WriteBufferSize: 1024,
// 	CheckOrigin:     func(r *http.Request) bool { return true },
// }

// //call back function for articles route with GET method
// func allArticles(w http.ResponseWriter, r *http.Request) {
// 	articles := Articles{
// 		Article{Title: "Test Title", Desc: "For Test", Content: "Hello world"},
// 	}

// 	fmt.Println("Endpoint hit: All Articles endpoint")
// 	json.NewEncoder(w).Encode(articles)
// }

// //call back function for articles route with POST method
// func testPostArticles(w http.ResponseWriter, r *http.Request) {
// 	fmt.Fprintf(w, "Test post endpoint")
// }

// func homePage(w http.ResponseWriter, r *http.Request) {
// 	fmt.Fprintf(w, "Homepage Endpoint Hit")
// }

// func reader(conn *websocket.Conn) {
// 	for {
// 		messageType, p, err := conn.ReadMessage()
// 		if err != nil {
// 			log.Println(err)
// 			return
// 		}

// 		log.Println(string(p))

// 		if err := conn.WriteMessage(messageType, p); err != nil {
// 			log.Println(err)
// 			return
// 		}
// 	}
// }

// func wsEndopoint(w http.ResponseWriter, r *http.Request) {
// 	ws, err := upgrader.Upgrade(w, r, nil)
// 	if err != nil {
// 		log.Println(err)
// 	}

// 	log.Println("Client Successfully connected...")

// 	reader(ws)
// }

// //function for specifying routes
// func handleRequests() {

// 	myRouter := mux.NewRouter().StrictSlash(true)

// 	myRouter.HandleFunc("/", homePage)
// 	myRouter.HandleFunc("/articles", allArticles).Methods("GET")
// 	myRouter.HandleFunc("/articles", testPostArticles).Methods("POST")
// 	myRouter.HandleFunc("/ws", wsEndopoint)
// 	log.Fatal(http.ListenAndServe(":8081", myRouter))
// }

// func main() {
// 	fmt.Println("HTTP server started nice")
// 	handleRequests()
// }

// //run docker run -it -p 8080:8081 go-webapi
