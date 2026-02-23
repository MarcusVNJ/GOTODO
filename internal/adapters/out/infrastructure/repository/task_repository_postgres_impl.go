package repository_impl

import (
	"context"
	"errors"
    "github.com/MarcusVNJ/GOTODO/internal/adapters/out/infrastructure/entity"
    "github.com/MarcusVNJ/GOTODO/internal/adapters/out/infrastructure/mappers"
    "github.com/MarcusVNJ/GOTODO/internal/adapters/out/infrastructure/repository/query_builder"
    "github.com/MarcusVNJ/GOTODO/internal/core/enums"
	"github.com/MarcusVNJ/GOTODO/internal/core/models"
    "github.com/MarcusVNJ/GOTODO/internal/core/ports"
    "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/xid"
	"github.com/samber/oops"
)

var (
	ErrTaskNotFound = errors.New("task not found") //TODO: esse ë um erro de negocio, o usecase q tem q lidar com isso
)

type PostgresTaskRepository struct {
	db      *pgxpool.Pool
	builder *task_query_builder.TaskQueryBuilder
}

func NewPostgresTaskRepository(db *pgxpool.Pool) repository.TaskRepository {
	return &PostgresTaskRepository{
		db:      db,
		builder: task_query_builder.NewTaskQueryBuilder(),
	}
}

func (repository *PostgresTaskRepository) Save(context context.Context, request *models.Task) error {
	task := mappers.DomainToEntity(request)

	query, args, err := repository.builder.QueryInsert(task)
	if err != nil {
		return oops.
			In("PostgresTaskRepository").
			Tags("database", "sql").
			Wrapf(err, "falha crítica ao tentar inserir task no banco")
	}

	_, err = repository.db.Exec(context, query, args...)
	if err != nil {
		return oops.
			In("PostgresTaskRepository").
			Tags("database", "postgres").
			Wrapf(err, "falha crítica ao tentar inserir task no banco")
	}

	return nil
}

func (repository *PostgresTaskRepository) FindByID(context context.Context, id xid.ID) (*models.Task, error) {

	query, args, err := repository.builder.QueryFindById(id.String())
	if err != nil {
		return nil, err
	}

	taskRow := repository.db.QueryRow(context, query, args...)

	return scanTask(taskRow)
}

func (repository *PostgresTaskRepository) FindAll(context context.Context, statusFilter string, minPriority int) ([]*models.Task, error) {
	query, args, err := repository.builder.QueryFindAllTasks(statusFilter, minPriority)
	if err != nil {
		return nil, err
	}

	tasksRows, err := repository.db.Query(context, query, args...)
	if err != nil {
		return nil, err
	}

	defer tasksRows.Close()

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

	_, err = repository.db.Exec(context, query, args...)
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

	_, err = repository.db.Exec(context, query, args...)
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

	tasksRows, err := repository.db.Query(context, query, args...)
	if err != nil {
		return nil, err
	}

	defer tasksRows.Close()

	tasksModel, err := scanTasks(tasksRows)
	if err != nil {
		return nil, err
	}

	return tasksModel, nil
}

func scanTasks(tasksRows pgx.Rows) ([]*models.Task, error) {
	var tasksModel []*models.Task

	for tasksRows.Next() {
		var taskEntity entity.TaskEntity

		err := tasksRows.Scan(
			&taskEntity.ID, &taskEntity.Title, &taskEntity.Description, &taskEntity.Status,
			&taskEntity.Priority, &taskEntity.CreatedAt, &taskEntity.UpdatedAt, &taskEntity.DeletedAt,
		)

		if err != nil {
			return nil, err // TODO: lembra da colocar o oops aqui
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

func scanTask(taskRow pgx.Row) (*models.Task, error) {

	var taskEntity entity.TaskEntity

	err := taskRow.Scan(
		&taskEntity.ID, &taskEntity.Title, &taskEntity.Description, &taskEntity.Status,
		&taskEntity.Priority, &taskEntity.CreatedAt, &taskEntity.UpdatedAt, &taskEntity.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrTaskNotFound
		}
		return nil, err
	}

	return mappers.EntityToDomain(&taskEntity)
}
