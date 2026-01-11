package frontend

import (
	"database/sql"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/diananeg01/url-shortner-atad/database"
	"github.com/diananeg01/url-shortner-atad/model"
	"github.com/google/uuid"
	g "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func UrlShortner(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_user")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	session, err := database.GetSession(cookie.Value)
	if err != nil || session.ExpiresAt.Before(time.Now()) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	user, err := database.GetUserByID(session.UserId)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	conn := database.GetDB()

	var (
		url          = ""
		chars        = ""
		shortenedUrl = ""
	)

	if r.Method == http.MethodPost {
		if err := r.ParseForm(); err == nil {
			url = r.FormValue("url")
			chars = r.FormValue("chars")
			expiresAtForm := r.FormValue("expiresAt")

			if chars == "" {
				const alphanumeric = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

				rand.Seed(time.Now().UnixNano())
				b := make([]byte, 8)
				for i := range b {
					b[i] = alphanumeric[rand.Intn(len(alphanumeric))]
				}
				chars = string(b)
			}
			shortenedUrl = "https://localhost:8080/redirect/" + chars + ""

			urlId := uuid.New()
			currentTime := time.Now()
			expiresAt := time.Time{}

			if expiresAtForm == "" {
				expiresAt = currentTime.Add(24 * time.Hour)
			} else {
				date, errDate := time.Parse("2006-01-02", expiresAtForm)
				if errDate != nil {
					http.Error(w, "Error when converting the date: "+errDate.Error(), http.StatusInternalServerError)
					return
				}
				expiresAt = time.Date(
					date.Year(), date.Month(), date.Day(),
					currentTime.Hour(), currentTime.Minute(), currentTime.Second(), currentTime.Nanosecond(),
					time.Local,
				)
			}

			_, errQuery := conn.Exec("INSERT INTO generated_url(url_id, url, short, crea, expires_at, user_id) values ($1, $2, $3, $4, $5, $6)", urlId, url, chars, currentTime, expiresAt, user.UserId)
			if errQuery != nil {
				http.Error(w, "Error when inserting the generated URL in the db: "+errQuery.Error(), http.StatusInternalServerError)
				return
			}
		}
	}

	_ = Page("Generate short URL",
		Div(Class("h-dvh bg-gray-900"),
			Div(Class("min-h-full"),
				Nav(Class("bg-gray-800/50"),
					Div(Class("max-w-7xl mx-auto px-4 sm:px-6 lg:px-8"),
						Div(Class("flex items-center justify-between h-16"),
							Div(Class("flex items-center"),
								Div(Class("flex-shrink-0"), A(Href("/"), Class("text-white font-semibold"), g.Text("URL Shortner"))),
								Div(Class("hidden md:block"),
									Div(Class("ml-10 flex items-baseline space-x-4"),
										A(Href("/"), Class("rounded-md bg-gray-950/50 px-3 py-2 text-sm font-medium text-white"), g.Text("Dashboard")),
										A(Href("/myStats"), Class("rounded-md px-3 py-2 text-sm font-medium text-gray-300 hover:bg-white/5 hover:text-white"), g.Text("User Analytics")),
										A(Href("/logout"), Class("rounded-md px-3 py-2 text-sm font-medium text-gray-300 hover:bg-white/5 hover:text-white"), g.Text("Logout")),
									),
								),
							),
							Div(Class("hidden md:block"),
								Div(Class("ml-4 flex items-center md:ml-6"), P(Class("ml-3 text-sm font-medium text-white"), g.Text("Signed in as "+user.Email))),
							),
						),
					),
				),

				Header(Class("relative bg-gray-800 after:pointer-events-none after:absolute after:inset-x-0 after:inset-y-0 after:border-y after:border-white/10"),
					Div(Class("max-w-7xl mx-auto px-4 py-6 sm:px-6 lg:px-8"),
						H1(Class("text-3xl font-bold text-white"), g.Text("Welcome to URL Shortner!")),
					),
				),

				Main(
					Div(Class("max-w-7xl mx-auto py-6 sm:px-6 lg:px-8"),
						Form(
							Method("POST"),
							Action("/submit"),
							Div(
								Label(For("url"), Class("block text-sm/6 font-medium text-gray-100"), g.Text("URL")),
								Div(Class("mt-2"), Input(Type("url"), ID("url"), Name("url"), Class("border rounded w-full p-2 mb-4"))),
							),
							Div(
								Label(For("chars"), Class("block text-sm/6 font-medium text-gray-100"), g.Text("Custom characters")),
								Div(Class("mt-2"), Input(Type("text"), ID("chars"), Name("chars"), Class("border rounded w-full p-2 mb-4"))),
							),
							Div(
								Label(For("expiresAt"), Class("block text-sm/6 font-medium text-gray-100"), g.Text("Date for expiration")),
								Div(Class("mt-2"), Input(Type("date"), ID("expiresAt"), Name("expiresAt"), Class("border rounded w-full p-2 mb-4"))),
							),
							Div(
								Button(Type("submit"), Class("flex w-50 justify-center rounded-md bg-indigo-500 px-3 py-1.5 text-sm/6 font-semibold text-white hover:bg-indigo-400 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-500"), g.Text("Submit")),
							),
						),

						Div(Class("border-t mt-10"),
							H2(Class("mt-10 text-center text-2xl/9 font-bold tracking-tight text-white"), g.Text("Submitted Data")),
							func() g.Node {
								if url == "" && chars == "" {
									return P(Class("mt-10 text-center text-sm/6 text-gray-400"), g.Text("No data submitted yet."))
								}
								return Div( //todo: add qr code
									P(Class("mb-2 text-center text-sm/6 text-gray-400"), g.Text(fmt.Sprintf("URL: %s", url))),
									P(Class("mb-4 text-center text-sm/6 text-gray-400"), g.Text("Shortened URL: "), A(Href(shortenedUrl), Target("_blank"), Class("text-blue-500 underline"), g.Text(shortenedUrl))),
								)
							}(),
						),
					),
				),
			),
		),
	).Render(w)
}

