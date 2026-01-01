package frontend

import (
	"math/rand"

	"github.com/go-echarts/go-echarts/v2/opts"
	g "maragu.dev/gomponents"
	c "maragu.dev/gomponents/components"
	. "maragu.dev/gomponents/html"
)

func Page(title string, body g.Node) g.Node {
	return c.HTML5(
		c.HTML5Props{
			Title: title,
			Head: []g.Node{
				// Example: add Tailwind CSS or your own stylesheet
				Script(Src("https://cdn.tailwindcss.com")),
			},
			Body: []g.Node{
				body,
			},
		},
	)
}

// generate random data for line chart
func generateLineItems() []opts.LineData {
	items := make([]opts.LineData, 0)
	for i := 0; i < 7; i++ {
		items = append(items, opts.LineData{Value: rand.Intn(300)})
	}
	return items
}
