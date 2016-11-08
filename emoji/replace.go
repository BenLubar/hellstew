// Package emoji handles Unicode emoji.
package emoji

import (
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

var skipElement = map[atom.Atom]bool{
	atom.Script: true,
	atom.Style:  true,
	atom.Pre:    true,
	atom.Code:   true,
}

// Replace finds Unicode emoji and emoji shortcodes (such as :tophat:) and
// replaces them with Unicode emoji with tooltips. Replace is idempotent.
func Replace(nodes ...*html.Node) []*html.Node {
	return defaultConfig.Replace(nodes...)
}

// Replace finds Unicode emoji and emoji shortcodes (such as :tophat:) and
// replaces them with Unicode emoji with tooltips. Replace is idempotent.
func (conf *Config) Replace(nodes ...*html.Node) []*html.Node {
	return conf.replace(true, nodes...)
}

func (conf *Config) replace(tooltip bool, nodes ...*html.Node) []*html.Node {
	result := make([]*html.Node, 0, len(nodes))

	for _, node := range nodes {
		switch node.Type {
		case html.ElementNode:
			if node.Namespace != "" || skipElement[node.DataAtom] {
				result = append(result, deepClone(node))
				break
			}

			result = append(result, conf.replaceElement(tooltip, node)...)
		case html.TextNode:
			result = append(result, conf.replaceText(tooltip, node)...)
		default:
			result = append(result, deepClone(node))
		}
	}

	return result
}

func (conf *Config) replaceElement(tooltip bool, node *html.Node) []*html.Node {
	if tooltip {
		for _, a := range node.Attr {
			if a.Namespace == "" && a.Key == "title" {
				tooltip = false
				break
			}
		}
	}
	for _, a := range node.Attr {
		if a.Namespace == "" && a.Key == "class" && a.Val == "emoji" {
			return []*html.Node{deepClone(node)}
		}
	}

	result := shallowClone(node)
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		for _, o := range conf.replace(tooltip, child) {
			result.AppendChild(o)
		}
	}
	return []*html.Node{result}
}

func (conf *Config) replaceText(tooltip bool, node *html.Node) []*html.Node {
	var matches [][2]int
	if conf.state == nil {
		matches = startState.match(node.Data)
	} else {
		matches = conf.state.match(node.Data)
	}
	if len(matches) == 0 {
		return []*html.Node{shallowClone(node)}
	}

	result := make([]*html.Node, 0, len(matches)*2+1)
	for i, match := range matches {
		if i == 0 {
			if match[0] != 0 {
				result = append(result, &html.Node{
					Type: html.TextNode,
					Data: node.Data[:match[0]],
				})
			}
		}
		name := node.Data[match[0]:match[1]]
		e, ok := conf.byName[name]
		if !ok {
			e = byName[name]
		}
		result = append(result, emojiToNode(tooltip, e, name))
		if i+1 == len(matches) {
			if match[1] != len(node.Data) {
				result = append(result, &html.Node{
					Type: html.TextNode,
					Data: node.Data[match[1]:],
				})
			}
		} else if next := matches[i+1]; match[1] != next[0] {
			result = append(result, &html.Node{
				Type: html.TextNode,
				Data: node.Data[match[1]:next[0]],
			})
		}
	}

	return result
}

func shallowClone(node *html.Node) *html.Node {
	result := &html.Node{
		Namespace: node.Namespace,
		Type:      node.Type,
		Attr:      make([]html.Attribute, len(node.Attr)),
		Data:      node.Data,
		DataAtom:  node.DataAtom,
	}
	copy(result.Attr, node.Attr)
	return result
}

func deepClone(node *html.Node) *html.Node {
	result := shallowClone(node)
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		result.AppendChild(deepClone(child))
	}
	return result
}

func emojiToNode(tooltip bool, e *emoji, name string) *html.Node {
	if e.emoji != "" {
		node := &html.Node{
			Type:     html.ElementNode,
			Data:     "span",
			DataAtom: atom.Span,
			Attr: []html.Attribute{
				{
					Key: "class",
					Val: "emoji",
				},
			},
		}
		if tooltip {
			node.Data = "abbr"
			node.DataAtom = atom.Abbr
			node.Attr = append(node.Attr, html.Attribute{
				Key: "title",
				Val: e.description,
			})
		}

		node.AppendChild(&html.Node{
			Type: html.TextNode,
			Data: e.emoji,
		})
		return node
	}
	img := &html.Node{
		Type:     html.ElementNode,
		Data:     "img",
		DataAtom: atom.Img,
		Attr: []html.Attribute{
			{
				Key: "src",
				Val: e.imageURL,
			},
			{
				Key: "alt",
				Val: name,
			},
			{
				Key: "class",
				Val: "emoji",
			},
		},
	}
	if tooltip {
		img.Attr = append(img.Attr, html.Attribute{
			Key: "title",
			Val: e.description,
		})
	}
	return img
}
