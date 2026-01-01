package frontend

import (
	"net/http"

	"github.com/diananeg01/url-shortner-atad/database"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	g "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func Register(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		Page("Register",
			Div(Class("h-dvh bg-gray-900"),
				Div(Class("flex min-h-full flex-col justify-center px-6 py-12 lg:px-8"),
					Div(Class("sm:mx-auto sm:w-full sm:max-w-sm"),
						Img(Src("https://tailwindcss.com/plus-assets/img/logos/mark.svg?color=indigo&shade=500"), Class("mx-auto h-10 w-auto")),
						H2(Class("mt-10 text-center text-2xl/9 font-bold tracking-tight text-white"), g.Text("Register")),
					),
					Div(Class("mt-10 sm:mx-auto sm:w-full sm:max-w-sm"),
						Form(Method("POST"),
							Action("/register"),
							Div(
								Label(For("email"), Class("block text-sm/6 font-medium text-gray-100"), g.Text("Email:")),
								Div(Class("mt-2"), Input(Type("email"), Name("email"), Class("border rounded w-full p-2 mb-4"), ID("email"))),
							),
							Div(
								Label(For("password"), Class("block text-sm/6 font-medium text-gray-100"), g.Text("Password:")),
								Div(Class("mt-2"), Input(Type("password"), Name("password"), Class("border rounded w-full p-2 mb-4"), ID("password"))),
							),

							Div(
								Button(Type("submit"), Class("flex w-full justify-center rounded-md bg-indigo-500 px-3 py-1.5 text-sm/6 font-semibold text-white hover:bg-indigo-400 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-500"), g.Text("Register")),
							),
						),

						P(Class("mt-10 text-center text-sm/6 text-gray-400"), g.Text("Already have an account? "), A(Href("/login"), Class("font-semibold text-indigo-400 hover:text-indigo-300"), g.Text("Log in"))),
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

	if email == "" || password == "" {
		http.Error(w, "Missing fields", http.StatusBadRequest)
		return
	}

	// Hash password
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	newUuid := uuid.New()
	conn := database.GetDB()
	_, err := conn.Exec(`
        INSERT INTO user_data (user_id, email, password_hash, crea, lupa, status)
        VALUES ($1, $2, $3, NOW(), NOW(), 'active')
    `, newUuid, email, string(hash))

	if err != nil {
		http.Error(w, "There was an error when creating the user: "+err.Error(), http.StatusConflict)
		return
	}

	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}
