package main

import (
	"fmt"
	"github.com/diananeg01/url-shortner-atad/database"
	p "github.com/diananeg01/url-shortner-atad/frontend"
	"net/http"
)

func main() {
	database.Init()
	defer database.Close()

	conn := database.GetDB()
	var variable any

	conn.QueryRow("SELECT id from test").Scan(&variable)
	fmt.Println(variable)

	http.HandleFunc("/", p.TestTitle)
	http.HandleFunc("/line", p.Httpserver)

	println("Server running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
