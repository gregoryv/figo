package figo

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
	css.Style("li.h1",
		"margin-top: 1.618em",
	)
	css.Style("li.h2",
		"margin-left: 1.618em",
	)
	css.Style("section.func:last-child",
		"margin-bottom: 3.618em",
	)
	css.Style("h1",
		"font-size: 1.618em",
		"border-bottom: 3px solid rgb(55, 94, 171)",
	)
	css.Style("h2",
		"font-size: 1em",
	)
	css.Style("a:link, a:visited",
		"color: rgb(55, 94, 171)", // godoc blue
		"text-decoration: none",
	)
	css.Style("a:hover",
		"text-decoration: underline",
	)
	css.Style("li",
		"list-style-type: none",
	)
	css.Style(".empty",
		"display: none",
	)
	return css
}
