package llm

import "errors"

var (
	ErrOverMaxQuestionLength = errors.New("maximum question length exceeded")
	ErrOverMaxAnswerLength   = errors.New("maximum answer length exceeded")
	ErrOverMaxTextLength     = errors.New("maximum text length exceeded")
	ErrOverMaxSequenceTimes  = errors.New("maximum number of sequence exceeded")
)
