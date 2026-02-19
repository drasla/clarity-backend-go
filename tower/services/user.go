package service

import (
	"context"
	"errors"
	"tower/model/maindb"
	"tower/repository"
)

type UserService interface {
	GetUser(ctx context.Context, id uint) (*maindb.User, error)
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) GetUser(ctx context.Context, id uint) (*maindb.User, error) {
	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("user not found")
	}

	if user.Status == maindb.StatusWithdrawn {
		return nil, errors.New("user has withdrawn")
	}
	if user.Status == maindb.StatusSuspended {
		return nil, errors.New("user is suspended")
	}

	return user, nil
}
