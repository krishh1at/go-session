package main

import (
	"net/http"

	"github.com/krishh1at/test/app/controllers"
	"github.com/krishh1at/test/db"
)

func main() {
	db.Open()
	// db.CreateTable()
	http.HandleFunc("/", controllers.Home)
	http.HandleFunc("/login", controllers.Login)
	http.HandleFunc("/signup", controllers.SignUp)
	http.HandleFunc("/signedup", controllers.CreateAccount)
	http.HandleFunc("/loggedin", controllers.LoggedIn)
	http.ListenAndServe(":8080", nil)
	defer db.DB.Close()
}
