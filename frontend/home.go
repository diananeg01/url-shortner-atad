package frontend

import (
	g "maragu.dev/gomponents"
	c "maragu.dev/gomponents/components"
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
