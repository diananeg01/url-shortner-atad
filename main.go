package main

import (
	"fmt"
	"net/http"

	"github.com/diananeg01/url-shortner-atad/database"
	"github.com/diananeg01/url-shortner-atad/frontend"
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

	http.HandleFunc("/", frontend.UrlShortner)
	http.HandleFunc("/line", frontend.TestChartRender)

	println("Server running on http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		return
	}
}
