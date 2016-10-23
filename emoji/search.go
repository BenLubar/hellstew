package emoji

import (
	"sort"
	"strings"
)

// SearchResult is a result from the Search function.
type SearchResult struct {
	emoji *emoji
	score int
}

// Emoji is the Unicode emoji.
func (s SearchResult) Emoji() string {
	return s.emoji.emoji
}

// Aliases is a slice of textual shortcodes that can be used between colons to
// represent the emoji.
func (s SearchResult) Aliases() []string {
	return s.emoji.aliases
}

// Description is the English textual description of the emoji.
func (s SearchResult) Description() string {
	return s.emoji.description
}

// Score is the likelihood of the result being correct. Higher is better.
func (s SearchResult) Score() int {
	return s.score
}

type searchResults []SearchResult

// Len implements sort.Interface.
func (s searchResults) Len() int {
	return len(s)
}

// Swap implements sort.Interface.
func (s searchResults) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// Less implements sort.Interface.
func (s searchResults) Less(i, j int) bool {
	if s[i].score == s[j].score {
		return strings.ToLower(s[i].emoji.description) < strings.ToLower(s[j].emoji.description)
	}
	return s[i].score > s[j].score
}

// Search returns a list of possible emoji for a query. The query is the text
// between the colon (:) and the user's cursor.
func Search(query string, max int) []SearchResult {
	if query == "" || max <= 0 {
		return nil
	}

	query = strings.ToLower(query)

	results := make(searchResults, 0, max)

	results = searchName(results, query, 3000)
	results = searchDescription(results, query, 2000)
	results = searchSet(results, query, 1000, byTag, tags)
	results = searchSet(results, query, 0, byCategory, categories)

	if len(results) < cap(results) {
		sort.Sort(results)
	}

	return results
}

func searchName(results searchResults, query string, bonus int) searchResults {
	for name, e := range byName {
		result, ok := match(query, strings.Trim(name, ":"), e, bonus)
		if !ok {
			continue
		}

		results = addResult(results, result)
	}

	return results
}

func searchDescription(results searchResults, query string, bonus int) searchResults {
	for i := range allEmoji {
		e := &allEmoji[i]
		result, ok := match(query, e.description, e, bonus)
		if !ok {
			continue
		}

		results = addResult(results, result)
	}

	return results
}

func searchSet(results searchResults, query string, bonus int, by [][]*emoji, names []string) searchResults {
	for i, es := range by {
		result, ok := match(query, names[i], es[0], bonus)
		if !ok {
			continue
		}

		results = addResult(results, result)
		for _, e := range es[1:] {
			results = addResult(results, SearchResult{e, result.score})
		}
	}

	return results
}

func addResult(results searchResults, result SearchResult) searchResults {
	for _, r := range results {
		if r.emoji == result.emoji {
			// new score is always lower
			return results
		}
	}

	if len(results) < cap(results) {
		results = append(results, result)
		if len(results) == cap(results) {
			sort.Sort(results)
		}
	} else {
		i := sort.Search(len(results), func(i int) bool {
			if results[i].score == result.score {
				return strings.ToLower(results[i].emoji.description) >= strings.ToLower(result.emoji.description)
			}
			return results[i].score < result.score
		})
		if i < len(results) {
			copy(results[i+1:], results[i:])
			results[i] = result
		}
	}
	return results
}

func match(query, actual string, e *emoji, bonus int) (SearchResult, bool) {
	actual = strings.ToLower(actual)

	if query == actual {
		return SearchResult{e, 500 + bonus}, true
	}

	if strings.HasPrefix(actual, query) {
		return SearchResult{e, len(query)*3 - len(actual) + bonus}, true
	}

	if strings.Contains(actual, query) {
		return SearchResult{e, len(query)*2 - len(actual) + bonus}, true
	}

	return SearchResult{}, false
}
