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
		"padding-right: 1.313rem",
		"padding-top: 1.313rem",
		"height: 2.4rem",
	)
	css.Style(".fi",
		"font-weight: bold",
	)
	css.Style("span.timestamp",
		"float: right",
		"font-style: italic",
		"font-size: 0.825rem",
		"padding-right: 1.313rem",
		"line-height: 1.3rem",
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
		"font-size: 1.25rem",
		"line-height: 1.25",
		gopherblue,
	)
	css.Style("dl",
		"font-size: 0.875rem",
		"line-height: 1.3",
	)
	css.Style("dd.method",
		"padding-left: 1.25rem",
	)

	css.Style("pre",
		"background: #EFEFEF",
		"padding: 0.625rem",
		"border-radius: 0.3125rem",
		"margin: 1.25rem",
	)
	css.Style("code",
		"font-family: Menlo, monospace",
		"font-size: 0.875rem",
	)
	css.Style("p",
		"margin: 1.25rem",
		"max-width: 900px",
	)
	css.Style("a",
		gopherblue,
		"text-decoration: none",
	)
	css.Style("a:hover",
		gopherblue,
		"text-decoration: underline",
	)
	css.Style(".title",
		"padding-left: 1.25rem",
	)
	return css
}
