package main

import (
	"errors"
	"net/url"
	"sync"
)

const base62Alphabet = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

var ErrEmptyURL = errors.New("url is empty")
var ErrInvalidURL = errors.New("url is invalid")

type URLShortener struct {
	mu        sync.RWMutex
	idCounter uint64
	codeToURL map[string]string
	urlToCode map[string]string
}

func NewURLShortener() *URLShortener {
	return &URLShortener{
		codeToURL: make(map[string]string),
		urlToCode: make(map[string]string),
	}
}

func (s *URLShortener) Shorten(longURL string) (string, error) {
	if longURL == "" {
		return "", ErrEmptyURL
	}

	if !isValidURL(longURL) {
		return "", ErrInvalidURL
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if code, exists := s.urlToCode[longURL]; exists {
		return code, nil
	}

	code := encodeBase62(s.idCounter)
	s.idCounter++

	s.codeToURL[code] = longURL
	s.urlToCode[longURL] = code

	return code, nil
}

func (s *URLShortener) Resolve(code string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	longURL, exists := s.codeToURL[code]
	return longURL, exists
}

func (s *URLShortener) Count() int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return len(s.codeToURL)
}

func isValidURL(rawURL string) bool {
	parsedURL, err := url.ParseRequestURI(rawURL)
	if err != nil {
		return false
	}

	return parsedURL.Scheme != "" && parsedURL.Host != ""
}

func encodeBase62(n uint64) string {
	if n == 0 {
		return string(base62Alphabet[0])
	}

	var result []byte

	for n > 0 {
		remainder := n % 62
		result = append(result, base62Alphabet[remainder])
		n = n / 62
	}

	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
	}

	return string(result)
}
