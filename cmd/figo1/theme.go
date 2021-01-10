package main

import . "github.com/gregoryv/web"

func theme() *CSS {
	css := NewCSS()
	gopherblue := "color: #375EAB"
	css.Style("html, body",
		"margin: 0 0",
		"padding: 0 0",
	)
	css.Style("body",
		"margin: 0",
		"font-family: Arial, sans-serif",
		"background-color: #fff",
		"line-height: 1.3",
		"color: #222",
	)
	css.Style("div.top",
		"background-color: #E0EBF5",
		"font-size: 1.25rem",
		"padding-left: 1.313rem",
		"padding-top: 1.313rem",
		"height: 2.4rem",
	)
	css.Style("article",
		"padding: 1.313rem 1.313rem",
	)
	css.Style("h1",
		"margin-top: 0",
		"font-size: 1.75rem",
		"line-height: 1",
		gopherblue,
	)
	css.Style("h2",
		"font-size: 1.25rem",
		"background: #E0EBF5",
		"padding: 0.5rem",
		"line-height: 1.25",
		"font-weight: normal",
		"overflow: auto",
		"overflow-wrap: break-word",
		gopherblue,
	)
	css.Style("h3",
		"font-size: 1em",
	)
	css.Style("dl",
		"margin: 1.25rem",
	)
	return css
}
