package main

import (
	"bytes"
	"strings"

	"github.com/BenLubar/hellstew/emoji"
	"github.com/gopherjs/gopherjs/js"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

func main() {
	js.Global.Call("$", "#emoji").Call("autoComplete", js.M{
		"cache":    false,
		"minChars": 1,
		"source": func(term string, response func([]string)) {
			results := emoji.Search(term, 10)
			output := make([]string, len(results))

			var buf bytes.Buffer

			for i, item := range results {
				div := &html.Node{
					Type:     html.ElementNode,
					Data:     "div",
					DataAtom: atom.Div,
					Attr: []html.Attribute{
						{
							Key: "class",
							Val: "autocomplete-suggestion",
						},
						{
							Key: "data-val",
							Val: item.Emoji(),
						},
					},
				}
				div.AppendChild(&html.Node{
					Type: html.TextNode,
					Data: item.Emoji() + " " + item.Description(),
				})

				buf.Reset()
				err := html.Render(&buf, div)
				if err != nil {
					panic(err)
				}
				output[i] = buf.String()
			}
			response(output)
		},
		"renderItem": func(item, search string) string {
			return item
		},
	})
	js.Global.Call("$", "#html-button").Call("on", "click", func() {
		text := js.Global.Call("$", "#html")
		nodes, err := html.ParseFragment(strings.NewReader(text.Call("val").String()), &html.Node{
			Type:     html.ElementNode,
			Data:     "div",
			DataAtom: atom.Div,
		})
		if err != nil {
			panic(err)
		}

		nodes = emoji.Replace(nodes...)

		var buf bytes.Buffer
		for _, node := range nodes {
			err = html.Render(&buf, node)
			if err != nil {
				panic(err)
			}
		}
		text.Call("val", buf.String())
	}).Call("removeAttr", "disabled")
}
