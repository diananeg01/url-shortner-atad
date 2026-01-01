package main

import (
	"fmt"
	"net/http"
	"time"

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
	http.HandleFunc("/register", frontend.Register)
	http.HandleFunc("/login", frontend.Login)
	http.HandleFunc("/line", frontend.TestChartRender)
	http.HandleFunc("/redirect/{url}", frontend.RedirectURL)
	http.HandleFunc("/logout", handleLogout)

	println("Server running on http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		return
	}
}

func handleLogout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "session_user",
		Value:    "",
		Expires:  time.Unix(0, 0),
		Path:     "/",
		HttpOnly: true,
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
