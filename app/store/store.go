package store

import (
	"crypto/md5"
	"fmt"
	"sync"

	"github.com/google/uuid"

	generated "github.com/shovanmaity/grpc-task/gen/go"
)

// New returns a new instance of Store
func New() *Store {
	st := &Store{
		mu:       sync.Mutex{},
		profiles: make(map[string]*generated.ProfileMessage),
		tokens:   make(map[string]string),
		sessions: make(map[string]string),
	}
	return st
}

// In memory store for the application
type Store struct {
	mu       sync.Mutex
	profiles map[string]*generated.ProfileMessage
	tokens   map[string]string
	sessions map[string]string
}

// InsertProfile adds a new profile in data store
func (s *Store) InsertProfile(profile *generated.ProfileMessage) (*generated.ProfileMessage, error) {
	if profile == nil {
		return nil, fmt.Errorf("got nil object")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.profiles[profile.Username] = profile
	return profile, nil
}

// UpdateProfile updates a profile in the store
func (s *Store) UpdateProfile(id string, profile *generated.ProfileMessage) (*generated.ProfileMessage, error) {
	if profile == nil {
		return nil, fmt.Errorf("got nil object")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.profiles[id]; !ok {
		return nil, fmt.Errorf("profile not found")
	}
	s.profiles[id] = profile
	return profile, nil
}

// GetProfile returns a profile from the store
func (s *Store) GetProfile(id string) (*generated.ProfileMessage, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	profile, ok := s.profiles[id]
	if !ok {
		return nil, fmt.Errorf("profile not found")
	}
	return profile, nil
}

// InsertToken adds a token in the store
func (s *Store) InsertToken(cred *generated.CredentialMessage) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.tokens[cred.Username]; ok {
		return fmt.Errorf("token already set for %s", cred.Username)
	}
	s.tokens[cred.Username] = string(md5.New().
		Sum([]byte(cred.Username + ":" + cred.Password)))
	return nil
}

// GetToken gets a token from the store
func (s *Store) GetToken(username string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.tokens[username]; !ok {
		return "", fmt.Errorf("%s not present in database", username)
	}
	return s.tokens[username], nil
}

// SetAndGetSession sets a new session if not and returns sessionID for a user
func (s *Store) SetAndGetSession(username string) string {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.sessions[username] == "" {
		s.sessions[username] = uuid.New().String()
	}
	return s.sessions[username]
}

// RemoveSession removes sessionID for given user from the store
func (s *Store) RemoveSession(username string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.sessions, username)
}

// IsValidSession checks if session and username is  valid or not
func (s *Store) IsValidSession(username, session string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.sessions[username] == session
}
