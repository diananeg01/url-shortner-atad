package frontend

import (
	"fmt"
	"net/http"
	"time"

	"github.com/diananeg01/url-shortner-atad/database"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/types"
	"github.com/google/uuid"
	g "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func TestChartRender(w http.ResponseWriter, _ *http.Request) {
	// create a new line instance
	line := charts.NewLine()
	// set some global options like Title/Legend/ToolTip or anything else
	line.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{Theme: types.ThemeWesteros}),
		charts.WithTitleOpts(opts.Title{
			Title:    "Line example in Westeros theme",
			Subtitle: "Line chart rendered by the http server this time",
		}))

	// Put data into instance
	line.SetXAxis([]string{"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"}).
		AddSeries("Category A", generateLineItems()).
		AddSeries("Category B", generateLineItems()).
		SetSeriesOptions(charts.WithLineChartOpts(opts.LineChart{Smooth: opts.Bool(true)}))
	line.Render(w)
}

func UrlShortner(w http.ResponseWriter, r *http.Request) {
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

			conn := database.GetDB()

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
			H1(Class("text-3xl font-bold mb-4"), g.Text("Welcome!")),
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
