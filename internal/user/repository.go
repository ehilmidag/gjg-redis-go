package user

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/go-redis/redis/v8"
	"gjg-redis-go/internal/user/models"
	"log"
	"time"
)

type DataRepository interface {
	RedisCreateUser(ctx context.Context, userid string) error
	RedisCreateUserByCountry(ctx context.Context, userid string, countryCode string) error
	RedisGetRankByID(ctx context.Context, userid string) (int64, error)
	RedisScoreUpdate(ctx context.Context, userid string, score float64) error
	RedisScoreUpdateByCountry(ctx context.Context, userid string, countryCode string, score float64) error
	MySQLCreateUser(ctx context.Context, user *models.UserCreateEntity) (*models.UserCreateDTO, error)
	MySQLGetUserByID(ctx context.Context, userid string) (*models.UserCreateDTO, error)
	MySQLUpdatePoint(ctx context.Context, update *models.SendScoreEntity) (*models.SendScoreDto, error)
}

type repository struct {
	mysql *sql.DB
	redis *redis.Client
}

func NewRepository(sqlDB *sql.DB, redisClient *redis.Client) *repository {
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

	_, err := r.redis.ZAdd(ctx, fmt.Sprint("leaderboard:"+countryCode), &redis.Z{
		Score:  0,
		Member: userid,
	}).Result()
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

func (r *repository) MySQLCreateUser(ctx context.Context, user *models.UserCreateEntity) (*models.UserCreateDTO, error) {

	query := `
		INSERT INTO user (user_id, display_name, hashed_password, points, user_rank,country_code,created_at,updated_at)
		VALUES (?, ?, ?, ?, ?, ?,?,?)
	`
	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	stmt, err := r.mysql.PrepareContext(ctx, query)
	if err != nil {
		log.Printf("Error %s when preparing SQL statement", err)
		return nil, err
	}
	defer stmt.Close()
	res, err := stmt.ExecContext(ctx, user.UserID, user.DisplayName, user.HashedPassword, user.Points, user.Rank, user.CountryCode, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		log.Printf("Error %s when inserting row into user table", err)
		return nil, err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		log.Printf("Error %s when finding rows affected", err)
		return nil, err
	}
	log.Printf("%d user created ", rows)
	userDto := &models.UserCreateDTO{
		UserID:      user.UserID,
		DisplayName: user.DisplayName,
		Points:      user.Points,
		Rank:        user.Rank,
	}
	return userDto, nil
}
func (r *repository) MySQLGetUserByID(ctx context.Context, userid string) (*models.UserCreateDTO, error) {
	query := "SELECT * FROM user WHERE user_id = ?"
	row := r.mysql.QueryRow(query, userid)

	var userEntity models.UserCreateEntity

	err := row.Scan(&userEntity.UserID, &userEntity.DisplayName, &userEntity.HashedPassword, &userEntity.Points, &userEntity.Rank, &userEntity.CountryCode, &userEntity.CreatedAt, &userEntity.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("Kullanıcı bulunamadı.")
		}
		return nil, err
	}

	userdto := &models.UserCreateDTO{
		UserID:      userEntity.UserID,
		DisplayName: userEntity.DisplayName,
		Points:      userEntity.Points,
		Rank:        userEntity.Rank,
	}
	fmt.Println(userdto)

	return userdto, nil
}

func (r *repository) MySQLUpdatePoint(ctx context.Context, update *models.SendScoreEntity) (*models.SendScoreDto, error) {

	stmt, err := r.mysql.Prepare("UPDATE user SET points = ?, user_rank = ?, updated_at = ? WHERE user_id = ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	_, err = stmt.Exec(&update.TotalScore, &update.Rank, &update.TimeStamp, &update.UserID)
	if err != nil {
		return nil, err
	}
	sendScoreDto := &models.SendScoreDto{
		UserID:     update.UserID,
		ScoreWorth: update.TotalScore,
		TimeStamp:  update.TimeStamp,
	}
	return sendScoreDto, nil
}
