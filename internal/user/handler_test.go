package user

import (
	"bytes"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gjg-redis-go/internal/user/models"
	"net/http/httptest"
	"testing"
)

func TestNewHandler(t *testing.T) {
	userHandler := NewHandler(nil, nil)

	assert.Implements(t, (*Handler)(nil), userHandler)
}

func TestHandler_CreateUser(t *testing.T) {
	TestUserModel := models.SignIn{
		DisplayName: "ehd",
		Password:    "12345",
		Country:     "tr",
	}

	mockController := gomock.NewController(t)
	defer mockController.Finish()

	t.Run("happy path", func(t *testing.T) {
		app := fiber.New()

		mockUserService := NewMockService(mockController)
		mockUserService.EXPECT().CreateUser(gomock.Any(), &TestUserModel).Return(&models.UserResponseModel{
			UserID:      "1234",
			DisplayName: "ehd",
			Point:       0,
			Rank:        1,
		}, nil)

		userHandler := NewHandler(mockUserService, nil)
		userHandler.RegisterRoutes(app)

		reqBody, err := json.Marshal(&TestUserModel)
		require.NoError(t, err)

		req := httptest.NewRequest(fiber.MethodPost, "/user/create", bytes.NewReader(reqBody))
		req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
	})
}
