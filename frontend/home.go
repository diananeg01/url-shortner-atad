package frontend

import (
	"fmt"
	"net/http"
	"time"

	"github.com/diananeg01/url-shortner-atad/database"
	"github.com/google/uuid"
	g "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func UrlShortner(w http.ResponseWriter, r *http.Request) {
	userID, ok := getSessionUser(r)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	conn := database.GetDB()

	var email string
	conn.QueryRow(`SELECT email FROM user_data WHERE user_id = $1`, userID).
		Scan(&email)

	var (
		url          = ""
		chars        = ""
		shortenedUrl = ""
	)

	if r.Method == http.MethodPost {
		if err := r.ParseForm(); err == nil {
			url = r.FormValue("url")
			chars = r.FormValue("chars")
			shortenedUrl = "https://localhost:8080/redirect" + chars + ""

			urlId := uuid.New()
			currentTime := time.Now()

			_, errQuery := conn.Exec("INSERT INTO generated_url(url_id, url, short, crea, lupa, status, user_id) values ($1, $2, $3, $4, $5, $6, $7)", urlId, url, chars, currentTime, currentTime, "active", "bda5e4b5-af8b-4c42-8733-09be226c8695")
			if errQuery != nil {
				return
			}
		}
	}

	_ = Page("Generate short URL",
		Div(
			H1(Class("text-3xl font-bold mb-4"), g.Text("Welcome, "+email+"!")),
			A(Href("/logout"), g.Text("Logout")),
			P(Class("text-gray-700 mb-6"), g.Text("This is a simple URL shortener.")),

			Form(
				Method("POST"),
				Action("/submit"),
				Div(
					Label(For("url"), Class("block mb-1 font-semibold"), g.Text("URL")),
					Input(Type("text"), ID("url"), Name("url"), Class("border rounded w-full p-2 mb-4")),
				),
				Div(
					Label(For("chars"), Class("block mb-1 font-semibold"), g.Text("Custom characters")),
					Input(Type("chars"), ID("chars"), Name("chars"), Class("border rounded w-full p-2 mb-4")),
				),
				Button(Type("submit"), Class("bg-blue-600 text-white px-4 py-2 rounded hover:bg-blue-700"), g.Text("Submit")),
			),

			Div(
				Class("border-t pt-4"),
				H2(Class("text-2xl font-semibold mb-2"), g.Text("Submitted Data")),
				func() g.Node {
					if url == "" && chars == "" {
						return P(Class("text-gray-500 italic"), g.Text("No data submitted yet."))
					}
					return Div(
						P(Class("mb-2"), g.Text(fmt.Sprintf("URL: %s", url))),
						P(Class("mb-4"), g.Text("Shortened URL: "), A(Href(url), Class("text-blue-500 underline"), g.Text(shortenedUrl))),
					)
				}(),
			),
		),
	).Render(w)
}

func RedirectURL(w http.ResponseWriter, r *http.Request) {
	shortUrl := r.PathValue("url")

	conn := database.GetDB()
	var url string

	errQuery := conn.QueryRow("SELECT url FROM generated_url where short = $1", shortUrl).Scan(&url)
	if errQuery != nil {
		return
	}
	fmt.Println(url)

	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// SESSION HELPER
func getSessionUser(r *http.Request) (int, bool) {
	cookie, err := r.Cookie("session_user")
	if err != nil || cookie.Value == "" {
		return 0, false
	}

	var id int
	fmt.Sscanf(cookie.Value, "%d", &id)
	return id, true
}
