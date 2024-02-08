package leesah

import (
	"github.com/google/uuid"
)

type (
	MessageType string
)

const (
	MessageTypeQuestion MessageType = "QUESTION"
	MessageTypeAnswer               = "ANSWER"

	LeesahTimeformat = "2006-01-02T15:04:05.999999"
)

type Message struct {
	Answer     string      `json:"answer,omitempty"`
	AnswerID   *uuid.UUID  `json:"answerId,omitempty"`
	MessageID  uuid.UUID   `json:"messageId,omitempty"`
	Question   string      `json:"question,omitempty"`
	QuestionID uuid.UUID   `json:"questionId,omitempty"`
	Category   string      `json:"category,omitempty"`
	Created    string      `json:"created,omitempty"`
	TeamName   string      `json:"teamName,omitempty"`
	Type       MessageType `json:"type,omitempty"`
}

func (m Message) ToQuestion() Question {
	return Question{
		Category: m.Category,
		Question: m.Question,
	}
}

type Question struct {
	Category string
	Question string
}
