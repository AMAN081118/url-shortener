package main

import "testing"

const testURL = "https://google.com"

func TestShortenAndResolve(t *testing.T) {
	shortener := NewURLShortener()

	code, err := shortener.Shorten(testURL)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	resolvedURL, ok := shortener.Resolve(code)
	if !ok {
		t.Fatalf("expected code to resolve, but it did not")
	}

	if resolvedURL != testURL {
		t.Fatalf("expected %q, got %q", testURL, resolvedURL)
	}
}

func TestShortenSameURLReturnsSameCode(t *testing.T) {
	shortener := NewURLShortener()

	code1, err1 := shortener.Shorten(testURL)
	if err1 != nil {
		t.Fatalf("expected no error on first shorten, got %v", err1)
	}

	code2, err2 := shortener.Shorten(testURL)
	if err2 != nil {
		t.Fatalf("expected no error on second shorten, got %v", err2)
	}

	if code1 != code2 {
		t.Fatalf("expected same URL to return same code, got %q and %q", code1, code2)
	}
}

func TestShortenEmptyURL(t *testing.T) {
	shortener := NewURLShortener()

	_, err := shortener.Shorten("")
	if err != ErrEmptyURL {
		t.Fatalf("expected ErrEmptyURL, got %v", err)
	}
}

func TestShortenInvalidURL(t *testing.T) {
	shortener := NewURLShortener()

	_, err := shortener.Shorten("not-a-url")
	if err != ErrInvalidURL {
		t.Fatalf("expected ErrInvalidURL, got %v", err)
	}
}

func TestResolveUnknownCode(t *testing.T) {
	shortener := NewURLShortener()

	_, ok := shortener.Resolve("unknown")
	if ok {
		t.Fatalf("expected unknown code to not resolve")
	}
}

func TestCount(t *testing.T) {
	shortener := NewURLShortener()

	_, err1 := shortener.Shorten(testURL)
	if err1 != nil {
		t.Fatalf("expected no error, got %v", err1)
	}

	_, err2 := shortener.Shorten("https://youtube.com")
	if err2 != nil {
		t.Fatalf("expected no error, got %v", err2)
	}

	_, err3 := shortener.Shorten(testURL)
	if err3 != nil {
		t.Fatalf("expected no error, got %v", err3)
	}

	cnt := shortener.Count()
	if cnt != 2 {
		t.Fatalf("expected count to be 2, got %d", cnt)
	}
}
func TestDifferentURLsReturnDifferentCodes(t *testing.T) {
	shortener := NewURLShortener()

	code1, err := shortener.Shorten("https://google.com")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	code2, err := shortener.Shorten("https://youtube.com")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if code1 == code2 {
		t.Fatalf("expected different URLs to have different codes, got same code %q", code1)
	}
}
func TestInitialCountIsZero(t *testing.T) {
	shortener := NewURLShortener()

	if shortener.Count() != 0 {
		t.Fatalf("expected initial count to be 0, got %d", shortener.Count())
	}
}
