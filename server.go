package main

import (
	"encoding/json"
	"net/http"
)

type Server struct {
	shortener *URLShortener
}

type shortenRequest struct {
	URL string `json:"url"`
}

type shortenResponse struct {
	Code string `json:"code"`
}

type errorResponse struct {
	Error string `json:"error"`
}

func NewServer(shortener *URLShortener) *Server {
	return &Server{
		shortener: shortener,
	}
}

func (s *Server) routes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/shorten", s.handleShorten)
	mux.HandleFunc("/", s.handleResolve)

	return mux
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func (s *Server) handleShorten(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, errorResponse{Error: "method not allowed"})
		return
	}

	var req shortenRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil || req.URL == "" {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid Request"})
		return
	}

	code, err := s.shortener.Shorten(req.URL)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: err.Error()})
		return
	}

	writeJSON(w, http.StatusCreated, shortenResponse{Code: code})
}

func (s *Server) handleResolve(w http.ResponseWriter, r *http.Request) {
	// Allow only GET
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Ignore /shorten
	if r.URL.Path == "/shorten" {
		http.NotFound(w, r)
		return
	}

	// Extract code
	code := r.URL.Path[1:]

	// If empty → health check
	if code == "" {
		w.Write([]byte("URL Shortener is running"))
		return
	}

	// Resolve
	longURL, exists := s.shortener.Resolve(code)
	if !exists {
		http.NotFound(w, r)
		return
	}

	// Redirect
	http.Redirect(w, r, longURL, http.StatusFound)
}
