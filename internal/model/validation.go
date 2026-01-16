package model

import (
	"errors"
	"strings"
)

var (
	ErrEmptyTaskName   = errors.New("task name is required")
	ErrInvalidEstimate = errors.New("estimate hours must be positive")
	ErrInvalidActual   = errors.New("actual hours must be positive")
)

// ValidateInput validates the log input
func ValidateInput(input *LogInput) error {
	if strings.TrimSpace(input.TaskName) == "" {
		return ErrEmptyTaskName
	}
	if input.EstimateHours <= 0 {
		return ErrInvalidEstimate
	}
	if input.ActualHours <= 0 {
		return ErrInvalidActual
	}
	return nil
}
