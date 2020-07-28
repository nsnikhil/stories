package store

type Store struct {
	sc StoriesCache
	ss StoriesStore
}

func NewStore(ss StoriesStore, sc StoriesCache) *Store {
	return &Store{
		ss: ss,
		sc: sc,
	}
}

func (s *Store) GetStoriesStore() StoriesStore {
	return s.ss
}

func (s *Store) GetStoriesCache() StoriesCache {
	return s.sc
}
