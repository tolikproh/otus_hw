-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE TABLE IF NOT EXISTS event (
                                     id serial PRIMARY KEY,
                                     title TEXT NOT NULL,
                                     date_time_start timestamptz NOT NULL,
                                     date_time_end timestamptz NOT NULL,
                                     description TEXT,
                                     user_id int NOT NULL,
                                     notification_duration bigint
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE IF EXISTS event;
-- +goose StatementEnd
