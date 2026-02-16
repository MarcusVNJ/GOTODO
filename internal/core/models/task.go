package models

import (
	"errors"
	"github.com/MarcusVNJ/GOTODO/internal/core/enums"
)

type Task struct {
	Audit
	title       string
	description string
	status      enums.Status
	priority    int
}

var (
	ErrTaskTitleEmpty  = errors.New("task title cannot be empty")
	ErrTaskAlreadyDone = errors.New("task is already completed")
	ErrTaskInitDone    = errors.New("task cannot start with a completed state")
	ErrInvalidPriority = errors.New("priority must be between 1 and 5")
)

func NewTask(title, description string, priority int) (*Task, error) {
	task := &Task{
		Audit:       NewAudit(),
		title:       title,
		description: description,
		priority:    priority,
		status:      enums.Pending,
	}

	if err := task.validade(); err != nil {
		return nil, err
	}

	return task, nil
}

func NewTaskInit(audit Audit, title, desc string, status enums.Status, priority int) *Task {
	return &Task{
		Audit:       audit,
		title:       title,
		description: desc,
		status:      status,
		priority:    priority,
	}
}

func (task *Task) validade() error {
	if err := task.validateTitle(); err != nil {
		return err
	}

	if err := task.validatePriority(); err != nil {
		return err
	}

	if err := task.verifyStatusIsCompleted(); err != nil {
		return ErrTaskInitDone
	}
	return nil
}

func (task *Task) validateTitle() error {
	if task.title == "" {
		return ErrTaskTitleEmpty
	}
	return nil
}

func (task *Task) validatePriority() error {
	if task.priority < 1 || task.priority > 5 {
		return ErrInvalidPriority
	}
	return nil
}

func (task *Task) verifyStatusIsCompleted() error {
	if task.status == enums.Completed {
		return ErrTaskAlreadyDone
	}
	return nil
}

func (task *Task) Complete() error {
	if err := task.verifyStatusIsCompleted(); err != nil {
		return err
	}

	task.status = enums.Completed
	task.UpdatedAudit()
	return nil
}

func (t *Task) Title() string        { return t.title }
func (t *Task) Description() string  { return t.description }
func (t *Task) Status() enums.Status { return t.status }
func (t *Task) Priority() int        { return t.priority }
