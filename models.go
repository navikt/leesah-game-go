package leesah

type (
	MessageType string
)

const (
	MessageTypeQuestion MessageType = "SPØRSMÅL"
	MessageTypeAnswer   MessageType = "SVAR"

	// LeesahTimeformat is Python's time format for Leesah messages, which is a form of RFC3339
	LeesahTimeformat = "2006-01-02T15:04:05.999999"
)

// Question is a simplified version of a Message with just the category and question
type Question struct {
	Category string
	Question string
}

// MinimalMessage is a simplified version of a Message with just the type
type MinimalMessage struct {
	Type MessageType `json:"type"`
}

// Message is a struct that represents a message in the Leesah Kafka topic
type Message struct {
	Answer     string      `json:"svar"`
	Category   string      `json:"kategorinavn"`
	Created    string      `json:"@opprettet"`
	MessageID  string      `json:"svarId"`
	Question   string      `json:"spørsmål,omitempty"`
	QuestionID string      `json:"spørsmålId"`
	TeamName   string      `json:"lagnavn"`
	Type       MessageType `json:"@event_name"`
}

// ToQuestion converts a Message to a simpler Question
func (m Message) ToQuestion() Question {
	return Question{
		Category: m.Category,
		Question: m.Question,
	}
}
