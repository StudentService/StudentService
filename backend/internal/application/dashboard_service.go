package application

import (
	"backend/internal/domain/activity"
	"context"
	"time"

	"github.com/google/uuid"

	"backend/internal/domain/calendar"
	"backend/internal/domain/challenge"
	"backend/internal/domain/course"
	"backend/internal/domain/dashboard"
	"backend/internal/domain/grade"
	"backend/internal/domain/group"
	"backend/internal/domain/questionnaire"
	"backend/internal/domain/semester"
	"backend/internal/domain/user"
)

type DashboardService struct {
	userRepo          user.Repository
	calendarRepo      calendar.Repository
	challengeRepo     challenge.Repository
	gradeRepo         grade.Repository
	questionnaireRepo questionnaire.Repository
	activityRepo      activity.Repository
	groupRepo         group.Repository
	courseRepo        course.Repository
	semesterRepo      semester.Repository
}

func NewDashboardService(
	userRepo user.Repository,
	calendarRepo calendar.Repository,
	challengeRepo challenge.Repository,
	gradeRepo grade.Repository,
	questionnaireRepo questionnaire.Repository,
	activityRepo activity.Repository,
	groupRepo group.Repository,
	courseRepo course.Repository,
	semesterRepo semester.Repository,
) *DashboardService {
	return &DashboardService{
		userRepo:          userRepo,
		calendarRepo:      calendarRepo,
		challengeRepo:     challengeRepo,
		gradeRepo:         gradeRepo,
		questionnaireRepo: questionnaireRepo,
		activityRepo:      activityRepo,
		groupRepo:         groupRepo,
		courseRepo:        courseRepo,
		semesterRepo:      semesterRepo,
	}
}

// GetStudentDashboard собирает все данные для дашборда студента
func (s *DashboardService) GetStudentDashboard(ctx context.Context, userID uuid.UUID) (*dashboard.DashboardResponse, error) {
	// Получаем информацию о студенте
	studentInfo, err := s.getStudentInfo(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Получаем ближайшие события
	upcomingEvents, err := s.getUpcomingEvents(ctx, userID)
	if err != nil {
		upcomingEvents = []*dashboard.UpcomingEvent{}
	}

	// Получаем активные вызовы
	activeChallenges, err := s.getActiveChallenges(ctx, userID)
	if err != nil {
		activeChallenges = []*dashboard.ActiveChallenge{}
	}

	// Получаем последние оценки
	recentGrades, err := s.getRecentGrades(ctx, userID)
	if err != nil {
		recentGrades = []*dashboard.RecentGrade{}
	}

	// Получаем статус анкеты
	questionnaireStatus, err := s.getQuestionnaireStatus(ctx, userID)
	if err != nil {
		questionnaireStatus = &dashboard.QuestionnaireStatus{Status: "not_found"}
	}

	// Получаем статистику
	statistics, err := s.getStatistics(ctx, userID)
	if err != nil {
		statistics = &dashboard.Statistics{}
	}

	return &dashboard.DashboardResponse{
		StudentInfo:         studentInfo,
		UpcomingEvents:      upcomingEvents,
		ActiveChallenges:    activeChallenges,
		RecentGrades:        recentGrades,
		QuestionnaireStatus: questionnaireStatus,
		Statistics:          statistics,
	}, nil
}

// getStudentInfo собирает информацию о студенте
func (s *DashboardService) getStudentInfo(ctx context.Context, userID uuid.UUID) (*dashboard.StudentInfo, error) {
	u, err := s.userRepo.GetByID(ctx, userID.String())
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, nil
	}

	info := &dashboard.StudentInfo{
		ID:   u.ID,
		Name: u.GetFullName(),
	}

	// Получаем информацию о группе
	if u.GroupID != nil {
		group, err := s.groupRepo.GetByID(ctx, *u.GroupID)
		if err == nil && group != nil {
			info.GroupName = group.Name

			// Получаем информацию о курсе
			course, err := s.courseRepo.GetByID(ctx, group.CourseID)
			if err == nil && course != nil {
				info.CourseName = course.Name
			}

			// Получаем информацию о семестре
			semester, err := s.semesterRepo.GetByID(ctx, group.SemesterID)
			if err == nil && semester != nil {
				info.Semester = semester.Name
			}
		}
	}

	return info, nil
}

