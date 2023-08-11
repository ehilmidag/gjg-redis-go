package main

import (
	"database/sql"
	"fmt"
	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"gjg-redis-go/internal/user"
	"log"
	"os"
	"path/filepath"

	"gjg-redis-go/pkg/config"
	"gjg-redis-go/pkg/logger"
	"gjg-redis-go/pkg/project"
	"gjg-redis-go/pkg/server"
)

func main() {
	log := logger.NewLogger()

	isAtRemote := os.Getenv(config.IsAtRemote)
	if isAtRemote == "" {
		rootDirectory := project.GetRootDirectory()
		fmt.Println(rootDirectory)
		dotenvPath := filepath.Join(rootDirectory, "/.env")
		err := godotenv.Load(dotenvPath)
		if err != nil {
			log.Fatal(err)
		}
	}

	cfg, err := config.ReadConfig()
	if err != nil {
		log.Fatal(err)
	}
	mySqlconf, err := config.ReadMySqlConfig()
	fmt.Println(mySqlconf)
	db, err := initializeDB(mySqlconf)
	if err != nil {
		log.Fatal(err)
	}
	redisconf, err := config.ReadRedisConfig()
	redisClient := RedisClient(redisconf)
	repository := user.NewRepository(db, redisClient)
	service := user.NewService(repository)
	handler := user.NewHandler(service, repository)
	var handlers []server.Handler
	handlers = append(handlers, handler)
	srv := server.NewServer(cfg, handlers)
	app := srv.GetFiberInstance()
	app.Use(logger.Middleware(log))
	app.Get("/", func(c *fiber.Ctx) error {
		// Render a template named 'index.html' with content
		return c.Render("index", fiber.Map{
			"Title":       "Hello, World!",
			"Description": "This is a template.",
		})
	})
	err = srv.Start()
	if err != nil {
		log.Fatal(err)
	}
}

func dsn(sqlConfig *config.MySQLConfig) string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s", sqlConfig.Username, sqlConfig.Password, sqlConfig.Hostname, "user")

}

func RedisClient(redisConfig *config.RedisConfig) *redis.Client {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     redisConfig.Addr,
		Password: redisConfig.Password,
		DB:       redisConfig.Db,
	})
	return redisClient
}

func initializeDB(sqlConfig *config.MySQLConfig) (*sql.DB, error) {
	dataSourceName := dsn(sqlConfig)
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		log.Fatal("Veritabanı bağlantısı oluşturulamadı:", err)
	}

	fmt.Println(sqlConfig.Dbname)

	// "user" tablosunu oluşturma
	createTableQuery := `
	CREATE TABLE IF NOT EXISTS user (
   user_id VARCHAR(255) PRIMARY KEY,
    display_name VARCHAR(255),
    hashed_password VARCHAR(255),
    points DOUBLE,
    user_rank BIGINT,
    country_code VARCHAR(255),
    created_at BIGINT,
    updated_at BIGINT
);
	`

	_, err = db.Exec(createTableQuery)
	if err != nil {
		log.Fatal("Tablo oluşturulurken hata oluştu:", err)
	}

	fmt.Println("user tablosu başarıyla oluşturuldu.")
	return db, err
}
