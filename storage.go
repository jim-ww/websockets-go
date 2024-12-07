package main

import (
	"sync"
)

type Message struct {
	Client *Client
	Text   string
}

func NewMessage(client *Client, text string) *Message {
	return &Message{
		Client: client,
		Text:   text,
	}
}

type Storage struct {
	messages      []*Message
	notifications []string
	mu            sync.Mutex
}

func NewStorage() *Storage {
	return &Storage{
		messages:      []*Message{},
		notifications: []string{},
	}
}

func (s *Storage) GetMessages() []*Message {
	return s.messages
}

func (s *Storage) GetNotifications() []string {
	return s.notifications
}

func (s *Storage) AddMessage(m *Message) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.messages = append(s.messages, m)
}

func (s *Storage) AddNotification(notification string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.notifications = append(s.notifications, notification)
}
