package emoji

import "strings"

// Config is a custom emoji set that extends the default set.
type Config struct {
	state      *state
	emoji      []*emoji
	byName     map[string]*emoji
	tags       []string
	byTag      [][]*emoji
	categories []string
	byCategory [][]*emoji
}

var defaultConfig = &Config{
	state: startState,
}

func (conf *Config) overrides(e *emoji) bool {
	/*
		// We don't currently have any images in the default set.
		if e.emoji == "" {
			_, ok := conf.byName[":"+e.aliases[0]+":"]
			return ok
		}
	*/
	_, ok := conf.byName[e.emoji]
	return ok
}

// AddEmoji adds a custom Unicode emoji to the Config.
func (conf *Config) AddEmoji(unicodeEmoji, description string, aliases []string, category string, tags []string) {
	if unicodeEmoji == "" {
		panic("emoji: emoji cannot be empty string")
	}
	if _, ok := conf.byName[unicodeEmoji]; ok {
		panic("emoji: already defined in this Config: " + unicodeEmoji)
	}
	conf.validateAliases(aliases)

	e := &emoji{
		emoji:       unicodeEmoji,
		description: description,
		aliases:     aliases,
	}

	conf.addEmoji(e, aliases, category, tags)
	conf.addName(unicodeEmoji, e)
}

// AddImage adds an image as a pseudo-emoji. At least one alias is required.
func (conf *Config) AddImage(imageURL, description string, aliases []string, category string, tags []string) {
	if len(aliases) == 0 {
		panic("emoji: image needs at least one alias")
	}
	conf.validateAliases(aliases)

	e := &emoji{
		imageURL:    imageURL,
		description: description,
		aliases:     aliases,
	}

	conf.addEmoji(e, aliases, category, tags)
}

func (conf *Config) validateAliases(aliases []string) {
	for _, a := range aliases {
		if a == "" {
			panic("emoji: alias cannot be empty string")
		}
		if strings.ContainsRune(a, ':') {
			panic("emoji: alias cannot contain ':'")
		}
		if _, ok := conf.byName[":"+a+":"]; ok {
			panic("emoji: already defined in this Config: :" + a + ":")
		}
	}
}

func (conf *Config) addEmoji(e *emoji, aliases []string, category string, tags []string) {
	conf.emoji = append(conf.emoji, e)
	if conf.byName == nil {
		conf.byName = make(map[string]*emoji)
	}
	for _, a := range aliases {
		conf.addName(":"+a+":", e)
	}
	addBy(&conf.byCategory, &conf.categories, e, category)
	for _, tag := range tags {
		addBy(&conf.byTag, &conf.tags, e, tag)
	}
}

func (conf *Config) addName(name string, e *emoji) {
	conf.byName[name] = e
	if conf.state == nil {
		conf.state = startState.add(strings.ToLower(name))
	} else {
		conf.state = conf.state.add(strings.ToLower(name))
	}
}

func addBy(by *[][]*emoji, names *[]string, e *emoji, name string) {
	for i, n := range *names {
		if n == name {
			(*by)[i] = append((*by)[i], e)
			return
		}
	}
	*names = append(*names, name)
	*by = append(*by, []*emoji{e})
}
