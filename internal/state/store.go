package state

import (
	"errors"
	"fmt"
	"sync"
)

// Store is a threadsafe runtime variable container.
type Store struct {
	mu   sync.RWMutex
	vars map[string]interface{}
}

func New() *Store {
	return &Store{vars: map[string]interface{}{}}
}

var Global = New()

func (s *Store) Get(name string) (interface{}, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	value, ok := s.vars[name]
	return value, ok
}

func (s *Store) MustGet(name string) interface{} {
	value, ok := s.Get(name)
	if !ok {
		panic(fmt.Sprintf("state variable not found: %s", name))
	}
	return value
}

func (s *Store) Set(name string, value interface{}) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.vars[name] = value
}

// Update mimics XSStrike updateVar behavior.
// mode "append": append string data to []string variable.
// mode "add": add string data into map[string]struct{} set.
// any other mode means direct Set.
func (s *Store) Update(name string, data interface{}, mode string) error {
	if mode == "" {
		s.Set(name, data)
		return nil
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	switch mode {
	case "append":
		slice, ok := s.vars[name].([]string)
		if !ok {
			return errors.New("append mode requires []string target")
		}
		entry, ok := data.(string)
		if !ok {
			return errors.New("append mode requires string data")
		}
		s.vars[name] = append(slice, entry)
		return nil
	case "add":
		set, ok := s.vars[name].(map[string]struct{})
		if !ok {
			return errors.New("add mode requires map[string]struct{} target")
		}
		entry, ok := data.(string)
		if !ok {
			return errors.New("add mode requires string data")
		}
		set[entry] = struct{}{}
		s.vars[name] = set
		return nil
	default:
		s.vars[name] = data
		return nil
	}
}
