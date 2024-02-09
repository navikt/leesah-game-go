package leesah

type (
	MessageType string
)

const (
	MessageTypeQuestion MessageType = "QUESTION"
	MessageTypeAnswer               = "ANSWER"

	LeesahTimeformat = "2006-01-02T15:04:05.999999"
)

type Question struct {
	Category string
	Question string
}

type MinimalMessage struct {
	Type MessageType `json:"type"`
}

type Message struct {
	Answer     string      `json:"answer"`
	Category   string      `json:"category"`
	Created    string      `json:"created"`
	MessageID  string      `json:"messageId"`
	Question   string      `json:"question,omitempty"`
	QuestionID string      `json:"questionId"`
	TeamName   string      `json:"teamName"`
	Type       MessageType `json:"type"`
}

func (m Message) ToQuestion() Question {
	return Question{
		Category: m.Category,
		Question: m.Question,
	}
}
