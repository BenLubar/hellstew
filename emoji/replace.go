// Package emoji handles Unicode emoji.
package emoji

import (
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

var skipElement = map[atom.Atom]bool{
	atom.Abbr:   true,
	atom.Script: true,
	atom.Style:  true,
	atom.Pre:    true,
	atom.Code:   true,
}

// Replace finds Unicode emoji and emoji shortcodes (such as :tophat:) and
// replaces them with Unicode emoji with tooltips.
func Replace(nodes ...*html.Node) []*html.Node {
	return defaultConfig.Replace(nodes...)
}

// Replace finds Unicode emoji and emoji shortcodes (such as :tophat:) and
// replaces them with Unicode emoji with tooltips.
func (conf *Config) Replace(nodes ...*html.Node) []*html.Node {
	result := make([]*html.Node, 0, len(nodes))

	for _, node := range nodes {
		switch node.Type {
		case html.ElementNode:
			if node.Namespace != "" || skipElement[node.DataAtom] {
				result = append(result, node)
				break
			}

			result = append(result, conf.replaceElement(node)...)
		case html.TextNode:
			result = append(result, conf.replaceText(node)...)
		default:
			result = append(result, node)
		}
	}

	return result
}

func (conf *Config) replaceElement(node *html.Node) []*html.Node {
	result := node
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		if out := conf.Replace(child); len(out) != 1 || out[0] != child || result != node {
			if result == node {
				result = shallowClone(node)

				for prev := node.FirstChild; prev != child && prev != nil; prev = prev.NextSibling {
					result.AppendChild(deepClone(prev))
				}
			}
			for _, o := range out {
				result.AppendChild(o)
			}
		}
	}
	return []*html.Node{result}
}

func (conf *Config) replaceText(node *html.Node) []*html.Node {
	var matches [][2]int
	if conf.state == nil {
		matches = startState.match(node.Data)
	} else {
		matches = conf.state.match(node.Data)
	}
	if len(matches) == 0 {
		return []*html.Node{node}
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
		result = append(result, emojiToNode(e, name))
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

func emojiToNode(e *emoji, name string) *html.Node {
	if e.emoji != "" {
		abbr := &html.Node{
			Type:     html.ElementNode,
			Data:     "abbr",
			DataAtom: atom.Abbr,
			Attr: []html.Attribute{
				{
					Key: "title",
					Val: e.description,
				},
				{
					Key: "class",
					Val: "emoji",
				},
			},
		}
		abbr.AppendChild(&html.Node{
			Type: html.TextNode,
			Data: e.emoji,
		})
		return abbr
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
				Key: "title",
				Val: e.description,
			},
			{
				Key: "class",
				Val: "emoji",
			},
		},
	}
	return img
}
