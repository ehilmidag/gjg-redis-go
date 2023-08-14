package user

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"gjg-redis-go/internal/user/models"
	"gorm.io/gorm"
)

type Repository interface {
	RedisCreateUser(ctx context.Context, userid string) error
	RedisCreateUserByCountry(ctx context.Context, userid string, countryCode string) error
	RedisGetRankByID(ctx context.Context, userid string) (int64, error)
	RedisScoreUpdate(ctx context.Context, userid string, score float64) error
	RedisScoreUpdateByCountry(ctx context.Context, userid string, countryCode string, score float64) error
	RedisLeaderBoard(ctx context.Context) ([]string, error)
	RedisLeaderBoardByCountry(ctx context.Context, country string) ([]string, error)
	MySQLCreateUser(ctx context.Context, user *models.UserCreateEntity) error
	MySQLGetUserByID(ctx context.Context, userid string) (*models.UserCreateEntity, error)
	MySQLUpdatePoint(ctx context.Context, update *models.SendScoreEntity) (*models.SendScoreDto, error)
}

type repository struct {
	mysql *gorm.DB
	redis *redis.Client
}

func NewRepository(sqlDB *gorm.DB, redisClient *redis.Client) *repository {
	return &repository{
		mysql: sqlDB,
		redis: redisClient,
	}
}

func (r *repository) RedisCreateUser(ctx context.Context, userid string) error {

	_, err := r.redis.ZAdd(ctx, "leaderboard", &redis.Z{
		Score:  0,
		Member: userid,
	}).Result()
	if err != nil {
		return err
	}
	return nil
}
func (r *repository) RedisCreateUserByCountry(ctx context.Context, userid string, countryCode string) error {
	fmt.Println(countryCode)
	fmt.Println(fmt.Sprintf("leaderboard:%s", countryCode))
	a, err := r.redis.ZAdd(ctx, fmt.Sprint("leaderboard:", countryCode), &redis.Z{
		Score:  0,
		Member: userid,
	}).Result()
	fmt.Println(a, "redis createden dönen")
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) RedisGetRankByID(ctx context.Context, userid string) (int64, error) {
	rank, err := r.redis.ZRevRank(ctx, "leaderboard", userid).Result()
	if err != nil {
		return -1, err
	}
	return rank, nil
}

func (r *repository) RedisScoreUpdate(ctx context.Context, userid string, score float64) error {
	_, err := r.redis.ZIncrBy(ctx, "leaderboard", score, userid).Result()
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) RedisScoreUpdateByCountry(ctx context.Context, userid string, countryCode string, score float64) error {
	_, err := r.redis.ZIncrBy(ctx, fmt.Sprint("leaderboard:"+countryCode), score, userid).Result()
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) MySQLCreateUser(ctx context.Context, user *models.UserCreateEntity) error {

	result := r.mysql.Create(user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *repository) MySQLGetUserByID(ctx context.Context, userid string) (*models.UserCreateEntity, error) {
	var userEntity models.UserCreateEntity

	result := r.mysql.Where("user_id = ?", userid).First(&userEntity)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("Kullanıcı bulunamadı.")
		}
		return nil, result.Error
	}

	return &userEntity, nil
}

func (r *repository) MySQLUpdatePoint(ctx context.Context, update *models.SendScoreEntity) (*models.SendScoreDto, error) {
	var sendScoreDto models.SendScoreDto

	result := r.mysql.Model(&models.SendScoreDto{}).Where("user_id = ?", update.UserID).
		Updates(map[string]interface{}{
			"total_score": update.TotalScore,
			"updated_at":  update.TimeStamp,
		})
	if result.Error != nil {
		return nil, result.Error
	}

	sendScoreDto.UserID = update.UserID
	sendScoreDto.ScoreWorth = update.TotalScore
	sendScoreDto.TimeStamp = update.TimeStamp
	return &sendScoreDto, nil
}

func (r *repository) RedisLeaderBoard(ctx context.Context) ([]string, error) {
	results, err := r.redis.ZRevRange(ctx, "leaderboard", 0, 999).Result()
	if err != nil {
		return nil, err
	}
	return results, nil
}

func (r *repository) RedisLeaderBoardByCountry(ctx context.Context, country string) ([]string, error) {
	results, err := r.redis.ZRevRange(ctx, fmt.Sprintf("leaderboard:%s", country), 0, 999).Result()
	if err != nil {
		return nil, err
	}
	return results, nil
}
