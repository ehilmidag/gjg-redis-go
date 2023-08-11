package user

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"gjg-redis-go/internal/user/models"
	"gjg-redis-go/pkg/cerror"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type Service interface {
	CreateUser(ctx context.Context, user *models.SignIn) (*models.UserCreateDTO, error)
	UpdateScore(ctx context.Context, scoreUpdate *models.SendScore) (*models.SendScoreDto, error)
}

type service struct {
	userRepository *repository
}

func NewService(
	userRepository *repository,
) Service {
	return &service{
		userRepository: userRepository,
	}
}

func (s *service) CreateUser(ctx context.Context, user *models.SignIn) (*models.UserCreateDTO, error) {
	userid := uuid.New().String()
	err := s.userRepository.RedisCreateUser(ctx, userid)
	if err != nil {
		cerr := cerror.NewError(500, "create user redis error")
		return nil, cerr
	}
	fmt.Println("redis user created")
	err = s.userRepository.RedisCreateUserByCountry(ctx, userid, user.CountryCode)
	if err != nil {
		cerr := cerror.NewError(500, "create user redis error")
		return nil, cerr
	}
	fmt.Println("redis user country base createdx")

	rank, err := s.userRepository.RedisGetRankByID(ctx, userid)
	if err != nil {
		cerr := cerror.NewError(500, "Error while getting redis rank")
		return nil, cerr
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		cerr := cerror.NewError(500, "Password Hashing Error")
		return nil, cerr
	}
	userCreateResponse, err := s.userRepository.MySQLCreateUser(ctx, &models.UserCreateEntity{
		UserID:         userid,
		DisplayName:    user.DisplayName,
		HashedPassword: string(hashedPassword),
		Points:         0,
		Rank:           rank + 1,
		CountryCode:    user.CountryCode,
		CreatedAt:      time.Now().Unix(),
		UpdatedAt:      time.Now().Unix(),
	})
	if err != nil {
		cerr := cerror.NewError(500, "Repository Failed")
		return nil, cerr
	}
	return userCreateResponse, nil
}

func (s *service) UpdateScore(ctx context.Context, scoreUpdate *models.SendScore) (*models.SendScoreDto, error) {
	err := s.userRepository.RedisScoreUpdate(ctx, scoreUpdate.UserID, scoreUpdate.ScoreWorth)
	if err != nil {
		cerr := cerror.NewError(500, "Repository Failed")
		return nil, cerr
	}
	err = s.userRepository.RedisScoreUpdateByCountry(ctx, scoreUpdate.UserID, scoreUpdate.CountryCode, scoreUpdate.ScoreWorth)
	if err != nil {
		cerr := cerror.NewError(500, "Repository Failed")
		return nil, cerr
	}

	rank, err := s.userRepository.RedisGetRankByID(ctx, scoreUpdate.UserID)
	if err != nil {
		cerr := cerror.NewError(500, "Error while getting redis rank")
		return nil, cerr
	}
	fmt.Println(rank, "AGA RANK BURDA")
	current, err := s.userRepository.MySQLGetUserByID(ctx, scoreUpdate.UserID)
	if err != nil {
		return nil, err
	}

	sendscoreEntity := &models.SendScoreEntity{
		UserID:     scoreUpdate.UserID,
		TotalScore: scoreUpdate.ScoreWorth + current.Points,
		Rank:       rank + 1,
		TimeStamp:  time.Now().Unix(),
	}
	SendScoreResponse, err := s.userRepository.MySQLUpdatePoint(ctx, sendscoreEntity)
	if err != nil {
		return nil, err
	}
	return SendScoreResponse, nil
}
