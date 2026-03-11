-- +goose Up
-- +goose StatementBegin
-- Статусы анкет
CREATE TYPE questionnaire_status AS ENUM (
    'draft', 'submitted', 'approved', 'rejected'
);

-- Таблица анкет студентов
CREATE TABLE questionnaires (
                                id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                                user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                                status questionnaire_status NOT NULL DEFAULT 'draft',
                                answers JSONB NOT NULL DEFAULT '{}',
                                submitted_at TIMESTAMP,
                                reviewed_by UUID REFERENCES users(id) ON DELETE SET NULL,
                                reviewed_at TIMESTAMP,
                                comment TEXT,
                                created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    -- У одного студента может быть только одна анкета
                                CONSTRAINT unique_user_questionnaire UNIQUE (user_id)
);

-- Таблица шаблонов анкет
CREATE TABLE questionnaire_templates (
                                         id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                                         name VARCHAR(255) NOT NULL,
                                         description TEXT,
                                         is_active BOOLEAN DEFAULT false,
                                         schema JSONB NOT NULL, -- JSON Schema для валидации
                                         fields JSONB NOT NULL, -- описание полей для UI
                                         created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                         updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Индексы
CREATE INDEX idx_questionnaires_user_id ON questionnaires(user_id);
CREATE INDEX idx_questionnaires_status ON questionnaires(status);
CREATE INDEX idx_questionnaires_submitted_at ON questionnaires(submitted_at);
CREATE INDEX idx_templates_is_active ON questionnaire_templates(is_active);

-- Добавляем начальный шаблон
INSERT INTO questionnaire_templates (name, description, is_active, schema, fields) VALUES (
                                                                                              'Анкета студента',
                                                                                              'Стандартная анкета для поступающих студентов',
                                                                                              true,
                                                                                              '{
                                                                                                "type": "object",
                                                                                                "required": ["full_name", "birth_date", "phone", "education"],
                                                                                                "properties": {
                                                                                                  "full_name": {"type": "string"},
                                                                                                  "birth_date": {"type": "string", "format": "date"},
                                                                                                  "phone": {"type": "string"},
                                                                                                  "email": {"type": "string", "format": "email"},
                                                                                                  "education": {"type": "string"},
                                                                                                  "interests": {"type": "string"},
                                                                                                  "experience": {"type": "string"},
                                                                                                  "motivation": {"type": "string"}
                                                                                                }
                                                                                              }',
                                                                                              '[
                                                                                                {
                                                                                                  "id": "full_name",
                                                                                                  "type": "text",
                                                                                                  "label": "ФИО",
                                                                                                  "required": true,
                                                                                                  "placeholder": "Иванов Иван Иванович"
                                                                                                },
                                                                                                {
                                                                                                  "id": "birth_date",
                                                                                                  "type": "date",
                                                                                                  "label": "Дата рождения",
                                                                                                  "required": true
                                                                                                },
                                                                                                {
                                                                                                  "id": "phone",
                                                                                                  "type": "text",
                                                                                                  "label": "Телефон",
                                                                                                  "required": true,
                                                                                                  "placeholder": "+7 (999) 123-45-67"
                                                                                                },
                                                                                                {
                                                                                                  "id": "email",
                                                                                                  "type": "text",
                                                                                                  "label": "Email",
                                                                                                  "required": false,
                                                                                                  "placeholder": "ivan@example.com"
                                                                                                },
                                                                                                {
                                                                                                  "id": "education",
                                                                                                  "type": "text",
                                                                                                  "label": "Образование",
                                                                                                  "required": true,
                                                                                                  "placeholder": "ВУЗ, факультет, год окончания"
                                                                                                },
                                                                                                {
                                                                                                  "id": "interests",
                                                                                                  "type": "textarea",
                                                                                                  "label": "Интересы",
                                                                                                  "required": false,
                                                                                                  "placeholder": "Чем увлекаетесь, хобби"
                                                                                                },
                                                                                                {
                                                                                                  "id": "experience",
                                                                                                  "type": "textarea",
                                                                                                  "label": "Опыт работы/проектов",
                                                                                                  "required": false
                                                                                                },
                                                                                                {
                                                                                                  "id": "motivation",
                                                                                                  "type": "textarea",
                                                                                                  "label": "Мотивация",
                                                                                                  "required": false,
                                                                                                  "placeholder": "Почему хотите участвовать в программе"
                                                                                                }
                                                                                              ]'::jsonb
                                                                                          );

-- Комментарии
COMMENT ON TABLE questionnaires IS 'Анкеты студентов';
COMMENT ON TABLE questionnaire_templates IS 'Шаблоны анкет';
COMMENT ON COLUMN questionnaires.answers IS 'Ответы в формате JSON';
COMMENT ON COLUMN questionnaire_templates.schema IS 'JSON Schema для валидации';
COMMENT ON COLUMN questionnaire_templates.fields IS 'Описание полей для UI';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS questionnaires;
DROP TABLE IF EXISTS questionnaire_templates;
DROP TYPE IF EXISTS questionnaire_status;
-- +goose StatementEnd