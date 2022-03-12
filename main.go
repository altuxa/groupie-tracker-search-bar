package main

import (
	"fmt"
	handlers "groupie-tracker/cmd/handlers"
	"log"
	"net/http"
)

func main() {
	// url := "https://groupietrackers.herokuapp.com/api/artists"
	// Parse(url, &B.Artists)
	// Parse("https://groupietrackers.herokuapp.com/api/relation", &B.Relation)
	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))
	mux.HandleFunc("/", handlers.Home)
	mux.HandleFunc("/artist/", handlers.Artist)
	mux.HandleFunc("/search", handlers.Search)
	fmt.Println("starting server at localhost:8080")
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatalln(err)
		return
	}
}
