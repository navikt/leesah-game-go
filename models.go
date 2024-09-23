package leesah

// MessageType is an enum for the type of message in the Leesah Kafka topic
type MessageType string

const (
	MessageTypeQuestion MessageType = "SPØRSMÅL"
	MessageTypeAnswer               = "SVAR"
)

const (
	// LeesahTimeformat is Python's time format for Leesah messages, which is a form of RFC3339
	LeesahTimeformat = "2006-01-02T15:04:05.999999"
)

// Question is a simplified version of a Message with just the category,
// question, and answer format
type Question struct {
	ID            string
	Category      string
	Question      string
	AnswerFormat  string
	Documentation string
}

// minimalMessage is a simplified version of a Message with just the type
// This is used to determine if the message is a question that needs to be answered
type minimalMessage struct {
	Type MessageType `json:"@event_name"`
}

// Message is a struct that represents a message in the Leesah Kafka topic
type Message struct {
	Answer        string      `json:"svar"`
	Category      string      `json:"kategorinavn"`
	Created       string      `json:"@opprettet"`
	AnswerID      string      `json:"svarId"`
	Question      string      `json:"spørsmål,omitempty"`
	QuestionID    string      `json:"spørsmålId"`
	AnswerFormat  string      `json:"svarformat,omitempty"`
	TeamName      string      `json:"lagnavn"`
	Type          MessageType `json:"@event_name"`
	Documentation string      `json:"dokumentasjon"`
}

// ToQuestion converts a Message to a simpler Question model
func (m Message) ToQuestion() Question {
	return Question{
		ID:            m.QuestionID,
		Category:      m.Category,
		Question:      m.Question,
		AnswerFormat:  m.AnswerFormat,
		Documentation: m.Documentation,
	}
}
