package main

import (
	p "github.com/diananeg01/url-shortner-atad/frontend"
	"net/http"
)

func main() {
	http.HandleFunc("/", p.TestTitle)
	http.HandleFunc("/line", p.Httpserver)

	println("Server running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
