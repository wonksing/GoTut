package user

import (
	"context"

	"github.com/wonksing/gotut/test/mock/go-mock/models"
)

type Repository interface {
	Find(c context.Context, name string, age int16) (*models.User, error)
}
