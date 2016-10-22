package emoji_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/BenLubar/hellstew/emoji"
)

type searchTest struct {
	query    string
	max      int
	expected []string
}

var searchTests = [...]searchTest{
	{
		query:    "",
		max:      10,
		expected: []string{},
	},
	{
		query:    "hap",
		max:      -1,
		expected: []string{},
	},
	{
		query: "mouse",
		max:   10,
		expected: []string{
			"mouse face",
			"mouse",
			"computer mouse",
		},
	},
	{
		query: "cl",
		max:   10,
		expected: []string{
			"CL button",
			"clapping hands",
			"clamp",
			"cloud",
			"club suit",
			"eight o’clock",
			"five o’clock",
			"four o’clock",
			"nine o’clock",
			"one o’clock",
		},
	},
	{
		query: "nature",
		max:   50,
		expected: []string{
			"ant",
			"baby chick",
			"bear face",
			"bird",
			"blossom",
			"blowfish",
			"boar",
			"bouquet",
			"bug",
			"cactus",
			"camel",
			"cat",
			"cat face",
			"cherry blossom",
			"chestnut",
			"chicken",
			"chipmunk",
			"Christmas tree",
			"cloud",
			"cloud with lightning",
			"cloud with lightning and rain",
			"cloud with rain",
			"cloud with snow",
			"collision",
			"comet",
			"cow",
			"cow face",
			"crab",
			"crescent moon",
			"crocodile",
			"dashing away",
			"deciduous tree",
			"dizzy",
			"dog",
			"dog face",
			"dolphin",
			"dove",
			"dragon",
			"dragon face",
			"droplet",
			"elephant",
			"evergreen tree",
			"fallen leaf",
			"fire",
			"first quarter moon",
			"first quarter moon with face",
			"fish",
			"fog",
			"four leaf clover",
			"frog face",
		},
	},
}

func TestSearch(t *testing.T) {
	for _, tt := range searchTests {
		tt := tt
		t.Run(fmt.Sprintf("%s_%d", tt.query, tt.max), func(t *testing.T) {
			results := emoji.Search(tt.query, tt.max)
			fail := len(results) != len(tt.expected)
			for i := 0; !fail && i < len(results); i++ {
				fail = results[i].Description() != tt.expected[i]
			}

			if fail {
				t.Errorf("#: actual / expected")
				for i := 0; i < len(results) || i < len(tt.expected); i++ {
					if i < len(results) && i < len(tt.expected) {
						t.Errorf("%d: %q / %q", i, results[i].Description(), tt.expected[i])
					} else if i < len(results) {
						t.Errorf("%d: %q / [none]", i, results[i].Description())
					} else {
						t.Errorf("%d: [none] / %q", i, tt.expected[i])
					}
				}
			}
		})
	}
}

func TestSearchResult(t *testing.T) {
	results := emoji.Search("minidisc", 1)
	if len(results) != 1 {
		t.Fatalf("no results for minidisc")
	}

	if expected, actual := []string{"minidisc"}, results[0].Aliases(); !reflect.DeepEqual(expected, actual) {
		t.Errorf("Aliases: %#v != %#v", expected, actual)
	}
	if expected, actual := "computer disk", results[0].Description(); expected != actual {
		t.Errorf("Description: %q != %q", expected, actual)
	}
	if expected, actual := "💽", results[0].Emoji(); expected != actual {
		t.Errorf("Emoji: %q != %q", expected, actual)
	}
	if expected, actual := 3500, results[0].Score(); expected != actual {
		t.Errorf("Score: %d != %d", expected, actual)
	}
}
