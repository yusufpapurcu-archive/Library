package main

import (
	"net/http"

	"github.com/yusufpapurcu/Library/database"
	"github.com/yusufpapurcu/Library/route"
)

func main() {
	database.Connect("mongodb://localhost:27017") //MongoDB Connection Started
	router := route.SetRouter()                   // Api Router's Set
	publicFiles := http.FileServer(http.Dir("pages/public"))
	router.PathPrefix("/public/").Handler(http.StripPrefix("/public/", publicFiles))
	http.Handle("/", router) // Start Server
	http.ListenAndServe(":8080", nil)
}
