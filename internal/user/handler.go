package user

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"gjg-redis-go/internal/user/models"
	"gjg-redis-go/pkg/cerror"
	"gjg-redis-go/pkg/server"
)

type Handler interface {
	server.Handler
	CreateUser(ctx *fiber.Ctx) error
	GetUserByID(ctx *fiber.Ctx) error
}

type handler struct {
	userService    Service
	userRepository Repository
}

func (h *handler) RegisterRoutes(app *fiber.App) {
	app.Post("/user/create", h.CreateUser)
	app.Get("/user/profile/:userID", h.GetUserByID)
	app.Post("/score/submit", h.ScoreSubmit)
	app.Get("/leaderboard", h.LeaderBoard)
	app.Get("/leaderboard/:country_iso_code", h.LeaderBoardByCountry)

}

func NewHandler(userService Service, userRepository Repository) Handler {
	return &handler{
		userService:    userService,
		userRepository: userRepository,
	}
}

func (h *handler) CreateUser(ctx *fiber.Ctx) error {
	var err error

	var userCreateRequest *models.SignIn
	err = ctx.BodyParser(&userCreateRequest)
	if err != nil {
		cerr := cerror.NewError(400, "Bad Request babuş")
		return cerr
	}
	userCreateResponse, err := h.userService.CreateUser(ctx.Context(), userCreateRequest)
	if err != nil {
		cerr := cerror.NewError(500, err.Error())
		return cerr
	}
	return ctx.Status(201).JSON(userCreateResponse)
}

func (h *handler) GetUserByID(ctx *fiber.Ctx) error {

	userId := ctx.Params("userID")
	fmt.Println(userId)

	userGetByIDResponse, err := h.userService.GetUserDetailsByID(ctx.Context(), userId)
	if err != nil {
		cerr := cerror.NewError(500, err.Error())
		return cerr
	}
	return ctx.Status(200).JSON(&userGetByIDResponse)
}

func (h *handler) ScoreSubmit(ctx *fiber.Ctx) error {
	var scoreUpdateRequest *models.SendScore
	err := ctx.BodyParser(&scoreUpdateRequest)

	if err != nil {
		cerr := cerror.NewError(400, "Parse error")
		return cerr
	}
	scoreUpdateResponse, err := h.userService.UpdateScore(ctx.Context(), scoreUpdateRequest)
	if err != nil {
		return err
	}
	return ctx.Status(200).JSON(&scoreUpdateResponse)
}

func (h *handler) LeaderBoard(ctx *fiber.Ctx) error {
	leaderBoard, err := h.userService.GetLeaderBoard(ctx.Context())
	if err != nil {
		cerr := cerror.NewError(500, err.Error())
		return cerr
	}
	return ctx.Status(200).JSON(&leaderBoard)
}
func (h *handler) LeaderBoardByCountry(ctx *fiber.Ctx) error {
	country := ctx.Params("country_iso_code")

	leaderBoard, err := h.userService.GetLeaderBoardByCountry(ctx.Context(), country)
	if err != nil {
		cerr := cerror.NewError(500, err.Error())
		return cerr
	}
	return ctx.Status(200).JSON(&leaderBoard)
}
