-- +goose Up
-- +goose StatementBegin
-- Типы активностей
CREATE TYPE activity_type AS ENUM (
    'class', 'workshop', 'meeting', 'task', 'project', 'event'
);

-- Статусы активностей
CREATE TYPE activity_status AS ENUM (
    'active', 'completed', 'cancelled', 'draft'
);

-- Статусы участия
CREATE TYPE participation_status AS ENUM (
    'enrolled', 'attended', 'completed', 'missed', 'cancelled'
);

-- Таблица активностей
CREATE TABLE activities (
                            id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                            title VARCHAR(255) NOT NULL,
                            description TEXT,
                            type activity_type NOT NULL,
                            status activity_status NOT NULL DEFAULT 'draft',

    -- Временные параметры
                            start_time TIMESTAMP,
                            end_time TIMESTAMP,
                            deadline TIMESTAMP,

    -- Место проведения
                            location VARCHAR(255),
                            online_link VARCHAR(500),

    -- Ограничения
                            max_participants INT CHECK (max_participants > 0),
                            current_participants INT NOT NULL DEFAULT 0 CHECK (current_participants >= 0),

    -- Баллы
                            points INT NOT NULL DEFAULT 0,
                            weight FLOAT NOT NULL DEFAULT 1.0,

    -- Для кого активность
                            course_id UUID REFERENCES courses(id) ON DELETE CASCADE,
                            group_id UUID REFERENCES groups(id) ON DELETE CASCADE,
                            created_by UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                            created_by_role VARCHAR(50) NOT NULL,

                            created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                            updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    -- Проверка дат
                            CONSTRAINT check_dates CHECK (
                                (start_time IS NULL OR end_time IS NULL OR start_time <= end_time)
                                )
);

-- Таблица участий
CREATE TABLE participations (
                                id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                                activity_id UUID NOT NULL REFERENCES activities(id) ON DELETE CASCADE,
                                user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                                status participation_status NOT NULL DEFAULT 'enrolled',

    -- Результаты
                                grade FLOAT,
                                feedback TEXT,
                                points_earned INT NOT NULL DEFAULT 0,

    -- Временные метки
                                enrolled_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                completed_at TIMESTAMP,

                                created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    -- Уникальность - студент может быть записан на активность только один раз
                                UNIQUE(activity_id, user_id)
);

-- Индексы
CREATE INDEX idx_activities_type ON activities(type);
CREATE INDEX idx_activities_status ON activities(status);
CREATE INDEX idx_activities_dates ON activities(start_time, end_time);
CREATE INDEX idx_activities_course_id ON activities(course_id);
CREATE INDEX idx_activities_group_id ON activities(group_id);
CREATE INDEX idx_activities_created_by ON activities(created_by);

CREATE INDEX idx_participations_activity_id ON participations(activity_id);
CREATE INDEX idx_participations_user_id ON participations(user_id);
CREATE INDEX idx_participations_status ON participations(status);
CREATE INDEX idx_participations_enrolled_at ON participations(enrolled_at);

-- Комментарии
COMMENT ON TABLE activities IS 'Каталог активностей (занятия, мероприятия, задачи)';
COMMENT ON TABLE participations IS 'Участия студентов в активностях';
COMMENT ON COLUMN activities.points IS 'Баллы за участие';
COMMENT ON COLUMN activities.weight IS 'Вес для расчёта метрик';
COMMENT ON COLUMN participations.grade IS 'Оценка за активность';
COMMENT ON COLUMN participations.points_earned IS 'Фактически полученные баллы';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS participations;
DROP TABLE IF EXISTS activities;
DROP TYPE IF EXISTS participation_status;
DROP TYPE IF EXISTS activity_status;
DROP TYPE IF EXISTS activity_type;
-- +goose StatementEnd