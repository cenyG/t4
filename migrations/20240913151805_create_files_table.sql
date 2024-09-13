-- +goose Up
-- +goose StatementBegin
CREATE TABLE files (
     id SERIAL PRIMARY KEY,
     name TEXT NOT NULL UNIQUE,
     servers TEXT NOT NULL,  -- list of server numbers divided by comma (like '1,2,3')
     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE files
-- +goose StatementEnd
