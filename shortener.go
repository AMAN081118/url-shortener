package main

import (
	"errors"
	"net/url"
	"sync/atomic"
)

const base62Alphabet = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

var ErrEmptyURL = errors.New("url is empty")
var ErrInvalidURL = errors.New("url is invalid")

type URLShortener struct {
	storage   Storage
	idCounter uint64
}

func NewURLShortener(storage Storage) *URLShortener {
	return &URLShortener{
		storage: storage,
	}
}

func (s *URLShortener) Shorten(longURL string) (string, error) {

	// validate empty URL
	if len(longURL) == 0 {
		return "", ErrEmptyURL
	}
	// validate format
	if !isValidURL(longURL) {
		return "", ErrInvalidURL
	}

	// check if URL already exists in storage
	existingCode, exists := s.storage.GetCodeByURL(longURL)

	// if exists, return existingCode, nil
	if exists {
		return existingCode, nil
	}

	// generate next ID
	id := atomic.AddUint64(&s.idCounter, 1) - 1

	// encode ID to base62
	code := encodeBase62(id)

	// save into storage
	s.storage.Save(code, longURL)

	return code, nil

}

func (s *URLShortener) Resolve(code string) (string, bool) {

	// delegate to storage
	resolvedURL, exists := s.storage.GetURLByCode(code)

	if exists {
		return resolvedURL, true
	}
	return "", false
}

func (s *URLShortener) Count() int {
	// TODO:
	// delegate to storage
	count := s.storage.Count()
	return count
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
