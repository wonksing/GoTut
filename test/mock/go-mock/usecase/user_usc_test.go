package usecase

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/wonksing/gotut/test/mock/go-mock/models"
	"github.com/wonksing/gotut/test/mock/go-mock/port/user/mocks"

	"github.com/stretchr/testify/assert"
)

func TestFind(t *testing.T) {
	user := &models.User{
		Name:   "won",
		Age:    99,
		Gender: "M",
	}

	c := context.Background()
	cto := 6 * time.Second

	ctrl := gomock.NewController(t)

	// Assert that Bar() is invoked.
	defer ctrl.Finish()

	m := mocks.NewMockUserRepository(ctrl)

	// Asserts that the first and only call to Bar() is passed 99.
	// Anything else will fail.
	m.
		EXPECT().
		Find(gomock.Any(), "won", int16(99)).
		Return(user, nil)

	// m := new(mocks.MockUserRepository)
	// m.On("Find", mock.Anything, "won", int16(99)).Return(user, nil)

	u := NewUserUsecase(m, cto)
	res, err := u.Find(c, "won", int16(99))
	assert.Nil(t, err)

	fmt.Println(res)
	// m.AssertExpectations(t)
}
