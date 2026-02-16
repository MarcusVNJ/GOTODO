CREATE TYPE task_status AS ENUM (
    'PENDING',
    'IN_PROCESS',
    'COMPLETED',
    'CANCELLED'
);

CREATE TABLE IF NOT EXISTS tasks (
    id          CHAR(20) PRIMARY KEY,
    title       VARCHAR(255) NOT NULL,
    description TEXT,
    status      task_status NOT NULL,
    priority    INT NOT NULL,
    created_at  TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at  TIMESTAMP WITH TIME ZONE NOT NULL,
    deleted_at  TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_tasks_status ON tasks(status);