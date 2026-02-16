package task_query_builder

import (
    "github.com/MarcusVNJ/GOTODO/internal/core/enums"
    "github.com/MarcusVNJ/GOTODO/internal/infrastructure/entity"
    sq "github.com/Masterminds/squirrel"
    "time"
)

type TaskQueryBuilder struct {
    psql sq.StatementBuilderType
}

func NewTaskQueryBuilder() *TaskQueryBuilder {
    return &TaskQueryBuilder{
        psql: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
    }
}

func (query_builder *TaskQueryBuilder) QueryInsert(entity entity.TaskEntity) (string, []interface{}, error) {
    return query_builder.psql.Insert("task").
        Columns("id", "title", "description", "status", "priority", "created_at", "updated_at", "deleted_at").
        Values(
            entity.ID,
            entity.Title,
            entity.Description,
            entity.Status,
            entity.Priority,
            entity.CreatedAt,
            entity.UpdatedAt,
            entity.DeletedAt,
        ).ToSql()
}

func (query_builder *TaskQueryBuilder) QueryFindById(id string) (string, []interface{}, error) {
    return query_builder.psql.Select("*").
        From("tasks").
        Where(sq.Eq{"id": id}).
        ToSql()
}

func (query_builder *TaskQueryBuilder) QueryFindAllTasks(statusFilter string, minPriority int) (string, []interface{}, error) {
    queryBuilder := query_builder.psql.Select("*").From("Task")

    if statusFilter != "" {
        queryBuilder.Where(sq.Eq{"status": statusFilter})
    }
    if minPriority > 0 {
        queryBuilder.Where(sq.GtOrEq{"priority": minPriority})
    }

    queryBuilder.OrderBy("created_at DESC")

    return queryBuilder.ToSql()
}

func (query_builder *TaskQueryBuilder) QueryUpdate(entity entity.TaskEntity) (string, []interface{}, error) {
    return query_builder.psql.Update("task").
        Set("title", entity.Title).
        Set("description", entity.Description).
        Set("status", entity.Status).
        Set("priority", entity.Priority).
        Set("updated_at", entity.UpdatedAt).
        Set("deleted_at", entity.DeletedAt).
        Where(sq.Eq{"id": entity.ID}).
        ToSql()
}

func (query_builder *TaskQueryBuilder) QueryDelete(id string) (string, []interface{}, error) {
    return query_builder.psql.Update("task").
        Set("deleted_at", time.Now()).
        Where(sq.Eq{"id": id}).
        ToSql()
}

func (query_builder *TaskQueryBuilder) QueryFindByStatus(status enums.Status) (string, []interface{}, error) {
    return query_builder.psql.Select("*").
        From("tasks").
        Where(sq.Eq{"status": status}).
        ToSql()
}
