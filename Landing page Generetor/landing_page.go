package main

import (
	"sync"

	"github.com/sashabaranov/go-openai"
)

// LandingPage represents the structure of a landing page
type LandingPage struct {
	HTML    string                         `json:"html"`
	ID      string                         `json:"id"`
	History []openai.ChatCompletionMessage `json:"history"`
}
type LandingPageStore struct {
	pages map[string]LandingPage
	mu    sync.RWMutex
}

func NewLandingPageStore() *LandingPageStore {
	return &LandingPageStore{
		pages: make(map[string]LandingPage),
	}
}

// Get retrieves a LandingPage by ID, returning a boolean for existence
func (s *LandingPageStore) Get(id string) (LandingPage, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	page, exists := s.pages[id]
	return page, exists
}

// Add inserts or updates a LandingPage
func (s *LandingPageStore) Add(page LandingPage) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.pages[page.ID] = page
}

// Delete removes a LandingPage by ID
func (s *LandingPageStore) Delete(id string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.pages, id)
}

// Exists checks if a LandingPage exists by ID
func (s *LandingPageStore) Exists(id string) bool {
	_, exists := s.Get(id)
	return exists
}
