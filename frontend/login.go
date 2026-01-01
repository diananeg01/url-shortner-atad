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
			Div(Class("h-dvh bg-gray-900"),
				Div(Class("flex min-h-full flex-col justify-center px-6 py-12 lg:px-8"),
					Div(Class("sm:mx-auto sm:w-full sm:max-w-sm"),
						Img(Src("https://tailwindcss.com/plus-assets/img/logos/mark.svg?color=indigo&shade=500"), Class("mx-auto h-10 w-auto")),
						H2(Class("mt-10 text-center text-2xl/9 font-bold tracking-tight text-white"), g.Text("Log in to your account")),
					),
					Div(Class("mt-10 sm:mx-auto sm:w-full sm:max-w-sm"),
						Form(Method("POST"),
							Action("/login"),
							Div(
								Label(For("email"), Class("block text-sm/6 font-medium text-gray-100"), g.Text("Email:")),
								Div(Class("mt-2"), Input(Type("email"), Name("email"), Class("border rounded w-full p-2 mb-4"), ID("email"))),
							),
							Div(
								Label(For("password"), Class("block text-sm/6 font-medium text-gray-100"), g.Text("Password:")),
								Div(Class("mt-2"), Input(Type("password"), Name("password"), Class("border rounded w-full p-2 mb-4"), ID("password"))),
							),

							Div(
								Button(Type("submit"), Class("flex w-full justify-center rounded-md bg-indigo-500 px-3 py-1.5 text-sm/6 font-semibold text-white hover:bg-indigo-400 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-500"), g.Text("Login")),
							),
						),

						P(Class("mt-10 text-center text-sm/6 text-gray-400"), g.Text("If you don't have an account: "), A(Href("/register"), Class("font-semibold text-indigo-400 hover:text-indigo-300"), g.Text("Register"))),
					),
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
	err := conn.QueryRow(`SELECT email, password_hash FROM user_data WHERE email = $1`, email).
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
		MaxAge:   3600, // seconds (1 hour)
		HttpOnly: true,
		Secure:   false, // set true in production
	})

	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}
