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

	errQuery := conn.QueryRow("SELECT 1").Scan(&variable)
	if errQuery != nil {
		return
	}
	fmt.Println(variable)

	http.HandleFunc("/", p.TestTitle)
	http.HandleFunc("/line", p.Httpserver)

	println("Server running on http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		return
	}
}
