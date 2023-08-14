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
	CreateUser(ctx context.Context, user *models.SignIn) (*models.UserResponseModel, error)
	UpdateScore(ctx context.Context, scoreUpdate *models.SendScore) (*models.SendScoreDto, error)
	GetUserDetailsByID(ctx context.Context, userid string) (*models.UserResponseModel, error)
	GetLeaderBoard(ctx context.Context) (*[]models.LeaderBoard, error)
	GetLeaderBoardByCountry(ctx context.Context, country string) (*[]models.LeaderBoard, error)
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

func (s *service) CreateUser(ctx context.Context, user *models.SignIn) (*models.UserResponseModel, error) {
	userid := uuid.New().String()
	err := s.userRepository.RedisCreateUser(ctx, userid)
	if err != nil {
		cerr := cerror.NewError(500, "create user redis error")
		return nil, cerr
	}
	fmt.Println("redis user created")
	err = s.userRepository.RedisCreateUserByCountry(ctx, userid, user.Country)
	if err != nil {
		cerr := cerror.NewError(500, "create user redis error")
		return nil, cerr
	}
	fmt.Println("redis user country base createdx")
	rank, err := s.userRepository.RedisGetRankByID(ctx, userid)
	if err != nil {
		return nil, err
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		cerr := cerror.NewError(500, "Password Hashing Error")
		return nil, cerr
	}
	err = s.userRepository.MySQLCreateUser(ctx, &models.UserCreateEntity{
		UserID:         userid,
		DisplayName:    user.DisplayName,
		HashedPassword: string(hashedPassword),
		Country:        user.Country,
		Points:         0,
		CreatedAt:      time.Now().Unix(),
		UpdatedAt:      time.Now().Unix(),
	})
	if err != nil {
		cerr := cerror.NewError(500, "Repository Failed")
		return nil, cerr
	}
	userResponse := &models.UserResponseModel{
		UserID:      userid,
		DisplayName: user.DisplayName,
		Point:       0,
		Rank:        rank,
	}

	return userResponse, nil
}

func (s *service) UpdateScore(ctx context.Context, scoreUpdate *models.SendScore) (*models.SendScoreDto, error) {
	err := s.userRepository.RedisScoreUpdate(ctx, scoreUpdate.UserID, scoreUpdate.ScoreWorth)
	if err != nil {
		cerr := cerror.NewError(500, "Repository Failed")
		return nil, cerr
	}
	err = s.userRepository.RedisScoreUpdateByCountry(ctx, scoreUpdate.UserID, scoreUpdate.Country, scoreUpdate.ScoreWorth)
	if err != nil {
		cerr := cerror.NewError(500, "Repository Failed")
		return nil, cerr
	}
	current, err := s.userRepository.MySQLGetUserByID(ctx, scoreUpdate.UserID)
	if err != nil {
		return nil, err
	}

	sendscoreEntity := &models.SendScoreEntity{
		UserID:     scoreUpdate.UserID,
		TotalScore: scoreUpdate.ScoreWorth + current.Points,
		TimeStamp:  time.Now().Unix(),
	}
	SendScoreResponse, err := s.userRepository.MySQLUpdatePoint(ctx, sendscoreEntity)
	if err != nil {
		return nil, err
	}
	SendScoreResponse.ScoreWorth = scoreUpdate.ScoreWorth
	return SendScoreResponse, nil
}

func (s *service) GetUserDetailsByID(ctx context.Context, userid string) (*models.UserResponseModel, error) {
	user, err := s.userRepository.MySQLGetUserByID(ctx, userid)
	if err != nil {
		return nil, err
	}
	rank, err := s.userRepository.RedisGetRankByID(ctx, userid)
	if err != nil {
		return nil, err
	}
	userResponse := &models.UserResponseModel{
		UserID:      user.UserID,
		DisplayName: user.DisplayName,
		Point:       user.Points,
		Rank:        rank,
	}
	return userResponse, nil
}

func (s *service) GetLeaderBoard(ctx context.Context) (*[]models.LeaderBoard, error) {
	usersFromRedis, err := s.userRepository.RedisLeaderBoard(ctx)
	if err != nil {
		return nil, err
	}
	leaderBoard := []models.LeaderBoard{}
	for i, userid := range usersFromRedis {
		user, err := s.userRepository.MySQLGetUserByID(ctx, userid)
		if err != nil {
			return nil, err
		}
		leaderBoard = append(leaderBoard, models.LeaderBoard{
			DisplayName: user.DisplayName,
			Point:       user.Points,
			Rank:        int64(i + 1),
			Country:     user.Country,
		})
	}
	return &leaderBoard, nil
}

func (s *service) GetLeaderBoardByCountry(ctx context.Context, country string) (*[]models.LeaderBoard, error) {
	usersFromRedis, err := s.userRepository.RedisLeaderBoardByCountry(ctx, country)
	if err != nil {
		return nil, err
	}
	leaderBoardByCountry := []models.LeaderBoard{}
	for _, userid := range usersFromRedis {
		rank, err := s.userRepository.RedisGetRankByID(ctx, userid)
		if err != nil {
			return nil, err
		}
		user, err := s.userRepository.MySQLGetUserByID(ctx, userid)
		if err != nil {
			return nil, err
		}
		leaderBoardByCountry = append(leaderBoardByCountry, models.LeaderBoard{
			DisplayName: user.DisplayName,
			Point:       user.Points,
			Rank:        rank + 1,
			Country:     user.Country,
		})
	}
	return &leaderBoardByCountry, nil
}
