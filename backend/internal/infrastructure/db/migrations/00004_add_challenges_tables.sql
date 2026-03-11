-- +goose Up
-- +goose StatementBegin
-- Статусы вызовов
CREATE TYPE challenge_status AS ENUM (
'draft', 'active', 'completed', 'overdue'
);

-- Таблица вызовов
CREATE TABLE challenges (
id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
title VARCHAR(255) NOT NULL,
description TEXT,
goal TEXT NOT NULL,
start_date TIMESTAMP NOT NULL,
end_date TIMESTAMP NOT NULL,
status challenge_status NOT NULL DEFAULT 'active',
progress INT NOT NULL DEFAULT 0 CHECK (progress >= 0 AND progress <= 100),
created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Таблица чекпоинтов
CREATE TABLE checkpoints (
id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
challenge_id UUID NOT NULL REFERENCES challenges(id) ON DELETE CASCADE,
title VARCHAR(255) NOT NULL,
description TEXT,
due_date TIMESTAMP NOT NULL,
is_completed BOOLEAN NOT NULL DEFAULT false,
completed_at TIMESTAMP,
order_num INT NOT NULL DEFAULT 0,
created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Таблица артефактов
CREATE TABLE artifacts (
id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
challenge_id UUID NOT NULL REFERENCES challenges(id) ON DELETE CASCADE,
type VARCHAR(50) NOT NULL CHECK (type IN ('file', 'link')),
name VARCHAR(255) NOT NULL,
url TEXT NOT NULL,
created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Таблица самооценок
CREATE TABLE self_assessments (
id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
challenge_id UUID UNIQUE NOT NULL REFERENCES challenges(id) ON DELETE CASCADE,
rating INT NOT NULL CHECK (rating >= 1 AND rating <= 10),
comment TEXT,
created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Индексы для оптимизации запросов
CREATE INDEX idx_challenges_user_id ON challenges(user_id);
CREATE INDEX idx_challenges_status ON challenges(status);
CREATE INDEX idx_challenges_dates ON challenges(start_date, end_date);
CREATE INDEX idx_checkpoints_challenge_id ON checkpoints(challenge_id);
CREATE INDEX idx_checkpoints_due_date ON checkpoints(due_date);
CREATE INDEX idx_artifacts_challenge_id ON artifacts(challenge_id);
CREATE INDEX idx_self_assessments_challenge_id ON self_assessments(challenge_id);

-- Комментарии к таблицам
COMMENT ON TABLE challenges IS 'Личные вызовы студентов';
COMMENT ON TABLE checkpoints IS 'Чекпоинты (этапы) выполнения вызова';
COMMENT ON TABLE artifacts IS 'Артефакты (файлы/ссылки) прикрепленные к вызову';
COMMENT ON TABLE self_assessments IS 'Самооценка студента по завершению вызова';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS self_assessments;
DROP TABLE IF EXISTS artifacts;
DROP TABLE IF EXISTS checkpoints;
DROP TABLE IF EXISTS challenges;
DROP TYPE IF EXISTS challenge_status;
-- +goose StatementEnd