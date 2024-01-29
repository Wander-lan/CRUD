package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	connectionStr := "golang:golang@/my_database?charset=utf8&parseTime=True&loc=Local"
	db, err := sql.Open("mysql", connectionStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}

	query, err := db.Query("select * from users")
	if err != nil {
		log.Fatal(err)
	}
	defer query.Close()

	fmt.Println(query)

}
