package store

import (
	"sync"

	"github.com/google/uuid"
)

type Message struct {
	ClientID uuid.UUID
	Text     string
}

func NewMessage(clientID uuid.UUID, text string) *Message {
	return &Message{
		ClientID: clientID,
		Text:     text,
	}
}

type Store struct {
	messages      []*Message
	notifications []string
	mu            sync.Mutex
}

func NewStore() *Store {
	return &Store{
		messages:      []*Message{},
		notifications: []string{},
	}
}

func (s *Store) GetMessages() []*Message {
	return s.messages
}

func (s *Store) GetNotifications() []string {
	return s.notifications
}

func (s *Store) AddMessage(m *Message) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.messages = append(s.messages, m)
}

func (s *Store) AddNotification(notification string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.notifications = append(s.notifications, notification)
}
