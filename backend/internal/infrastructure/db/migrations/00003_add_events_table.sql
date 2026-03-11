-- +goose Up
-- +goose StatementBegin
CREATE TYPE event_type AS ENUM (
    'class', 'meeting', 'deadline', 'activity', 'exam', 'holiday'
);

CREATE TABLE events (
                        id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                        title VARCHAR(255) NOT NULL,
                        description TEXT,
                        type event_type NOT NULL,
                        start_time TIMESTAMP NOT NULL,
                        end_time TIMESTAMP NOT NULL,
                        all_day BOOLEAN DEFAULT false,
                        location VARCHAR(255),
                        online_link VARCHAR(500),

    -- Связи
                        course_id UUID REFERENCES courses(id) ON DELETE CASCADE,
                        group_id UUID REFERENCES groups(id) ON DELETE CASCADE,
                        user_id UUID REFERENCES users(id) ON DELETE CASCADE,
                        created_by UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                        created_by_role VARCHAR(50) NOT NULL, -- кто создал (admin/teacher/holder)

                        created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                        updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Индексы
CREATE INDEX idx_events_start_time ON events(start_time);
CREATE INDEX idx_events_end_time ON events(end_time);
CREATE INDEX idx_events_course_id ON events(course_id);
CREATE INDEX idx_events_group_id ON events(group_id);
CREATE INDEX idx_events_user_id ON events(user_id);
CREATE INDEX idx_events_created_by ON events(created_by);
CREATE INDEX idx_events_type ON events(type);
CREATE INDEX idx_events_date_range ON events(start_time, end_time);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS events;
DROP TYPE IF EXISTS event_type;
-- +goose StatementEnd