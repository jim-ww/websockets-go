package handler

type InputType string

const (
	ChatMessage InputType = "chat_message"
)

type RequestPayload struct {
	Type    InputType `json:"type"`
	Content string    `json:"content"`
}
