package emoji

import (
	"sort"
	"strings"
)

// SearchResult is a result from the Search function.
//
// Either Emoji or ImageURL will return a non-empty string, but not both.
type SearchResult struct {
	emoji *emoji
	score int
}

// Emoji is the Unicode emoji.
func (s SearchResult) Emoji() string {
	return s.emoji.emoji
}

// ImageURL is the URL of an image representing this emoji.
func (s SearchResult) ImageURL() string {
	return s.emoji.imageURL
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
	return defaultConfig.Search(query, max)
}

// Search returns a list of possible emoji for a query. The query is the text
// between the colon (:) and the user's cursor.
func (conf *Config) Search(query string, max int) []SearchResult {
	if query == "" || max <= 0 {
		return nil
	}

	query = strings.ToLower(query)

	results := make(searchResults, 0, max)

	results = conf.searchName(results, query, 3000)
	results = conf.searchDescription(results, query, 2000)
	results = conf.searchSet(results, query, 1000, conf.byTag, conf.tags, byTag, tags)
	results = conf.searchSet(results, query, 0, conf.byCategory, conf.categories, byCategory, categories)

	if len(results) < cap(results) {
		sort.Sort(results)
	}

	return results
}

func (conf *Config) searchName(results searchResults, query string, bonus int) searchResults {
	for name, e := range conf.byName {
		if result, ok := match(query, strings.Trim(name, ":"), e, bonus); ok {
			results = addResult(results, result)
		}
	}
	for name, e := range byName {
		if conf.overrides(e) {
			continue
		}

		if result, ok := match(query, strings.Trim(name, ":"), e, bonus); ok {
			results = addResult(results, result)
		}
	}

	return results
}

func (conf *Config) searchDescription(results searchResults, query string, bonus int) searchResults {
	for _, e := range conf.emoji {
		if result, ok := match(query, e.description, e, bonus); ok {
			results = addResult(results, result)
		}
	}

	for i := range allEmoji {
		e := &allEmoji[i]

		if conf.overrides(e) {
			continue
		}

		if result, ok := match(query, e.description, e, bonus); ok {
			results = addResult(results, result)
		}
	}

	return results
}

func (conf *Config) searchSet(results searchResults, query string, bonus int, byLocal [][]*emoji, namesLocal []string, by [][]*emoji, names []string) searchResults {
	for i, es := range byLocal {
		if result, ok := match(query, namesLocal[i], es[0], bonus); ok {
			results = addResult(results, result)
			for _, e := range es[1:] {
				results = addResult(results, SearchResult{e, result.score})
			}
		}
	}

	for i, es := range by {
		if result, ok := match(query, names[i], es[0], bonus); ok {
			if !conf.overrides(es[0]) {
				results = addResult(results, result)
			}
			for _, e := range es[1:] {
				if !conf.overrides(e) {
					results = addResult(results, SearchResult{e, result.score})
				}
			}
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
