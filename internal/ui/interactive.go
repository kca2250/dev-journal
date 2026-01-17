package ui

import (
	"errors"
	"fmt"

	"github.com/charmbracelet/huh"
	"github.com/kca2250/djou/internal/model"
)

// Action represents the user's selected action
type Action string

const (
	ActionEdit   Action = "edit"
	ActionDelete Action = "delete"
	ActionCancel Action = "cancel"
)

// LogSelector displays logs for selection
type LogSelector struct {
	localizer *Localizer
	logs      []model.Log
}

// NewLogSelector creates a new LogSelector
func NewLogSelector(l *Localizer, logs []model.Log) *LogSelector {
	return &LogSelector{
		localizer: l,
		logs:      logs,
	}
}

// formatLogOption formats a log entry for display in the select menu
func formatLogOption(log model.Log) string {
	return fmt.Sprintf("%s | %s | Est: %.1fh Act: %.1fh",
		log.CreatedAt.Format("2006-01-02"),
		TruncateString(log.TaskName, 30),
		log.EstimateHours,
		log.ActualHours,
	)
}

// Run displays the log selection form and returns the selected log
// Returns nil if the user selects "Exit"
func (s *LogSelector) Run() (*model.Log, error) {
	if len(s.logs) == 0 {
		return nil, nil
	}

	// Build options including exit option
	options := make([]huh.Option[int], 0, len(s.logs)+1)

	// Add exit option first
	options = append(options, huh.NewOption(s.localizer.Get(MsgExitInteractive), -1))

	// Add log options
	for i, log := range s.logs {
		options = append(options, huh.NewOption(formatLogOption(log), i))
	}

	var selected int
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[int]().
				Title(s.localizer.Get(MsgSelectLog)).
				Options(options...).
				Value(&selected),
		),
	).WithKeyMap(NewKeyMapWithEsc())

	err := form.Run()
	if err != nil {
		if errors.Is(err, huh.ErrUserAborted) {
			return nil, nil
		}
		return nil, err
	}

	if selected == -1 {
		return nil, nil
	}

	return &s.logs[selected], nil
}

// ActionSelector displays action options for a selected log
type ActionSelector struct {
	localizer *Localizer
	log       model.Log
}

// NewActionSelector creates a new ActionSelector
func NewActionSelector(l *Localizer, log model.Log) *ActionSelector {
	return &ActionSelector{
		localizer: l,
		log:       log,
	}
}

// Run displays the action selection form
func (a *ActionSelector) Run() (Action, error) {
	options := []huh.Option[Action]{
		huh.NewOption(a.localizer.Get(MsgActionEdit), ActionEdit),
		huh.NewOption(a.localizer.Get(MsgActionDelete), ActionDelete),
		huh.NewOption(a.localizer.Get(MsgActionCancel), ActionCancel),
	}

	var selected Action
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[Action]().
				Title(a.localizer.Get(MsgSelectAction)).
				Options(options...).
				Value(&selected),
		),
	).WithKeyMap(NewKeyMapWithEsc())

	err := form.Run()
	if err != nil {
		if errors.Is(err, huh.ErrUserAborted) {
			return ActionCancel, nil
		}
		return ActionCancel, err
	}

	return selected, nil
}

// DeleteConfirm shows a delete confirmation dialog
type DeleteConfirm struct {
	localizer *Localizer
	log       model.Log
}

// NewDeleteConfirm creates a new DeleteConfirm
func NewDeleteConfirm(l *Localizer, log model.Log) *DeleteConfirm {
	return &DeleteConfirm{
		localizer: l,
		log:       log,
	}
}

// Run displays the delete confirmation form
func (d *DeleteConfirm) Run() (bool, error) {
	var confirmed bool
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title(d.localizer.Get(MsgDeleteConfirm)).
				Description(formatLogOption(d.log)).
				Value(&confirmed),
		),
	).WithKeyMap(NewKeyMapWithEsc())

	err := form.Run()
	if err != nil {
		if errors.Is(err, huh.ErrUserAborted) {
			return false, nil
		}
		return false, err
	}

	return confirmed, nil
}

// EditForm creates and runs the edit form with pre-populated values
type EditForm struct {
	localizer *Localizer
	log       model.Log
	input     *FormInput
}

// NewEditForm creates a new EditForm with pre-populated values
func NewEditForm(l *Localizer, log model.Log) *EditForm {
	input := &FormInput{
		TaskName:      log.TaskName,
		EstimateHours: fmt.Sprintf("%.1f", log.EstimateHours),
		ActualHours:   fmt.Sprintf("%.1f", log.ActualHours),
	}
	if log.Memo != nil {
		input.Memo = *log.Memo
	}

	return &EditForm{
		localizer: l,
		log:       log,
		input:     input,
	}
}

// Run executes the edit form and returns the updated input
func (e *EditForm) Run() (*model.LogInput, error) {
	l := e.localizer

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title(l.Get(MsgTaskNameLabel)).
				Value(&e.input.TaskName).
				Validate(ValidateTaskName),

			huh.NewInput().
				Title(l.Get(MsgEstimateHoursLabel)).
				Value(&e.input.EstimateHours).
				Validate(ValidatePositiveFloat),

			huh.NewInput().
				Title(l.Get(MsgActualHoursLabel)).
				Value(&e.input.ActualHours).
				Validate(ValidatePositiveFloat),

			huh.NewText().
				Title(l.Get(MsgMemoLabel)).
				Value(&e.input.Memo),
		),
	).WithKeyMap(NewKeyMapWithEsc())

	err := form.Run()
	if err != nil {
		if errors.Is(err, huh.ErrUserAborted) {
			return nil, ErrFormCancelled
		}
		return nil, err
	}

	return e.input.ToLogInput()
}
