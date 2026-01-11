package frontend

import (
	"net/http"
	"strconv"
	"time"

	"github.com/diananeg01/url-shortner-atad/database"
	"github.com/diananeg01/url-shortner-atad/model"
	g "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func UrlAnalytics(w http.ResponseWriter, r *http.Request) {
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

	result, errSelect := conn.Query(
		`select u.url_id, u.url, u.short,
				   case
					   when expires_at > now() then 'false'
					   else 'true'
				   end as expired,
					coalesce(
						(select a1.value::int from url_analytics a1 left join generated_url u1 on u1.url_id = a1.url_id where a1.keyname = 'clicks' and u1.url_id = u.url_id), 0
					) as clicks,
					(
						select COALESCE(COUNT(DISTINCT a2.value), 0) from url_analytics a2 left join generated_url u2 on u2.url_id = a2.url_id where a2.keyname = 'visitor' and u2.url_id = u.url_id
					) as visitors
				from generated_url u
					where u.user_id = $1
				group by u.url_id`,
		user.UserId)

	if errSelect != nil {
		http.Error(w, "Error when selecting the analytics for current user: "+errSelect.Error(), http.StatusInternalServerError)
		return
	}

	var rows []model.TableRow

	for result.Next() {
		var s model.TableRow

		errScan := result.Scan(
			&s.UrlID,
			&s.LongURL,
			&s.ShortURL,
			&s.Expired,
			&s.TotalClicks,
			&s.UniqueClicks,
		)
		if errScan != nil {
			http.Error(w, "Error in for when converting the select result into table data: "+errScan.Error(), http.StatusInternalServerError)
			return
		}

		rows = append(rows, s)
	}

	if result.Err() != nil {
		http.Error(w, "Error after for when converting the select result into table data: "+result.Err().Error(), http.StatusInternalServerError)
		return
	}

	_ = Page(
		"URL Analytics for current user",
		Div(Class("h-dvh bg-gray-900"),
			Div(Class("min-h-full"),
				Nav(Class("bg-gray-800/50"),
					Div(Class("max-w-7xl mx-auto px-4 sm:px-6 lg:px-8"),
						Div(Class("flex items-center justify-between h-16"),
							Div(Class("flex items-center"),
								Div(Class("flex-shrink-0"), A(Href("/"), Class("text-white font-semibold"), g.Text("URL Shortner"))),
								Div(Class("hidden md:block"),
									Div(Class("ml-10 flex items-baseline space-x-4"),
										A(Href("/"), Class("rounded-md px-3 py-2 text-sm font-medium text-gray-300 hover:bg-white/5 hover:text-white"), g.Text("Dashboard")),
										A(Href("/myStats"), Class("rounded-md bg-gray-950/50 px-3 py-2 text-sm font-medium text-white"), g.Text("User Analytics")),
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
						H1(Class("text-3xl font-bold text-white"), g.Text("URL Analytics for current user")),
					),
				),

				Main(
					Div(Class("max-w-7xl mx-auto py-6 sm:px-6 lg:px-8"),
						Div(Class("relative flex flex-col w-full h-full overflow-scroll text-slate-300 bg-slate-800 shadow-md rounded-lg bg-clip-border"),
							DataTable(rows),
						),
					),
				),
			),
		),
	).Render(w)
}

func DataTable(rows []model.TableRow) g.Node {
	return Table(Class("w-full text-left table-auto min-w-max"),
		Tr(
			Th(Class("p-4 border-b border-slate-600 bg-slate-700"),
				P(Class("text-sm font-normal leading-none text-slate-300"), g.Text("Short URL"))),
			Th(Class("p-4 border-b border-slate-600 bg-slate-700"),
				P(Class("text-sm font-normal leading-none text-slate-300"), g.Text("Original URL"))),
			Th(Class("p-4 border-b border-slate-600 bg-slate-700"),
				P(Class("text-sm font-normal leading-none text-slate-300"), g.Text("Total Clicks"))),
			Th(Class("p-4 border-b border-slate-600 bg-slate-700"),
				P(Class("text-sm font-normal leading-none text-slate-300"), g.Text("Unique Visitors"))),
			Th(Class("p-4 border-b border-slate-600 bg-slate-700"),
				P(Class("text-sm font-normal leading-none text-slate-300"), g.Text("Expired"))),
		),
		g.Map(rows, func(r model.TableRow) g.Node {
			return Tr(Class("hover:bg-slate-700"),
				Td(Class("p-4 border-b border-slate-700"),
					P(Class("text-sm text-slate-100 font-semibold"), g.Text(r.ShortURL))),
				Td(Class("p-4 border-b border-slate-700"),
					P(Class("text-sm text-slate-300"), g.Text(r.LongURL))),
				Td(Class("p-4 border-b border-slate-700"),
					P(Class("text-sm text-slate-300"), g.Text(strconv.Itoa(r.TotalClicks)))),
				Td(Class("p-4 border-b border-slate-700"),
					P(Class("text-sm text-slate-300"), g.Text(strconv.Itoa(r.UniqueClicks)))),
				Td(Class("p-4 border-b border-slate-700"),
					P(Class("text-sm text-slate-300"), g.Text(r.Expired))),
			)
		}),
	)
}
