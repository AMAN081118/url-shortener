package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleShorten_Success(t *testing.T) {
	shortener := NewURLShortener()
	server := NewServer(shortener)

	body := bytes.NewBuffer([]byte(`{"url":"https://google.com"}`))
	req := httptest.NewRequest(http.MethodPost, "/shorten", body)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	server.handleShorten(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("expected %d got %d", http.StatusCreated, rr.Code)
	}

	if rr.Body.Len() == 0 {
		t.Fatalf("expected non-empty response body")
	}
}

func TestHandleShorten_InvalidMethod(t *testing.T) {
	shortener := NewURLShortener()
	server := NewServer(shortener)

	req := httptest.NewRequest(http.MethodGet, "/shorten", nil)
	rr := httptest.NewRecorder()

	server.handleShorten(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected %d got %d", http.StatusMethodNotAllowed, rr.Code)
	}
}

func TestHandleShorten_InvalidJSON(t *testing.T) {
	shortener := NewURLShortener()
	server := NewServer(shortener)

	body := bytes.NewBuffer([]byte(`{"url":}`))
	req := httptest.NewRequest(http.MethodPost, "/shorten", body)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	server.handleShorten(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected %d got %d", http.StatusBadRequest, rr.Code)
	}
}

func TestHandleShorten_InvalidURL(t *testing.T) {
	shortener := NewURLShortener()
	server := NewServer(shortener)

	body := bytes.NewBuffer([]byte(`{"url":"not-real-url"}`))
	req := httptest.NewRequest(http.MethodPost, "/shorten", body)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	server.handleShorten(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected %d got %d", http.StatusBadRequest, rr.Code)
	}
}

func TestHandleResolve_Success(t *testing.T) {
	shortener := NewURLShortener()

	code, err := shortener.Shorten("https://www.google.com")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	server := NewServer(shortener)

	req := httptest.NewRequest(http.MethodGet, "/"+code, nil)
	rr := httptest.NewRecorder()

	server.handleResolve(rr, req)

	if rr.Code != http.StatusFound {
		t.Fatalf("expected %d got %d", http.StatusFound, rr.Code)
	}

	location := rr.Header().Get("Location")
	if location != "https://www.google.com" {
		t.Fatalf("expected redirect to %q, got %q", "https://www.google.com", location)
	}
}

func TestHandleResolve_NotFound(t *testing.T) {
	shortener := NewURLShortener()
	server := NewServer(shortener)

	req := httptest.NewRequest(http.MethodGet, "/unknown", nil)
	rr := httptest.NewRecorder()

	server.handleResolve(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Fatalf("expected %d got %d", http.StatusNotFound, rr.Code)
	}
}

func TestHandleResolve_Root(t *testing.T) {
	shortener := NewURLShortener()
	server := NewServer(shortener)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()

	server.handleResolve(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected %d got %d", http.StatusOK, rr.Code)
	}

	if rr.Body.String() != "URL Shortener is running" {
		t.Fatalf("expected body %q, got %q", "URL Shortener is running", rr.Body.String())
	}
}
