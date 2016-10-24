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
	var conf emoji.Config

	addGitHubEmoji := func(name string, aliases ...string) {
		conf.AddImage("https://assets-cdn.github.com/images/icons/emoji/"+name+".png", ":"+name+":", append([]string{name}, aliases...), "GitHub", nil)
	}

	addGitHubEmoji("basecamp")
	addGitHubEmoji("basecampy")
	addGitHubEmoji("bowtie")
	addGitHubEmoji("feelsgood")
	addGitHubEmoji("finnadie")
	addGitHubEmoji("goberserk")
	addGitHubEmoji("godmode")
	addGitHubEmoji("hurtrealbad")
	addGitHubEmoji("neckbeard")
	addGitHubEmoji("octocat")
	addGitHubEmoji("rage1")
	addGitHubEmoji("rage2")
	addGitHubEmoji("rage3")
	addGitHubEmoji("rage4")
	addGitHubEmoji("shipit", "squirrel")
	addGitHubEmoji("suspect")
	addGitHubEmoji("trollface")

	js.Global.Call("$", "#emoji").Call("autoComplete", js.M{
		"cache":    false,
		"minChars": 1,
		"source": func(term string, response func([]string)) {
			results := conf.Search(term, 10)
			output := make([]string, len(results))

			var buf bytes.Buffer

			for i, item := range results {
				replacement := item.Emoji()
				if replacement == "" {
					replacement = ":" + item.Aliases()[0] + ":"
				}
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
							Val: replacement,
						},
					},
				}
				for _, node := range conf.Replace(&html.Node{
					Type: html.TextNode,
					Data: replacement,
				}) {
					div.AppendChild(node)
				}
				div.AppendChild(&html.Node{
					Type: html.TextNode,
					Data: " " + item.Description(),
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

		nodes = conf.Replace(nodes...)

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
