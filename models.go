package leesah

import "fmt"

// MessageType is an enum for the type of message in the Leesah Kafka topic.
type MessageType string

// FeedbackType is an enum for the type of feedback in the Leesah Kafka topic.
type FeedbackType string

const (
	MessageTypeQuestion MessageType = "SPØRSMÅL"
	MessageTypeAnswer   MessageType = "SVAR"
	MessageTypeFeedback MessageType = "KORREKTUR"

	FeedbackTypeCorrect   FeedbackType = "KORREKT"
	FeedbackTypeIncorrect FeedbackType = "FEIL"
)

const (
	// LeesahTimeformat is Python's time format for Leesah messages, which is a form of RFC3339.
	LeesahTimeformat = "2006-01-02T15:04:05.999999"
)

// Question is a simplified version of a Message with just the fields you need to play the game.
type Question struct {
	ID            string
	Category      string
	Question      string
	AnswerFormat  string
	Documentation string
}

// FeedbackMessage is a message with just the fields needed when you have received feedback.
type FeedbackMessage struct {
	QuestionID string       `json:"spørsmålId"`
	Feedback   FeedbackType `json:"korrektur"`
	Category   string       `json:"kategori"`
}

// minimalMessage is a simplified version of a Message with just the message type.
// This is used to determine if the message is a question that needs to be answered.
type minimalMessage struct {
	Type     MessageType `json:"@event_name"`
	TeamName string      `json:"lagnavn"`
	Category string      `json:"kategori"`
}

// Message is a struct that represents a message in the Leesah Kafka topic.
type Message struct {
	Answer        string       `json:"svar"`
	Category      string       `json:"kategori"`
	Created       string       `json:"@opprettet"`
	AnswerID      string       `json:"svarId"`
	Question      string       `json:"spørsmål,omitempty"`
	QuestionID    string       `json:"spørsmålId"`
	AnswerFormat  string       `json:"svarformat,omitempty"`
	TeamName      string       `json:"lagnavn"`
	Type          MessageType  `json:"@event_name"`
	Documentation string       `json:"dokumentasjon"`
	Feedback      FeedbackType `json:"korrektur,omitempty"`
}

// ToQuestion converts a Message to a simpler Question model.
func (m Message) ToQuestion() Question {
	return Question{
		ID:            m.QuestionID,
		Category:      m.Category,
		Question:      m.Question,
		AnswerFormat:  m.AnswerFormat,
		Documentation: m.Documentation,
	}
}

type ErrorIncorrectAnswer struct {
	// QuestionID is the ID of the question that was answered incorrectly.
	QuestionID string
	Category   string
}

func NewErrorIncorrectAnswer(questionID, category string) *ErrorIncorrectAnswer {
	return &ErrorIncorrectAnswer{
		QuestionID: questionID,
		Category:   category,
	}
}

func (e *ErrorIncorrectAnswer) Error() string {
	return fmt.Sprintf("Incorrect answer for question %s in category %s", e.QuestionID, e.Category)
}
