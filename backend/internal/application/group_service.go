package application

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"backend/internal/domain/group"
)

type GroupService struct {
	repo group.Repository
}

func NewGroupService(repo group.Repository) *GroupService {
	return &GroupService{repo: repo}
}

func (s *GroupService) GetAll(ctx context.Context) ([]*group.Group, error) {
	return s.repo.GetAll(ctx)
}

func (s *GroupService) GetByID(ctx context.Context, id uuid.UUID) (*group.Group, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *GroupService) Create(ctx context.Context, req *group.CreateGroupRequest) (*group.Group, error) {
	g := &group.Group{
		Name:       req.Name,
		CourseID:   req.CourseID,
		SemesterID: req.SemesterID,
		HolderID:   req.HolderID,
	}
	if err := s.repo.Create(ctx, g); err != nil {
		return nil, err
	}
	return g, nil
}

func (s *GroupService) Update(ctx context.Context, id uuid.UUID, req *group.UpdateGroupRequest) (*group.Group, error) {
	g, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if g == nil {
		return nil, errors.New("group not found")
	}

	if req.Name != nil {
		g.Name = *req.Name
	}
	if req.CourseID != nil {
		g.CourseID = *req.CourseID
	}
	if req.SemesterID != nil {
		g.SemesterID = *req.SemesterID
	}
	if req.HolderID != nil {
		g.HolderID = req.HolderID
	}

	if err := s.repo.Update(ctx, g); err != nil {
		return nil, err
	}
	return g, nil
}

func (s *GroupService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}
