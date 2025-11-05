package frontend

import (
	"fmt"
	"net/http"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/types"
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

func TestTitle(w http.ResponseWriter, r *http.Request) {
	// Render the HTML page
	_ = Page("Hello", H1(g.Text("Hello, Gomponents!"))).Render(w)
}

func UrlShortner(w http.ResponseWriter, r *http.Request) {
	var (
		url   = ""
		chars = ""
	)

	if r.Method == http.MethodPost {
		if err := r.ParseForm(); err == nil {
			url = r.FormValue("url")
			chars = r.FormValue("chars")
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
						P(Class("mb-4"), g.Text(fmt.Sprintf("Custom characters: %s", chars))),
					)
				}(),
			),
		),
	).Render(w)
}
