package frontend

import (
	"fmt"
	"net/http"

	"github.com/diananeg01/url-shortner-atad/database"
	"golang.org/x/crypto/bcrypt"
	g "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		Page("Login",
			Div(
				H1(Class("text-3xl font-bold mb-4"), g.Text("Log in")),
				Form(
					Method("POST"),
					Action("/login"),

					Label(For("email"), Class("block mb-1 font-semibold"), g.Text("Email:")),
					Input(Type("email"), Name("email"), Class("border rounded w-full p-2 mb-4"), ID("email")), Br(),

					Label(For("password"), Class("block mb-1 font-semibold"), g.Text("Password:")),
					Input(Type("password"), Name("password"), Class("border rounded w-full p-2 mb-4"), ID("password")), Br(),

					Button(Type("submit"), Class("bg-blue-600 text-white px-4 py-2 rounded hover:bg-blue-700"), g.Text("Login")),
				),
			),
		).Render(w)
		return
	}

	// POST
	r.ParseForm()
	email := r.FormValue("email")
	password := r.FormValue("password")

	var userID string
	var hash string
	conn := database.GetDB()
	err := conn.QueryRow(`SELECT user_id, password_hash FROM user_data WHERE email = $1`, email).
		Scan(&userID, &hash)

	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Check password
	if bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Set session cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session_user",
		Value:    fmt.Sprint(userID),
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // set true in production
	})

	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}
