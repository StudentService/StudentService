-- +goose Up
-- +goose StatementBegin
-- Типы оценок
CREATE TYPE grade_type AS ENUM (
    'exam', 'test', 'homework', 'project', 'activity'
);

-- Таблица оценок
CREATE TABLE grades (
                        id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                        user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                        course_id UUID NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
                        type grade_type NOT NULL,
                        value FLOAT NOT NULL CHECK (value >= 0),
                        max_value FLOAT NOT NULL CHECK (max_value > 0),
                        weight FLOAT NOT NULL DEFAULT 1.0 CHECK (weight >= 0),
                        comment TEXT,
                        date TIMESTAMP NOT NULL,

    -- Метаданные источника
                        source_type VARCHAR(50) NOT NULL DEFAULT 'manual',
                        source_id UUID,

                        created_by UUID NOT NULL REFERENCES users(id),
                        created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                        updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Индексы
CREATE INDEX idx_grades_user_id ON grades(user_id);
CREATE INDEX idx_grades_course_id ON grades(course_id);
CREATE INDEX idx_grades_date ON grades(date);
CREATE INDEX idx_grades_user_course ON grades(user_id, course_id);
CREATE INDEX idx_grades_created_by ON grades(created_by);

-- Представление для сводки по студенту (опционально)
CREATE VIEW student_grade_summary AS
SELECT
    user_id,
    course_id,
    COUNT(*) as grades_count,
    AVG(value / max_value * 100) as average_percentage,
    SUM(weight) as total_weight
FROM grades
GROUP BY user_id, course_id;

-- Комментарии
COMMENT ON TABLE grades IS 'Оценки студентов';
COMMENT ON COLUMN grades.value IS 'Полученная оценка';
COMMENT ON COLUMN grades.max_value IS 'Максимальная оценка';
COMMENT ON COLUMN grades.weight IS 'Вес оценки для расчёта среднего';
COMMENT ON COLUMN grades.source_type IS 'Источник: manual, import, system';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP VIEW IF EXISTS student_grade_summary;
DROP TABLE IF EXISTS grades;
DROP TYPE IF EXISTS grade_type;
-- +goose StatementEnd