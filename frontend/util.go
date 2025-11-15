package frontend

import (
	"github.com/go-echarts/go-echarts/v2/opts"
	g "maragu.dev/gomponents"
	c "maragu.dev/gomponents/components"
	"math/rand"
)

func Page(title string, body g.Node) g.Node {
	return c.HTML5(
		c.HTML5Props{
			Title: title,
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
