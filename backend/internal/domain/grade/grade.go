package grade

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type GradeType string

const (
	GradeTypeExam     GradeType = "exam"     // экзамен
	GradeTypeTest     GradeType = "test"     // зачёт
	GradeTypeHomework GradeType = "homework" // домашнее задание
	GradeTypeProject  GradeType = "project"  // проект
	GradeTypeActivity GradeType = "activity" // активность на занятии
)

type Grade struct {
	ID       uuid.UUID `json:"id"`
	UserID   uuid.UUID `json:"user_id"`
	CourseID uuid.UUID `json:"course_id"`
	Type     GradeType `json:"type"`
	Value    float64   `json:"value"`     // оценка (может быть 5, 100, 4.5 и т.д.)
	MaxValue float64   `json:"max_value"` // максимальная оценка (5, 100, 10)
	Weight   float64   `json:"weight"`    // вес оценки (для расчёта среднего)
	Comment  string    `json:"comment"`
	Date     time.Time `json:"date"` // дата получения оценки

	// Откуда оценка
	SourceType string     `json:"source_type"`         // "manual", "import", "system"
	SourceID   *uuid.UUID `json:"source_id,omitempty"` // ID импорта или др.

	CreatedBy uuid.UUID `json:"created_by"` // кто выставил (преподаватель)
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CourseSummary сводка по курсу
type CourseSummary struct {
	CourseID     uuid.UUID `json:"course_id"`
	CourseName   string    `json:"course_name"`
	Average      float64   `json:"average"`       // средняя оценка
	TotalCredits int       `json:"total_credits"` // всего кредитов
	GradesCount  int       `json:"grades_count"`  // количество оценок
}

// StudentSummary общая сводка успеваемости студента
type StudentSummary struct {
	OverallAverage float64          `json:"overall_average"` // общий средний балл
	Courses        []*CourseSummary `json:"courses"`
	TotalCredits   int              `json:"total_credits"` // всего кредитов
	LastUpdated    time.Time        `json:"last_updated"`
}

type Repository interface {
	// Для студента
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*Grade, error)
	GetByUserAndCourse(ctx context.Context, userID, courseID uuid.UUID) ([]*Grade, error)
	GetByUserAndPeriod(ctx context.Context, userID uuid.UUID, from, to time.Time) ([]*Grade, error)

	// Для преподавателя/админа
	GetByID(ctx context.Context, id uuid.UUID) (*Grade, error)
	Create(ctx context.Context, grade *Grade) error
	CreateBatch(ctx context.Context, grades []*Grade) error
	Update(ctx context.Context, grade *Grade) error
	Delete(ctx context.Context, id uuid.UUID) error

	// Для расчётов
	GetAverageByUser(ctx context.Context, userID uuid.UUID) (float64, error)
	GetAverageByUserAndCourse(ctx context.Context, userID, courseID uuid.UUID) (float64, error)
	GetSummaryByUser(ctx context.Context, userID uuid.UUID) (*StudentSummary, error)
}
