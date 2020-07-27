package store

type Store struct {
	ss StoriesStore
	sc StoriesCache
}

func NewStore(sc StoriesCache, ss StoriesStore) *Store {
	return &Store{
		sc: sc,
		ss: ss,
	}
}

func (s *Store) GetStoriesStore() StoriesStore {
	return s.ss
}

func (s *Store) GetStoriesCache() StoriesCache {
	return s.sc
}
