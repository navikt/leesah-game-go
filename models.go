package leesah

import (
	"github.com/google/uuid"
)

type (
	MessageType      string
	AssessmentStatus string
)

const (
	MessageTypeQuestion   MessageType = "QUESTION"
	MessageTypeAnswer                 = "ANSWER"
	MessageTypeAssessment             = "ASSESSMENT"

	AssessmentStatusSuccess AssessmentStatus = "SUCCESS"
	AssessmentStatusFailure                  = "FAILURE"

	LeesahTimeformat = "2006-01-02T15:04:05.999999"
)

type Message struct {
	Answer     string           `json:"answer,omitempty"`
	AnswerID   *uuid.UUID       `json:"answerId,omitempty"`
	MessageID  uuid.UUID        `json:"messageId,omitempty"`
	Question   string           `json:"question,omitempty"`
	QuestionID uuid.UUID        `json:"questionId,omitempty"`
	Category   string           `json:"category,omitempty"`
	Created    string           `json:"created,omitempty"`
	TeamName   string           `json:"teamName,omitempty"`
	Type       MessageType      `json:"type,omitempty"`
	Status     AssessmentStatus `json:"status,omitempty"`
}

func (m Message) ToQuestion() Question {
	return Question{
		Category: m.Category,
		Question: m.Question,
	}
}

func (m Message) ToAssessment() Assessment {
	return Assessment{
		Category: m.Category,
		Status:   m.Status,
	}
}

type Assessment struct {
	Category string
	Status   AssessmentStatus
	TeamName string
}

type Question struct {
	Category string
	Question string
}
