package store

import (
	"crypto/md5"
	"fmt"
	"sync"

	"github.com/google/uuid"

	generated "github.com/shovanmaity/grpc-task/gen/go"
)

func NewStore() *Store {
	st := &Store{
		mu:       sync.Mutex{},
		profiles: make(map[string]*generated.ProfileMessage),
		tokens:   make(map[string]string),
		sessions: make(map[string]string),
	}
	st.InsertProfile(&generated.ProfileMessage{Username: "shovan", Name: "Kuchbhi", Email: "test@gmail.com"})
	return st
}

type Store struct {
	mu       sync.Mutex
	profiles map[string]*generated.ProfileMessage
	tokens   map[string]string
	sessions map[string]string
}

func (s *Store) InsertProfile(profile *generated.ProfileMessage) (*generated.ProfileMessage, error) {
	if profile == nil {
		return nil, fmt.Errorf("got nil object")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.profiles[profile.Username] = profile
	return profile, nil
}

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

func (s *Store) GetProfile(id string) (*generated.ProfileMessage, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	profile, ok := s.profiles[id]
	if !ok {
		return nil, fmt.Errorf("profile not found")
	}
	return profile, nil
}

func (s *Store) InsertToken(cred *generated.CredentialMessage) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.tokens[cred.Username]; ok {
		return fmt.Errorf("token already set for %s", cred.Username)
	}
	s.tokens[cred.Username] = string(md5.New().
		Sum([]byte(cred.Username + ":" + cred.Password)))
	fmt.Println(s.tokens[cred.Username])
	return nil
}

func (s *Store) GetToken(username string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.tokens[username]; !ok {
		return "", fmt.Errorf("%s not present in database", username)
	}
	return s.tokens[username], nil
}

func (s *Store) GetSession(username string) string {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.sessions[username] == "" {
		s.sessions[username] = uuid.New().String()
	}
	return s.sessions[username]
}

func (s *Store) RemoveSession(username string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.sessions, username)
}

func (s *Store) IsValidSession(username, session string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.sessions[username] == session
}
