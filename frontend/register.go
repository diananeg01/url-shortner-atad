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
			Div(
				H1(Class("text-3xl font-bold mb-4"), g.Text("Register")),

				Form(
					Method("POST"),
					Action("/register"),

					Label(For("email"), Class("block mb-1 font-semibold"), g.Text("Email:")),
					Input(Type("email"), Name("email"), Class("border rounded w-full p-2 mb-4"), ID("email")), Br(),

					Label(For("password"), Class("block mb-1 font-semibold"), g.Text("Password:")),
					Input(Type("password"), Name("password"), Class("border rounded w-full p-2 mb-4"), ID("password")), Br(),

					Button(Type("submit"), Class("bg-blue-600 text-white px-4 py-2 rounded hover:bg-blue-700"), g.Text("Register")),
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
