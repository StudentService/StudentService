package teacher

import (
	"time"

	"github.com/google/uuid"
)

// DashboardResponse дашборд преподавателя
type DashboardResponse struct {
	Groups           []*GroupSummary    `json:"groups"`
	RecentActivities []*ActivitySummary `json:"recent_activities"`
	PendingReviews   int                `json:"pending_reviews"`
	TotalStudents    int                `json:"total_students"`
}

// GroupSummary краткая информация о группе
type GroupSummary struct {
	ID            uuid.UUID `json:"id"`
	Name          string    `json:"name"`
	CourseName    string    `json:"course_name"`
	StudentsCount int       `json:"students_count"`
	AverageGrade  float64   `json:"average_grade"`
}

// ActivitySummary краткая информация об активности
type ActivitySummary struct {
	ID           uuid.UUID `json:"id"`
	Title        string    `json:"title"`
	Type         string    `json:"type"`
	Date         time.Time `json:"date"`
	Participants int       `json:"participants"`
}

// StudentProfile профиль студента для преподавателя
type StudentProfile struct {
	ID           uuid.UUID `json:"id"`
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	Email        string    `json:"email"`
	GroupName    string    `json:"group_name"`
	AverageGrade float64   `json:"average_grade"`
	TotalPoints  int       `json:"total_points"`
}

// GradeWithStudent оценка с данными студента
type GradeWithStudent struct {
	ID          uuid.UUID `json:"id"`
	StudentID   uuid.UUID `json:"student_id"`
	StudentName string    `json:"student_name"`
	Value       float64   `json:"value"`
	MaxValue    float64   `json:"max_value"`
	Type        string    `json:"type"`
	Date        time.Time `json:"date"`
	Comment     string    `json:"comment"`
}

// ImportGradesRequest запрос на импорт оценок
type ImportGradesRequest struct {
	GroupID    uuid.UUID         `json:"group_id" binding:"required"`
	ActivityID uuid.UUID         `json:"activity_id,omitempty"`
	Grades     []GradeImportItem `json:"grades" binding:"required"`
}

// GradeImportItem элемент импорта
type GradeImportItem struct {
	StudentID uuid.UUID `json:"student_id" binding:"required"`
	Value     float64   `json:"value" binding:"required"`
	MaxValue  float64   `json:"max_value" binding:"required"`
	Comment   string    `json:"comment"`
}

// AttendanceRequest запрос на отметку посещаемости
type AttendanceRequest struct {
	StudentIDs []uuid.UUID `json:"student_ids" binding:"required"`
	Attended   bool        `json:"attended"`
}