func RedirectURL(w http.ResponseWriter, r *http.Request) {
	shortUrl := r.PathValue("url")

	conn := database.GetDB()

	var url string
	var urlId uuid.UUID
	errQuery := conn.QueryRow("SELECT url_id, url FROM generated_url where short = $1 AND expires_at > $2",
		shortUrl, time.Now()).Scan(&urlId, &url)
	if errQuery != nil {
		http.Error(w, "URL expired or not found: "+errQuery.Error(), http.StatusNotFound)
		return
	}
	fmt.Println(url)

	var clickAnalyticId uuid.UUID
	errClick := conn.QueryRow("update url_analytics set value = value::int + 1 where url_id = $1 and keyname = $2 returning analytic_id", urlId, model.ClickAnalyticKeyname).
		Scan(&clickAnalyticId)
	if errors.Is(errClick, sql.ErrNoRows) {
		_, errUpsertClick := conn.Exec("INSERT INTO url_analytics(analytic_id, keyname, value, url_id) values ($1, $2, $3, $4) on conflict (analytic_id) do UPDATE SET value = excluded.value",
			uuid.New(), model.ClickAnalyticKeyname, 1, urlId)
		if errUpsertClick != nil {
			http.Error(w, "There was an error inserting the click analytic: "+errUpsertClick.Error(), http.StatusInternalServerError)
			return
		}
	} else if errClick != nil {
		http.Error(w, "There was an error updating the click analytic: "+errClick.Error(), http.StatusInternalServerError)
		return
	}

	var ip string
	forwarded := r.Header.Get("X-Forwarded-For")
	if forwarded != "" {
		ip = strings.Split(forwarded, ",")[0]
	} else {
		ip = strings.Split(r.RemoteAddr, ":")[0]
	}

	_, errUpsertVisitor := conn.Exec("INSERT INTO url_analytics(analytic_id, keyname, value, url_id) select $1, $2, $3, $4 where not exists (SELECT 1 FROM url_analytics WHERE url_id = $4 AND value = $3)",
		uuid.New(), model.VisitAnalyticKeyname, ip, urlId)
	if errUpsertVisitor != nil {
		http.Error(w, "There was an error inserting the click analytic: "+errUpsertVisitor.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}
