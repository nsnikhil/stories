package config

type StoryConfig struct {
	titleMaxLength int
	bodyMaxLength  int
}

func newStoryConfig() StoryConfig {
	return StoryConfig{
		titleMaxLength: getInt("TITLE_MAX_LENGTH"),
		bodyMaxLength:  getInt("BODY_MAX_LENGTH"),
	}
}

func (bc StoryConfig) TitleMaxLength() int {
	return bc.titleMaxLength
}

func (bc StoryConfig) BodyMaxLength() int {
	return bc.bodyMaxLength
}
