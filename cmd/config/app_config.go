package config

type BlogConfig struct {
	titleMaxLength int
	bodyMaxLength  int
	cacheTitle     bool
	cacheBody      bool
}

func newBlogConfig() BlogConfig {
	return BlogConfig{
		titleMaxLength: getInt("TITLE_MAX_LENGTH"),
		bodyMaxLength:  getInt("BODY_MAX_LENGTH"),
		cacheTitle:     getBool("CACHE_TITLE"),
		cacheBody:      getBool("CACHE_BODY"),
	}
}

func (bc BlogConfig) TitleMaxLength() int {
	return bc.titleMaxLength
}

func (bc BlogConfig) BodyMaxLength() int {
	return bc.bodyMaxLength
}
func (bc BlogConfig) CacheTitle() bool {
	return bc.cacheTitle
}

func (bc BlogConfig) CacheBody() bool {
	return bc.cacheBody
}
