package main

import "sync"

type Storage interface {
	GetCodeByURL(longURL string) (string, bool)

	GetURLByCode(code string) (string, bool)

	GetOrCreate(longURL string, generate func() string) string

	Save(code, longURL string)

	Count() int
}

type InMemoryStorage struct {
	mu        sync.RWMutex
	codeToURL map[string]string
	urlToCode map[string]string
}

func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		codeToURL: make(map[string]string),
		urlToCode: make(map[string]string),
	}
}

func (s *InMemoryStorage) GetCodeByURL(longURL string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	code, exists := s.urlToCode[longURL]
	if exists {
		return code, true
	}

	return "", false
}

func (s *InMemoryStorage) GetOrCreate(longURL string, generate func() string) string {
	s.mu.Lock()
	defer s.mu.Unlock()

	if code, exists := s.urlToCode[longURL]; exists {
		return code
	}

	code := generate()

	s.urlToCode[longURL] = code
	s.codeToURL[code] = longURL

	return code
}

func (s *InMemoryStorage) GetURLByCode(code string) (string, bool) {

	s.mu.RLock()
	defer s.mu.RUnlock()
	longURL, exists := s.codeToURL[code]
	if exists {
		return longURL, true
	}
	return "", false
}

func (s *InMemoryStorage) Save(code, longURL string) {

	s.mu.Lock()
	defer s.mu.Unlock()
	s.urlToCode[longURL] = code
	s.codeToURL[code] = longURL
}

func (s *InMemoryStorage) Count() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.codeToURL)
}
