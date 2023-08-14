package user

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gjg-redis-go/internal/user/models"
	"testing"
)

const (
	TestUserId        = "abcd-abcd-abcd-abcd-abcd"
	TestCryptPassword = "$2a$07$21Py6b8E1XWLlpSS1ASxK.RhNpvm1n3q34G9uqysCwx/ciP0vSaEm"
)

func TestNewService(t *testing.T) {
	userService := NewService(nil)

	assert.Implements(t, (*Service)(nil), userService)
}

func TestService_CreateUser(t *testing.T) {
	mockController := gomock.NewController(t)
	defer mockController.Finish()
	TestUserModel := models.SignIn{
		DisplayName: "ehd",
		Password:    "12345",
		Country:     "tr",
	}

	t.Run("happy path", func(t *testing.T) {
		ctx := context.Background()
		mockUserRepository := NewMockRepository(mockController)
		mockUserRepository.
			EXPECT().RedisCreateUser(ctx, gomock.Any()).
			Return(nil)

		mockUserRepository.
			EXPECT().
			RedisCreateUserByCountry(ctx, gomock.Any(), "tr").
			Return(nil)
		mockUserRepository.EXPECT().RedisGetRankByID(ctx, gomock.Any()).Return(int64(1), nil)
		mockUserRepository.EXPECT().MySQLCreateUser(ctx, gomock.Any()).Return(nil)

		userService := NewService(mockUserRepository)
		response, err := userService.CreateUser(ctx, &TestUserModel)

		assert.NoError(t, err)
		assert.NotEmpty(t, response)

	})
}
