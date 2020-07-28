package config

type BlogConfig struct {
	titleMaxLength int
	bodyMaxLength  int
}

func newBlogConfig() BlogConfig {
	return BlogConfig{
		titleMaxLength: getInt("TITLE_MAX_LENGTH"),
		bodyMaxLength:  getInt("BODY_MAX_LENGTH"),
	}
}

func (bc BlogConfig) TitleMaxLength() int {
	return bc.titleMaxLength
}

func (bc BlogConfig) BodyMaxLength() int {
	return bc.bodyMaxLength
}
