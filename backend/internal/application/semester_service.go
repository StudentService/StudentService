package application

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"backend/internal/domain/semester"
)

type SemesterService struct {
	repo semester.Repository
}

func NewSemesterService(repo semester.Repository) *SemesterService {
	return &SemesterService{repo: repo}
}

func (s *SemesterService) GetAll(ctx context.Context) ([]*semester.Semester, error) {
	return s.repo.GetAll(ctx)
}

func (s *SemesterService) GetActive(ctx context.Context) (*semester.Semester, error) {
	return s.repo.GetActive(ctx)
}

func (s *SemesterService) Create(ctx context.Context, req *semester.CreateSemesterRequest) (*semester.Semester, error) {
	sm := &semester.Semester{
		Name:      req.Name,
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
		IsActive:  req.IsActive,
	}
	if err := s.repo.Create(ctx, sm); err != nil {
		return nil, err
	}
	return sm, nil
}

func (s *SemesterService) Update(ctx context.Context, id uuid.UUID, req *semester.UpdateSemesterRequest) (*semester.Semester, error) {
	sm, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if sm == nil {
		return nil, errors.New("semester not found")
	}

	if req.Name != nil {
		sm.Name = *req.Name
	}
	if req.StartDate != nil {
		sm.StartDate = *req.StartDate
	}
	if req.EndDate != nil {
		sm.EndDate = *req.EndDate
	}
	if req.IsActive != nil {
		sm.IsActive = *req.IsActive
	}

	if err := s.repo.Update(ctx, sm); err != nil {
		return nil, err
	}
	return sm, nil
}

func (s *SemesterService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}
