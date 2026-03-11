-- +goose Up
-- +goose StatementBegin
-- Создаём таблицу курсов
CREATE TABLE courses (
                         id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                         name VARCHAR(255) NOT NULL,
                         description TEXT,
                         credits INT NOT NULL DEFAULT 0,
                         created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                         updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Создаём таблицу семестров
CREATE TABLE semesters (
                           id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                           name VARCHAR(255) NOT NULL,
                           start_date DATE NOT NULL,
                           end_date DATE NOT NULL,
                           is_active BOOLEAN DEFAULT false,
                           created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                           updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Создаём таблицу групп
CREATE TABLE groups (
                        id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                        name VARCHAR(255) NOT NULL,
                        course_id UUID NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
                        semester_id UUID NOT NULL REFERENCES semesters(id) ON DELETE CASCADE,
                        holder_id UUID REFERENCES users(id) ON DELETE SET NULL,
                        created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                        updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Добавляем group_id в таблицу users
ALTER TABLE users
    ADD COLUMN group_id UUID REFERENCES groups(id) ON DELETE SET NULL;

-- Создаём индексы
CREATE INDEX idx_courses_name ON courses(name);
CREATE INDEX idx_semesters_is_active ON semesters(is_active);
CREATE INDEX idx_semesters_dates ON semesters(start_date, end_date);
CREATE INDEX idx_groups_course_id ON groups(course_id);
CREATE INDEX idx_groups_semester_id ON groups(semester_id);
CREATE INDEX idx_groups_holder_id ON groups(holder_id);
CREATE INDEX idx_users_group_id ON users(group_id);

-- Добавляем тестовые данные (опционально)
INSERT INTO courses (name, description, credits) VALUES
                                                     ('Программирование на Go', 'Базовый курс по Go', 4),
                                                     ('Веб-разработка', 'HTML, CSS, JavaScript', 3),
                                                     ('Базы данных', 'SQL и проектирование БД', 3);

INSERT INTO semesters (name, start_date, end_date, is_active) VALUES
                                                                  ('Весна 2026', '2026-02-01', '2026-05-31', true),
                                                                  ('Осень 2026', '2026-09-01', '2026-12-31', false);

INSERT INTO groups (name, course_id, semester_id) VALUES
                                                      ('Поток Go-1', (SELECT id FROM courses WHERE name = 'Программирование на Go'),
                                                       (SELECT id FROM semesters WHERE is_active = true)),
                                                      ('Веб-разработка-1', (SELECT id FROM courses WHERE name = 'Веб-разработка'),
                                                       (SELECT id FROM semesters WHERE is_active = true));
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Удаляем в обратном порядке
ALTER TABLE users DROP COLUMN IF EXISTS group_id;
DROP TABLE IF EXISTS groups CASCADE;
DROP TABLE IF EXISTS semesters CASCADE;
DROP TABLE IF EXISTS courses CASCADE;
-- +goose StatementEnd