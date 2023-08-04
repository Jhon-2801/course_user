package user

import (
	"context"

	"github.com/Jhon-2801/courses-domain/domain"
)

type (
	Filters struct {
		FirstName string
		LastName  string
	}
	Service interface {
		Create(ctx context.Context, firsName, LastName, email, phone string) (*domain.User, error)
		Get(ctx context.Context, id string) (*domain.User, error)
		GetAll(ctx context.Context, filters Filters, offset, limit int) ([]domain.User, error)
		Delete(ctx context.Context, id string) error
		Update(ctx context.Context, id string, firstName *string, lastName *string, email *string, phone *string) error
		Count(ctx context.Context, filters Filters) (int, error)
	}
	service struct {
		repo Repository
	}
)

func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}

func (s service) Create(ctx context.Context, firstName, LastName, email, phone string) (*domain.User, error) {
	user := domain.User{
		FirstName: firstName,
		LastName:  LastName,
		Email:     email,
		Phone:     phone,
	}
	if err := s.repo.Create(ctx, &user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (s service) GetAll(ctx context.Context, filters Filters, offsite, limit int) ([]domain.User, error) {
	users, err := s.repo.GetAll(ctx, filters, offsite, limit)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (s service) Get(ctx context.Context, id string) (*domain.User, error) {
	user, err := s.repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}
func (s service) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s service) Update(ctx context.Context, id string, firstName *string, lastName *string, email *string, phone *string) error {
	return s.repo.Update(ctx, id, firstName, lastName, email, phone)
}
func (s service) Count(ctx context.Context, filters Filters) (int, error) {
	return s.repo.Count(ctx, filters)
}
