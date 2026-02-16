package repository_impl

import (
	"context"
	"database/sql"
	"errors"
	repository "github.com/MarcusVNJ/GOTODO/internal/adapters/out"
	"github.com/MarcusVNJ/GOTODO/internal/infrastructure/entity"
	"github.com/MarcusVNJ/GOTODO/internal/infrastructure/repository/query_builder"
	"github.com/MarcusVNJ/GOTODO/internal/adapters/mappers"
	"github.com/MarcusVNJ/GOTODO/internal/core/enums"
	"github.com/MarcusVNJ/GOTODO/internal/core/models"
	"github.com/rs/xid"
	"log"
)

var (
	ErrTaskNotFound = errors.New("task not found")
)

type PostgresTaskRepository struct {
	db *sql.DB
	builder *task_query_builder.TaskQueryBuilder
}

func NewPostgresTaskRepository(db *sql.DB) repository.TaskRepository {
	return &PostgresTaskRepository{
		db: db,
		builder: task_query_builder.NewTaskQueryBuilder(),
		}
}

func (repository *PostgresTaskRepository) Save(context context.Context, task entity.TaskEntity) error {

	query, args, err := repository.builder.QueryInsert(task)
	if err != nil {
		return err
	}

	_, err = repository.db.ExecContext(context, query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (repository *PostgresTaskRepository) FindByID(context context.Context, id xid.ID) (*models.Task, error) {

	query, args, err := repository.builder.QueryFindById(id.String())
	if err != nil {
		return nil, err
	}

	taskRow := repository.db.QueryRowContext(context, query, args...)

	return scanTask(taskRow)
}

func (repository *PostgresTaskRepository) FindAll(context context.Context, statusFilter string, minPriority int) ([]*models.Task, error) {
	query, args, err := repository.builder.QueryFindAllTasks(statusFilter, minPriority)
	if err != nil {
		return nil, err
	}

	tasksRows, err := repository.db.QueryContext(context, query, args...)
	if err != nil {
		return nil, err
	}

	defer func() {
	    if err := tasksRows.Close(); err != nil {
	        log.Printf("erro ao fechar rows: %v", err)
	    }
	}()

	tasksModel, err := scanTasks(tasksRows)
	if err != nil {
		return nil, err
	}

	return tasksModel, nil

}

func (repository *PostgresTaskRepository) Update(context context.Context, task entity.TaskEntity) error {

	query, args, err := repository.builder.QueryUpdate(task)
	if err != nil {
		return err
	}

	_, err = repository.db.ExecContext(context, query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (repository *PostgresTaskRepository) Delete(context context.Context, id xid.ID) error {

	query, args, err := repository.builder.QueryDelete(id.String())
	if err != nil {
		return err
	}

	_, err = repository.db.ExecContext(context, query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (repository *PostgresTaskRepository) FindByStatus(context context.Context, status enums.Status) ([]*models.Task, error) {

	query, args, err := repository.builder.QueryFindByStatus(status)
	if err != nil {
		return nil, err
	}

	tasksRows, err := repository.db.QueryContext(context, query, args...)
	if err != nil {
		return nil, err
	}

	defer func() {
	    if err := tasksRows.Close(); err != nil {
	        log.Printf("erro ao fechar rows: %v", err)
	    }
	}()

	tasksModel, err := scanTasks(tasksRows)
	if err != nil {
		return nil, err
	}

	return tasksModel, nil
}

func scanTasks(tasksRows *sql.Rows) ([]*models.Task, error) {
	var tasksModel []*models.Task

	for tasksRows.Next() {
		var taskEntity entity.TaskEntity

		err := tasksRows.Scan(
			&taskEntity.ID, &taskEntity.Title, &taskEntity.Description, &taskEntity.Status,
			&taskEntity.Priority, &taskEntity.CreatedAt, &taskEntity.UpdatedAt, &taskEntity.DeletedAt,
		)

		if err != nil {
			return nil, err
		}

		taskDomain, err := mappers.EntityToDomain(&taskEntity)
		if err != nil {
			return nil, err
		}

		tasksModel = append(tasksModel, taskDomain)
	}

	if err := tasksRows.Err(); err != nil {
		return nil, err
	}

	return tasksModel, nil
}

func scanTask(taskRow *sql.Row) (*models.Task, error) {

	var taskEntity entity.TaskEntity

	err := taskRow.Scan(
		&taskEntity.ID, &taskEntity.Title, &taskEntity.Description, &taskEntity.Status,
		&taskEntity.Priority, &taskEntity.CreatedAt, &taskEntity.UpdatedAt, &taskEntity.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrTaskNotFound
		}
		return nil, err
	}

	return mappers.EntityToDomain(&taskEntity)
}
