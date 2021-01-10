package main

import . "github.com/gregoryv/web"

func theme() *CSS {
	css := NewCSS()
	css.Style("html, body",
		"margin: 0 0",
		"padding: 0 0",
	)
	css.Style("body",
		"padding: 1em 1.618em",
		"padding-bottom: 1640px",
		"font-family: sans-serif",
	)
	css.Style("h1",
		"font-size: 1.618em",
		"border-bottom: 3px solid rgb(55, 94, 171)",
	)
	css.Style("h2",
		"font-size: 1em",
	)
	css.Style("h3",
		"font-size: 1em",
	)
	return css
}
