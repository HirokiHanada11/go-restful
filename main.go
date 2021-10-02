package main

import (
	"encoding/json" //encodes json
	"fmt"
	"log"
	"net/http" //for using http

	//external libraries must be donwloaded using go get command first
	"github.com/gorilla/mux" //this external library allows you to specify http method
)

type Article struct {
	Title   string `json:"Title"`
	Desc    string `json:"desc"`
	Content string `json:"content"`
}

type Articles []Article

//call back function for articles route with GET method
func allArticles(w http.ResponseWriter, r *http.Request) {
	articles := Articles{
		Article{Title: "Test Title", Desc: "For Test", Content: "Hello world"},
	}

	fmt.Println("Endpoint hit: All Articles endpoint")
	json.NewEncoder(w).Encode(articles)
}

//call back function for articles route with POST method
func testPostArticles(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Test post endpoint")
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Homepage Endpoint Hit")
}

//function for specifying routes
func handleRequests() {

	myRouter := mux.NewRouter().StrictSlash(true)

	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/articles", allArticles).Methods("GET")
	myRouter.HandleFunc("/articles", testPostArticles).Methods("POST")
	log.Fatal(http.ListenAndServe(":8081", myRouter))
}

func main() {
	handleRequests()
}
