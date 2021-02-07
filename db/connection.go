package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func Open() *sql.DB {
	DB, err := sql.Open("mysql", "krishna:Krishna@501@/test?charset=utf8")
	if err != nil {
		log.Fatalln("Unnable to open database:" + err.Error())
	}

	if err = DB.Ping(); err == nil {
		log.Println("Database connection stablished...")
	}

	return DB
}

func CreateTable() sql.Result {
	DB, err := sql.Open("mysql", "krishna:Krishna@501@/test?charset=utf8")

	if err != nil {
		fmt.Println("Connection not stablished!")
	}

	if err = DB.Ping(); err != nil {
		fmt.Println("Connection not stablished!")
	}

	stmt, err := DB.Prepare("CREATE TABLE IF NOT EXISTS `test`.`users` (id int, first_name varchar(255), last_name varchar(255), email varchar(255), password varchar(255))")
	if err != nil {
		log.Fatalf("Unable to prepare create table statement: %w", err)
	}

	defer stmt.Close()

	rslt, err := stmt.Exec()
	if err != nil {
		log.Fatalf("Unnable to create table: %w", err)
	}

	return rslt
}
