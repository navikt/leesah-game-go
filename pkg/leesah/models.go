package leesah

import (
	"github.com/google/uuid"
)

type (
	message_type      string
	assessment_status string
)

const (
	message_type_question   message_type = "QUESTION"
	message_type_answer                  = "ANSWER"
	message_type_assessment              = "ASSESSMENT"

	assessment_status_success assessment_status = "SUCCESS"
	assessment_status_failure                   = "FAILURE"
)

type LeesahMessage struct {
	Answer     string            `json:"answer"`
	AnswerID   uuid.UUID         `json:"answerId"`
	MessageID  uuid.UUID         `json:"messageId"`
	QuestionID uuid.UUID         `json:"questionId"`
	Category   string            `json:"category"`
	Created    string            `json:"created"`
	TeamName   string            `json:"teamName"`
	Type       message_type      `json:"type"`
	Status     assessment_status `json:"status"`
	Sign       string            `json:"sign"`
}

func (m LeesahMessage) ToQuestion() Question {
	return Question{
		Category:  m.Category,
		MessageID: m.MessageID,
		Question:  m.Sign,
	}
}

func (m LeesahMessage) ToAssessment() Assessment {
	return Assessment{
		AnswerID:   m.AnswerID,
		MessageID:  m.MessageID,
		QuestionID: m.QuestionID,
		Category:   m.Category,
		Status:     m.Status,
		Sign:       m.Sign,
	}
}

type Answer struct {
	Answer     string
	Category   string
	MessageID  uuid.UUID
	QuestionID uuid.UUID
}

type Assessment struct {
	AnswerID   uuid.UUID
	MessageID  uuid.UUID
	QuestionID uuid.UUID
	Category   string
	Status     assessment_status
	Sign       string
}

type Question struct {
	Category  string
	MessageID uuid.UUID
	Question  string
}
