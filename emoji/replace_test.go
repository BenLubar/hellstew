package emoji_test

import (
	"bytes"
	"strings"
	"testing"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"

	"github.com/BenLubar/hellstew/emoji"
)

type replaceTest struct {
	name   string
	input  string
	output string
}

var replaceTests = [...]replaceTest{
	{
		name:   "Empty",
		input:  ``,
		output: ``,
	},
	{
		name:   "PlainText",
		input:  `Hello, world!`,
		output: `Hello, world!`,
	},
	{
		name:   "Comment",
		input:  `:<!-- skip -->horse_racing:`,
		output: `:<!-- skip -->horse_racing:`,
	},
	{
		name:   "Colons",
		input:  `:horse_racing:`,
		output: `<abbr title="horse racing" class="emoji">🏇</abbr>`,
	},
	{
		name:   "Unicode",
		input:  `🏇`,
		output: `<abbr title="horse racing" class="emoji">🏇</abbr>`,
	},
	{
		name:   "Mixed",
		input:  `<em>To the 🍿 thread!</em> :musical_note:`,
		output: `<em>To the <abbr title="popcorn" class="emoji">🍿</abbr> thread!</em> <abbr title="musical note" class="emoji">🎵</abbr>`,
	},
	{
		name:   "Garbage",
		input:  `:po:popcor:corn:n:`,
		output: `:po:popcor<abbr title="ear of corn" class="emoji">🌽</abbr>n:`,
	},
	{
		name:   "Code",
		input:  `<code>:tangerine:</code>`,
		output: `<code>:tangerine:</code>`,
	},
	{
		name:   "Attribute",
		input:  `<a href=":book:">:book:</a>`,
		output: `<a href=":book:"><abbr title="open book" class="emoji">📖</abbr></a>`,
	},
	{
		name:   "Nested",
		input:  `<p><a href="https://www.google.com/"><img src="https://www.google.com/favicon.ico" alt="Google"/></a> :mag: Look it up:exclamation:</p>`,
		output: `<p><a href="https://www.google.com/"><img src="https://www.google.com/favicon.ico" alt="Google"/></a> <abbr title="left-pointing magnifying glass" class="emoji">🔍</abbr> Look it up<abbr title="exclamation mark" class="emoji">❗️</abbr></p>`,
	},
}

func TestReplace(t *testing.T) {
	t.Run("Global", func(t *testing.T) {
		testReplace(t, emoji.Replace)
	})
	t.Run("Config", func(t *testing.T) {
		testReplace(t, testConfig.Replace)
	})
	t.Run("EmptyConfig", func(t *testing.T) {
		var conf emoji.Config
		testReplace(t, conf.Replace)
	})
	t.Run("ConfigSpecific", func(t *testing.T) {
		nodes := testConfig.Replace(&html.Node{
			Type: html.TextNode,
			Data: ":wtf:",
		})

		if len(nodes) != 1 {
			t.Fatalf("unexpected len(nodes) == %d", len(nodes))
		}
		if nodes[0].Type != html.ElementNode {
			t.Errorf("unexpected nodes[0].Type == %d", nodes[0].Type)
		}
		if nodes[0].Data != "img" {
			t.Errorf("unexpected nodes[0].Data == %q", nodes[0].Data)
		}
	})
}

func testReplace(t *testing.T, replace func(...*html.Node) []*html.Node) {
	for _, tt := range replaceTests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			input, err := html.ParseFragment(strings.NewReader(tt.input), &html.Node{
				Type:     html.ElementNode,
				Data:     "div",
				DataAtom: atom.Div,
			})
			if err != nil {
				t.Fatal(err)
			}

			nodes := replace(input...)

			var buf bytes.Buffer
			for _, n := range nodes {
				err = html.Render(&buf, n)
				if err != nil {
					t.Fatal(err)
				}
			}

			if output := buf.String(); tt.output != output {
				t.Errorf("input %q\nexpected %q\nactual   %q", tt.input, tt.output, output)
			}
		})
	}
}
