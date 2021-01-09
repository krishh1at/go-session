package controllers

import (
	"crypto/hmac"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/krishh1at/test/app/models"
	"github.com/krishh1at/test/config"
	"golang.org/x/crypto/bcrypt"
)

var Key = []byte("Hello World!")

func Login(w http.ResponseWriter, r *http.Request) {
	msg := r.FormValue("msg")
	config.Template.ExecuteTemplate(w, "login.html", msg)
}

func SignUp(w http.ResponseWriter, r *http.Request) {
	msg := r.FormValue("msg")
	config.Template.ExecuteTemplate(w, "signup.html", msg)
}

func CreateAccount(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		msg := url.QueryEscape("Please use post method to send data.")
		http.Redirect(w, r, "/signup?msg="+msg, http.StatusSeeOther)
		return
	}

	firstName := r.FormValue("firstName")
	lastName := r.FormValue("lastName")
	email := r.FormValue("email")
	password := r.FormValue("password")
	confirmPassword := r.FormValue("confirmPassword")

	if !strings.EqualFold(password, confirmPassword) {
		msg := url.QueryEscape("Password and Confirm Password didn't match.")
		http.Redirect(w, r, "/signup?msg="+msg, http.StatusSeeOther)
		return
	}

	db, err := sql.Open("mysql", "krishna:Krishna@501@/test?charset=utf8")

	if err != nil {
		fmt.Println("Connection not stablished!")
	}

	if err = db.Ping(); err != nil {
		fmt.Println("Connection not stablished!")
	}

	stmt, err := db.Prepare("SELECT id FROM users ORDER BY ID DESC LIMIT 1")
	if err != nil {
		log.Fatalln("Unable to prepare db select query")
	}

	defer stmt.Close()
	id, err := stmt.Exec()
	if err != nil {
		log.Fatalln("Unable to process db select query")
	}

	stmt1, err := db.Prepare("INSERT INTO users VALUES(?, ?, ?, ?, ?)")
	if err != nil {
		log.Fatalln("Unable to prepare insert db query")
	}

	iD, _ := id.LastInsertId()

	defer stmt1.Close()
	_, err = stmt1.Exec(iD, firstName, lastName, email, EncryptPassword(password))
	if err != nil {
		log.Fatalln("Unable to process insert db query")
	}

	http.Redirect(w, r, "/login", http.StatusPermanentRedirect)
}

func LoggedIn(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		msg := url.QueryEscape("Please use post method to send data.")
		http.Redirect(w, r, "/signup?msg="+msg, http.StatusSeeOther)
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")

	if !ComparePassword(email, password) {
		msg := url.QueryEscape("Password and Confirm Password didn't match.")
		http.Redirect(w, r, "/login?msg="+msg, http.StatusSeeOther)
		return
	}

	token := CreateToken(email)
	c := http.Cookie{
		Name:  "sessionId",
		Value: token,
	}

	http.SetCookie(w, &c)
	http.Redirect(w, r, "/", http.StatusPermanentRedirect)
}

func EncryptPassword(password string) string {
	encryptedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalln("Getting error while encrypting password!")
	}

	return base64.StdEncoding.EncodeToString(encryptedPass)
}

func ComparePassword(email, password string) bool {
	DB, err := sql.Open("mysql", "krishna:Krishna@501@/test?charset=utf8")

	if err != nil {
		fmt.Println("Connection not stablished!")
	}

	if err = DB.Ping(); err != nil {
		fmt.Println("Connection not stablished!")
	}

	user := models.User{}

	err = DB.QueryRow("SELECT password FROM users WHERE email = ?", email).Scan(&user.Password)
	if err != nil {
		log.Fatalf("Unable to prepare select password statement: %w", err)
	}

	pass, err := base64.StdEncoding.DecodeString(user.Password)
	if err != nil {
		log.Fatalln("Getting Error while decoding password")
	}

	err = bcrypt.CompareHashAndPassword(pass, []byte(password))
	if err != nil {
		log.Fatalln("Password don't match")
		return false
	}

	return true
}

func CreateToken(sid string) string {
	mac := hmac.New(sha256.New, Key)
	mac.Write([]byte(sid))
	signedMac := base64.StdEncoding.EncodeToString(mac.Sum(nil))

	return signedMac + "|" + sid
}
