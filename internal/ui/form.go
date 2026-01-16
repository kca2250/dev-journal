package ui

import (
	"errors"
	"strconv"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/kca2250/djou/internal/model"
)

// ErrFormCancelled is returned when the user cancels the form
var ErrFormCancelled = errors.New("form cancelled")

// FormInput holds the raw string input from the form
type FormInput struct {
	TaskName      string
	EstimateHours string
	ActualHours   string
	Memo          string
}

// ToLogInput converts FormInput to model.LogInput
func (f *FormInput) ToLogInput() (*model.LogInput, error) {
	estimate, err := ParsePositiveFloat(f.EstimateHours)
	if err != nil {
		return nil, err
	}

	actual, err := ParsePositiveFloat(f.ActualHours)
	if err != nil {
		return nil, err
	}

	return &model.LogInput{
		TaskName:      strings.TrimSpace(f.TaskName),
		EstimateHours: estimate,
		ActualHours:   actual,
		Memo:          strings.TrimSpace(f.Memo),
	}, nil
}

// ParsePositiveFloat parses a string to a positive float64
func ParsePositiveFloat(s string) (float64, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, errors.New("value is required")
	}

	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, errors.New("invalid number")
	}

	if v <= 0 {
		return 0, errors.New("must be positive")
	}

	return v, nil
}

// ValidateTaskName validates the task name field
func ValidateTaskName(s string) error {
	if strings.TrimSpace(s) == "" {
		return errors.New(T(MsgTaskNameRequired))
	}
	return nil
}

// ValidatePositiveFloat validates that a string is a positive float
func ValidatePositiveFloat(s string) error {
	_, err := ParsePositiveFloat(s)
	if err != nil {
		if err.Error() == "value is required" {
			return errors.New(T(MsgPositiveRequired))
		}
		if err.Error() == "invalid number" {
			return errors.New(T(MsgInvalidNumber))
		}
		return errors.New(T(MsgPositiveRequired))
	}
	return nil
}

// RecordForm creates and runs the interactive record form
type RecordForm struct {
	localizer *Localizer
	input     *FormInput
}

// NewRecordForm creates a new RecordForm
func NewRecordForm(l *Localizer) *RecordForm {
	return &RecordForm{
		localizer: l,
		input:     &FormInput{},
	}
}

// Run executes the interactive form and returns the input
func (r *RecordForm) Run() (*model.LogInput, error) {
	l := r.localizer

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title(l.Get(MsgTaskNameLabel)).
				Value(&r.input.TaskName).
				Validate(ValidateTaskName),

			huh.NewInput().
				Title(l.Get(MsgEstimateHoursLabel)).
				Value(&r.input.EstimateHours).
				Validate(ValidatePositiveFloat),

			huh.NewInput().
				Title(l.Get(MsgActualHoursLabel)).
				Value(&r.input.ActualHours).
				Validate(ValidatePositiveFloat),

			huh.NewText().
				Title(l.Get(MsgMemoLabel)).
				Value(&r.input.Memo),
		),
	)

	err := form.Run()
	if err != nil {
		if errors.Is(err, huh.ErrUserAborted) {
			return nil, ErrFormCancelled
		}
		return nil, err
	}

	return r.input.ToLogInput()
}

// GetInput returns the current form input (for testing)
func (r *RecordForm) GetInput() *FormInput {
	return r.input
}
