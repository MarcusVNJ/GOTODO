package repository_impl

import (
	"context"
	"errors"
	"github.com/MarcusVNJ/GOTODO/internal/adapters/out/entity"
	"github.com/MarcusVNJ/GOTODO/internal/adapters/out/mappers"
	"github.com/MarcusVNJ/GOTODO/internal/adapters/out/repository/query_builder"
	"github.com/MarcusVNJ/GOTODO/internal/core/enums"
	"github.com/MarcusVNJ/GOTODO/internal/core/exceptions"
	"github.com/MarcusVNJ/GOTODO/internal/core/exceptions/codes"
	"github.com/MarcusVNJ/GOTODO/internal/core/models"
	"github.com/MarcusVNJ/GOTODO/internal/core/ports"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/samber/oops"
)

type PostgresTaskRepository struct {
	db      *pgxpool.Pool
	builder *task_query_builder.TaskQueryBuilder
}

func NewPostgresTaskRepository(db *pgxpool.Pool) ports.TaskRepository {
	return &PostgresTaskRepository{
		db:      db,
		builder: task_query_builder.NewTaskQueryBuilder(),
	}
}

func (repository *PostgresTaskRepository) Save(context context.Context, task *models.Task) error {
	taskEntity := mappers.DomainToEntity(task)

	query, args, err := repository.builder.QueryInsert(taskEntity)
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

func (repository *PostgresTaskRepository) ExistByID(context context.Context, id string) (bool, error) {
	query, args, err := repository.builder.QueryExistsById(id)
	if err != nil {
		return false, oops.
			In("PostgresTaskRepository").
			Tags("database", "postgres").
			Wrapf(err, "falha crítica ao tentar criar o sql de verificação se existe ou não uma task")
	}
	var exist int

	err = repository.db.QueryRow(context, query, args...).Scan(&exist)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, oops.
			In("PostgresTaskRepository").
			Tags("database", "postgres").
			Wrapf(err, "falha crítica ao tentar executar o sql de verificação se existe ou não uma task")
	}
	return true, nil
}

func (repository *PostgresTaskRepository) FindByID(context context.Context, id string) (*models.Task, error) {

	query, args, err := repository.builder.QueryFindById(id)
	if err != nil {
		return nil, oops.
			In("PostgresTaskRepository").
			Tags("database", "postgres").
			Wrapf(err, "falha crítica ao tentar criar sql de busca de task")
	}

	taskRow := repository.db.QueryRow(context, query, args...)

	return scanTask(taskRow)
}

func (repository *PostgresTaskRepository) FindAll(context context.Context, statusFilter string, minPriority int) ([]*models.Task, error) {
	query, args, err := repository.builder.QueryFindAllTasks(statusFilter, minPriority)
	if err != nil {
		return nil, oops.
			In("PostgresTaskRepository").
			Tags("database", "postgres").
			Wrapf(err, "falha ao criar sql para buscar todas as tasks")
	}

	tasksRows, err := repository.db.Query(context, query, args...)
	if err != nil {
		return nil, oops.
			In("PostgresTaskRepository").
			Tags("database", "postgres").
			Wrapf(err, "falha crítica ao executar a query de busca de todas as tasks")
	}

	defer tasksRows.Close()

	tasksModel, err := scanTasks(tasksRows)
	if err != nil {
		return nil, oops.
			In("PostgresTaskRepository").
			Tags("database", "postgres").
			Wrapf(err, "falha ao escanear resultados das tasks")
	}

	return tasksModel, nil

}

func (repository *PostgresTaskRepository) Update(context context.Context, task *models.Task) error {
	taskEntity := mappers.DomainToEntity(task)

	query, args, err := repository.builder.QueryUpdate(taskEntity)
	if err != nil {
		return oops.
			In("PostgresTaskRepository").
			Tags("database", "postgres").
			Wrapf(err, "falha ao criar sql para atualizar task")
	}

	_, err = repository.db.Exec(context, query, args...)
	if err != nil {
		return oops.
			In("PostgresTaskRepository").
			Tags("database", "postgres").
			Wrapf(err, "falha crítica ao executar atualização no banco")
	}

	return nil
}

func (repository *PostgresTaskRepository) Delete(context context.Context, id string) error {

	query, args, err := repository.builder.QueryDelete(id)
	if err != nil {
		return oops.
			In("PostgresTaskRepository").
			Tags("database", "postgres").
			Wrapf(err, "falha crítica ao tentar criar o sql de exclusão de uma task")
	}

	_, err = repository.db.Exec(context, query, args...)
	if err != nil {
		return oops.
			In("PostgresTaskRepository").
			Tags("database", "postgres").
			Wrapf(err, "falha crítica ao tentar executar o sql de exclusão de uma task")
	}

	return nil
}

func (repository *PostgresTaskRepository) FindByStatus(context context.Context, status enums.Status) ([]*models.Task, error) {

	query, args, err := repository.builder.QueryFindByStatus(status)
	if err != nil {
		return nil, oops.
			In("PostgresTaskRepository").
			Tags("database", "postgres").
			Wrapf(err, "falha ao criar sql para buscar tasks por status")
	}

	tasksRows, err := repository.db.Query(context, query, args...)
	if err != nil {
		return nil, oops.
			In("PostgresTaskRepository").
			Tags("database", "postgres").
			Wrapf(err, "falha crítica ao executar a query de busca por status")
	}

	defer tasksRows.Close()

	tasksModel, err := scanTasks(tasksRows)
	if err != nil {
		return nil, oops.
			In("PostgresTaskRepository").
			Tags("database", "postgres").
			Wrapf(err, "falha ao escanear resultados das tasks por status")
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
			return nil, oops.
				In("PostgresTaskRepository").
				Tags("database", "postgres").
				Wrapf(err, "falha crítica ao tentar tyransformar a task em uma entidade")
		}

		taskDomain := mappers.EntityToDomain(&taskEntity)

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
			return nil, exceptions.NewBusinessException(codes.TaskNotFound)
		}
		return nil, oops.
			In("PostgresTaskRepository").
			Tags("database", "postgres").
			Wrapf(err, "falha crítica ao tentar executar o sql de busca da task")
	}

	return mappers.EntityToDomain(&taskEntity), nil
}
