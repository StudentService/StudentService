-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS casbin_rule (
                                           id SERIAL PRIMARY KEY,
                                           ptype VARCHAR(100),
    v0 VARCHAR(255),
    v1 VARCHAR(255),
    v2 VARCHAR(255),
    v3 VARCHAR(255),
    v4 VARCHAR(255),
    v5 VARCHAR(255),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
    );

CREATE INDEX idx_casbin_ptype ON casbin_rule(ptype);
CREATE INDEX idx_casbin_v0 ON casbin_rule(v0);
CREATE INDEX idx_casbin_v1 ON casbin_rule(v1);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS casbin_rule;
-- +goose StatementEnd