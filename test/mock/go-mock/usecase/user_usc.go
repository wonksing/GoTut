package usecase

import (
	"context"
	"time"

	"github.com/wonksing/gotut/test/mock/go-mock/port/user"

	"github.com/wonksing/gotut/test/mock/go-mock/models"
)

type Usecase interface {
	Find(c context.Context, name string, age int16) (*models.User, error)
}

type userUsecase struct {
	repo user.Repository
	cto  time.Duration
}

func NewUserUsecase(repo user.Repository, cto time.Duration) Usecase {
	return &userUsecase{
		repo: repo,
		cto:  cto,
	}
}

func (u *userUsecase) Find(c context.Context, name string, age int16) (*models.User, error) {
	c, cancel := context.WithTimeout(c, u.cto)
	defer cancel()

	user, err := u.repo.Find(c, name, age)
	if err != nil {
		return nil, err
	}

	return user, nil
}
