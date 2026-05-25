package models

import (
	"github.com/MarcusVNJ/GOTODO/internal/core/enums"
	"github.com/MarcusVNJ/GOTODO/internal/core/exceptions"
	"github.com/MarcusVNJ/GOTODO/internal/core/exceptions/codes"
)

type Task struct {
	Audit
	title       string
	description string
	status      enums.Status
	priority    int
}

func NewTaskWithoutAudit(title, description string, priority int) (*Task, error) {
	task := &Task{
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
	return new(Task{
		Audit:       audit,
		title:       title,
		description: desc,
		status:      status,
		priority:    priority,
	})
}

func (task *Task) validade() error {
	if err := task.validateTitle(); err != nil {
		return err
	}

	if err := task.validatePriority(); err != nil {
		return err
	}

	if err := task.verifyStatusIsCompleted(); err != nil {
		return exceptions.NewBusinessException(codes.TaskInitDone)
	}
	return nil
}

func (task *Task) validateTitle() error {
	if task.title == "" {
		return exceptions.NewBusinessException(codes.TaskTitleEmpty)
	}
	return nil
}

func (task *Task) validatePriority() error {
	if task.priority < 1 || task.priority > 5 {
		return exceptions.NewBusinessException(codes.InvalidPriority)
	}
	return nil
}

func (task *Task) verifyStatusIsCompleted() error {
	if task.status == enums.Completed {
		return exceptions.NewBusinessException(codes.TaskAlreadyDone)
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

func (task *Task) Title() string        { return task.title }
func (task *Task) Description() string  { return task.description }
func (task *Task) Status() enums.Status { return task.status }
func (task *Task) Priority() int        { return task.priority }
