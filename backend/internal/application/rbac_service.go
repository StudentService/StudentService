package application

import (
	"log"

	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type RBACService struct {
	enforcer *casbin.Enforcer
}

func NewRBACService(dbURL string) (*RBACService, error) {
	// Подключаемся к PostgreSQL через GORM
	gormDB, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Адаптер для хранения политик в БД
	adapter, err := gormadapter.NewAdapterByDB(gormDB)
	if err != nil {
		return nil, err
	}

	// Загружаем модель из файла
	enforcer, err := casbin.NewEnforcer("pkg/configs/rbac_model.conf", adapter)
	if err != nil {
		return nil, err
	}

	// Загружаем политики из БД
	if err := enforcer.LoadPolicy(); err != nil {
		return nil, err
	}

	// Инициализируем базовые политики
	if err := initPolicies(enforcer); err != nil {
		log.Printf("Warning: failed to init policies: %v", err)
	}

	return &RBACService{enforcer: enforcer}, nil
}

// CheckPermission проверяет права для роли
func (s *RBACService) CheckPermission(role, resource, action string) (bool, error) {
	return s.enforcer.Enforce(role, resource, action)
}

// CheckUserPermission проверяет права для роли (без ошибки)
func (s *RBACService) CheckUserPermission(userRole, resource, action string) bool {
	allowed, err := s.enforcer.Enforce(userRole, resource, action)
	if err != nil {
		log.Printf("RBAC check error: %v", err)
		return false
	}
	return allowed
}

// GetRolePermissions получает все разрешения для роли
func (s *RBACService) GetRolePermissions(role string) [][]string {
	return s.enforcer.GetPermissionsForUser(role)
}

// AddPolicy добавляет разрешение для роли
func (s *RBACService) AddPolicy(role, resource, action string) error {
	_, err := s.enforcer.AddPolicy(role, resource, action)
	if err == nil {
		s.enforcer.SavePolicy()
	}
	return err
}

// RemovePolicy удаляет разрешение
func (s *RBACService) RemovePolicy(role, resource, action string) error {
	_, err := s.enforcer.RemovePolicy(role, resource, action)
	if err == nil {
		s.enforcer.SavePolicy()
	}
	return err
}

// initPolicies инициализирует базовые политики
func initPolicies(enforcer *casbin.Enforcer) error {
	// Проверяем, есть ли уже политики (GetPolicy возвращает 1 значение)
	policies := enforcer.GetPolicy()
	if len(policies) > 0 {
		return nil
	}

	// Матрица прав
	type policyRule struct {
		role     string
		resource string
		action   string
	}

	rules := []policyRule{
		// СТУДЕНТ
		{"student", "profile", "read"},
		{"student", "profile", "write"},
		{"student", "challenge", "read"},
		{"student", "challenge", "write"},
		{"student", "challenge", "delete"},
		{"student", "questionnaire", "read"},
		{"student", "questionnaire", "write"},
		{"student", "grade", "read"},
		{"student", "activity", "read"},
		{"student", "activity", "enroll"},
		{"student", "calendar", "read"},
		{"student", "dashboard", "read"},

		// ПРЕПОДАВАТЕЛЬ
		{"teacher", "teacher", "access"},
		{"teacher", "profile", "read"},
		{"teacher", "profile", "write"},
		{"teacher", "student", "read"},
		{"teacher", "grade", "read"},
		{"teacher", "grade", "write"},
		{"teacher", "grade", "delete"},
		{"teacher", "activity", "read"},
		{"teacher", "activity", "write"},
		{"teacher", "activity", "delete"},
		{"teacher", "calendar", "read"},
		{"teacher", "calendar", "write"},
		{"teacher", "attendance", "write"},
		{"teacher", "group", "read"},

		// ДЕРЖАТЕЛЬ
		{"holder", "profile", "read"},
		{"holder", "profile", "write"},
		{"holder", "student", "read"},
		{"holder", "student", "write"},
		{"holder", "challenge", "read"},
		{"holder", "challenge", "comment"},
		{"holder", "grade", "read"},
		{"holder", "activity", "read"},
		{"holder", "activity", "enroll"},
		{"holder", "calendar", "read"},
		{"holder", "calendar", "write"},
		{"holder", "note", "write"},

		// АДМИН
		{"admin", "*", "*"},
	}

	for _, r := range rules {
		if _, err := enforcer.AddPolicy(r.role, r.resource, r.action); err != nil {
			return err
		}
	}

	return enforcer.SavePolicy()
}
