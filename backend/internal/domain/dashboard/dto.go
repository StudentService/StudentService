package dashboard

import (
	"time"

	"github.com/google/uuid"
)

// DashboardResponse главный DTO для дашборда студента
type DashboardResponse struct {
	StudentInfo         *StudentInfo         `json:"student_info"`
	UpcomingEvents      []*UpcomingEvent     `json:"upcoming_events"`
	ActiveChallenges    []*ActiveChallenge   `json:"active_challenges"`
	RecentGrades        []*RecentGrade       `json:"recent_grades"`
	QuestionnaireStatus *QuestionnaireStatus `json:"questionnaire_status"`
	Statistics          *Statistics          `json:"statistics"`
}

// StudentInfo базовая информация о студенте
type StudentInfo struct {
	ID         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	GroupName  string    `json:"group_name,omitempty"`
	CourseName string    `json:"course_name,omitempty"`
	Semester   string    `json:"semester,omitempty"`
	AvatarURL  string    `json:"avatar_url,omitempty"`
}

// UpcomingEvent ближайшее событие в календаре
type UpcomingEvent struct {
	ID         uuid.UUID `json:"id"`
	Title      string    `json:"title"`
	Type       string    `json:"type"` // class, meeting, deadline, etc
	StartTime  time.Time `json:"start_time"`
	EndTime    time.Time `json:"end_time"`
	Location   string    `json:"location,omitempty"`
	OnlineLink string    `json:"online_link,omitempty"`
	IsToday    bool      `json:"is_today"`
	IsTomorrow bool      `json:"is_tomorrow"`
}

// ActiveChallenge активный личный вызов
type ActiveChallenge struct {
	ID       uuid.UUID `json:"id"`
	Title    string    `json:"title"`
	Goal     string    `json:"goal"`
	Progress int       `json:"progress"` // 0-100
	DaysLeft int       `json:"days_left"`
	EndDate  time.Time `json:"end_date"`
	Status   string    `json:"status"`
}

// RecentGrade последняя оценка
type RecentGrade struct {
	ID         uuid.UUID `json:"id"`
	CourseName string    `json:"course_name"`
	Type       string    `json:"type"` // exam, test, homework, etc
	Value      float64   `json:"value"`
	MaxValue   float64   `json:"max_value"`
	Date       time.Time `json:"date"`
	IsPassed   bool      `json:"is_passed"` // значение >= порога
}

// QuestionnaireStatus статус анкеты
type QuestionnaireStatus struct {
	Status      string     `json:"status"` // draft, submitted, approved, rejected
	SubmittedAt *time.Time `json:"submitted_at,omitempty"`
	ReviewedAt  *time.Time `json:"reviewed_at,omitempty"`
	Comment     string     `json:"comment,omitempty"`
}

// Statistics общая статистика
type Statistics struct {
	TotalChallenges     int     `json:"total_challenges"`
	CompletedChallenges int     `json:"completed_challenges"`
	AverageGrade        float64 `json:"average_grade"`
	TotalPoints         int     `json:"total_points"`
	UpcomingEvents      int     `json:"upcoming_events_count"`
}
