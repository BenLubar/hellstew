package emoji

import "strings"

type state struct {
	next [256]*state
	term bool
}

var startState = func() *state {
	var root state
	for name := range byName {
		root.addInPlace(strings.ToLower(name))
	}
	return &root
}()

func (s *state) addInPlace(name string) {
	for {
		if name == "" {
			s.term = true
			return
		}
		if s.next[name[0]] == nil {
			s.next[name[0]] = makeState(name[1:])
			return
		}
		s = s.next[name[0]]
		name = name[1:]
	}
}

func (s *state) add(name string) *state {
	if name == "" {
		return &state{
			next: s.next,
			term: true,
		}
	}
	clone := *s
	if clone.next[name[0]] != nil {
		clone.next[name[0]] = clone.next[name[0]].add(name[1:])
	} else {
		clone.next[name[0]] = makeState(name[1:])
	}
	return &clone
}

func makeState(name string) *state {
	s := &state{
		term: true,
	}

	for i := len(name) - 1; i >= 0; i-- {
		next := &state{}
		next.next[name[i]] = s
		s = next
	}

	return s
}

func (s *state) match(str string) [][2]int {
	var matches [][2]int

	for i := 0; i < len(str); i++ {
		term := -1
		for j, c := i, s; ; j++ {
			if c.term {
				term = j
			}
			if j == len(str) || c.next[str[j]] == nil {
				break
			}
			c = c.next[str[j]]
		}
		if term != -1 {
			matches = append(matches, [2]int{i, term})
			i = term - 1
		}
	}

	return matches
}
