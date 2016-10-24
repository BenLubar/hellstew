package emoji_test

import (
	"testing"

	"github.com/BenLubar/hellstew/emoji"
)

var testConfig = func() *emoji.Config {
	var conf emoji.Config

	conf.AddImage("https://assets-cdn.github.com/images/icons/emoji/trollface.png", "trollface", []string{"trollface"}, "GitHub", nil)
	conf.AddImage("https://assets-cdn.github.com/images/icons/emoji/shipit.png", "ship it!", []string{"shipit", "squirrel"}, "GitHub", nil)
	conf.AddImage("https://assets-cdn.github.com/images/icons/emoji/octocat.png", "octocat", []string{"cat"}, "GitHub", nil)
	conf.AddImage("http://thedailywtf.com/favicon.ico", "wtf", []string{"wtf"}, "The Daily WTF", nil)
	conf.AddEmoji("⁉️", "backwards interrobang", []string{"wrongterrobang"}, "The Daily WTF", []string{"wrong"})

	return &conf
}()

func TestConfig(t *testing.T) {
	t.Run("AlreadyDefinedEmoji", func(t *testing.T) {
		defer expectPanic(t, "emoji: already defined in this Config: ⁉️")

		testConfig.AddEmoji("⁉️", "exclamation mark question mark", nil, "", nil)
	})

	t.Run("AlreadyDefinedEmojiAlias", func(t *testing.T) {
		defer expectPanic(t, "emoji: already defined in this Config: :trollface:")

		testConfig.AddEmoji("LOL", "trollface", []string{"trollface"}, "", nil)
	})

	t.Run("NoImageAliases", func(t *testing.T) {
		defer expectPanic(t, "emoji: image needs at least one alias")

		testConfig.AddImage("/images/unaliased.png", "", nil, "", nil)
	})

	t.Run("AlreadyDefinedImageAlias", func(t *testing.T) {
		defer expectPanic(t, "emoji: already defined in this Config: :wtf:")

		testConfig.AddImage("/images/wtf.png", "", []string{"worsethanfailure", "wtf"}, "", nil)
	})

	t.Run("EmptyEmoji", func(t *testing.T) {
		defer expectPanic(t, "emoji: emoji cannot be empty string")

		testConfig.AddEmoji("", "empty", []string{"empty"}, "", nil)
	})

	t.Run("EmptyAlias", func(t *testing.T) {
		defer expectPanic(t, "emoji: alias cannot be empty string")

		testConfig.AddImage("/images/emptyalias.png", "", []string{"emptyalias", ""}, "", nil)
	})

	t.Run("ColonAlias", func(t *testing.T) {
		defer expectPanic(t, "emoji: alias cannot contain ':'")

		testConfig.AddImage("/images/colonalias.png", "", []string{":alias:"}, "", nil)
	})
}

func expectPanic(t *testing.T, expected interface{}) {
	if r := recover(); r == nil {
		t.Errorf("Expected panic: %v", expected)
	} else if r != expected {
		t.Errorf("Expected panic: %v", expected)
		t.Errorf("Got panic: %v", r)
	}
}
