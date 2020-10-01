package config

type StoryConfig struct {
	titleMaxLength int
	bodyMaxLength  int
	//cacheTitle     bool
	//cacheBody      bool
}

func newStoryConfig() StoryConfig {
	return StoryConfig{
		titleMaxLength: getInt("TITLE_MAX_LENGTH"),
		bodyMaxLength:  getInt("BODY_MAX_LENGTH"),
		//cacheTitle:     getBool("CACHE_TITLE"),
		//cacheBody:      getBool("CACHE_BODY"),
	}
}

func (bc StoryConfig) TitleMaxLength() int {
	return bc.titleMaxLength
}

func (bc StoryConfig) BodyMaxLength() int {
	return bc.bodyMaxLength
}

//func (bc StoryConfig) CacheTitle() bool {
//	return bc.cacheTitle
//}
//
//func (bc StoryConfig) CacheBody() bool {
//	return bc.cacheBody
//}