// getUpcomingEvents получает ближайшие события (максимум 5)
func (s *DashboardService) getUpcomingEvents(ctx context.Context, userID uuid.UUID) ([]*dashboard.UpcomingEvent, error) {
	now := time.Now()
	weekLater := now.AddDate(0, 0, 7)

	events, err := s.calendarRepo.GetStudentEvents(ctx, userID, now, weekLater)
	if err != nil {
		return nil, err
	}

	// Ограничиваем до 5 событий
	if len(events) > 5 {
		events = events[:5]
	}

	result := make([]*dashboard.UpcomingEvent, len(events))
	for i, e := range events {
		isToday := e.StartTime.Year() == now.Year() &&
			e.StartTime.YearDay() == now.YearDay()
		isTomorrow := e.StartTime.Year() == now.Year() &&
			e.StartTime.YearDay() == now.YearDay()+1

		result[i] = &dashboard.UpcomingEvent{
			ID:         e.ID,
			Title:      e.Title,
			Type:       string(e.Type),
			StartTime:  e.StartTime,
			EndTime:    e.EndTime,
			Location:   e.Location,
			OnlineLink: e.OnlineLink,
			IsToday:    isToday,
			IsTomorrow: isTomorrow,
		}
	}

	return result, nil
}

// getActiveChallenges получает активные вызовы
func (s *DashboardService) getActiveChallenges(ctx context.Context, userID uuid.UUID) ([]*dashboard.ActiveChallenge, error) {
	challenges, err := s.challengeRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	var active []*dashboard.ActiveChallenge
	now := time.Now()

	for _, c := range challenges {
		// Только активные или в процессе
		if c.Status == challenge.StatusActive || c.Status == challenge.StatusDraft {
			daysLeft := int(c.EndDate.Sub(now).Hours() / 24)
			if daysLeft < 0 {
				daysLeft = 0
			}

			active = append(active, &dashboard.ActiveChallenge{
				ID:       c.ID,
				Title:    c.Title,
				Goal:     c.Goal,
				Progress: c.Progress,
				DaysLeft: daysLeft,
				EndDate:  c.EndDate,
				Status:   string(c.Status),
			})
		}
	}

	// Ограничиваем до 5
	if len(active) > 5 {
		active = active[:5]
	}

	return active, nil
}

// getRecentGrades получает последние оценки (максимум 5)
func (s *DashboardService) getRecentGrades(ctx context.Context, userID uuid.UUID) ([]*dashboard.RecentGrade, error) {
	grades, err := s.gradeRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Сортируем по дате (самые свежие первые)
	if len(grades) > 5 {
		grades = grades[:5]
	}

	result := make([]*dashboard.RecentGrade, len(grades))
	for i, g := range grades {
		// Получаем название курса
		courseName := "Неизвестный курс"
		course, err := s.courseRepo.GetByID(ctx, g.CourseID)
		if err == nil && course != nil {
			courseName = course.Name
		}

		result[i] = &dashboard.RecentGrade{
			ID:         g.ID,
			CourseName: courseName,
			Type:       string(g.Type),
			Value:      g.Value,
			MaxValue:   g.MaxValue,
			Date:       g.Date,
			IsPassed:   g.Value >= (g.MaxValue * 0.6), // 60% от максимума
		}
	}

	return result, nil
}

// getQuestionnaireStatus получает статус анкеты
func (s *DashboardService) getQuestionnaireStatus(ctx context.Context, userID uuid.UUID) (*dashboard.QuestionnaireStatus, error) {
	q, err := s.questionnaireRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if q == nil {
		return &dashboard.QuestionnaireStatus{Status: "not_found"}, nil
	}

	return &dashboard.QuestionnaireStatus{
		Status:      string(q.Status),
		SubmittedAt: q.SubmittedAt,
		ReviewedAt:  q.ReviewedAt,
		Comment:     q.Comment,
	}, nil
}

// getStatistics собирает общую статистику
func (s *DashboardService) getStatistics(ctx context.Context, userID uuid.UUID) (*dashboard.Statistics, error) {
	stats := &dashboard.Statistics{}

	// Статистика по вызовам
	challenges, err := s.challengeRepo.GetByUserID(ctx, userID)
	if err == nil {
		stats.TotalChallenges = len(challenges)
		for _, c := range challenges {
			if c.Status == challenge.StatusCompleted {
				stats.CompletedChallenges++
			}
		}
	}

	// Средняя оценка
	avg, err := s.gradeRepo.GetAverageByUser(ctx, userID)
	if err == nil {
		stats.AverageGrade = avg
	}

	// Количество предстоящих событий
	now := time.Now()
	monthLater := now.AddDate(0, 1, 0)
	events, err := s.calendarRepo.GetStudentEvents(ctx, userID, now, monthLater)
	if err == nil {
		stats.UpcomingEvents = len(events)
	}

	// Сумма баллов (из активностей) - теперь работает!
	participations, err := s.activityRepo.GetMyParticipations(ctx, userID)
	if err == nil {
		total := 0
		for _, p := range participations {
			total += p.PointsEarned
		}
		stats.TotalPoints = total
	}

	return stats, nil
}
