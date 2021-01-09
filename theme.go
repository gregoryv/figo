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
		"font-family: sans-serif",
	)
	css.Style("li.h3",
		"margin-left: 1.618em",
	)
	css.Style("a:link, a:visited",
		"color: rgb(55, 94, 171)", // godoc blue
		"text-decoration: none",
	)
	css.Style("a:hover",
		"text-decoration: underline",
	)
	return css
}
