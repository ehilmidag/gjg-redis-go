package main

import (
	"fmt"
	"github.com/go-redis/redis/v8"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"gjg-redis-go/internal/user"
	"gjg-redis-go/internal/user/models"
	_ "gjg-redis-go/pkg/cerror"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

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
	db, err := initializeDB()
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

func RedisClient(redisConfig *config.RedisConfig) *redis.Client {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     redisConfig.Addr,
		Password: redisConfig.Password,
		DB:       redisConfig.Db,
	})
	return redisClient
}

func initializeDB() (*gorm.DB, error) {
	dsn := "root:12345@tcp(127.0.0.1:3306)/user"
	fmt.Println(dsn, "birinci dsn")
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		fmt.Println("Veritabanına bağlanırken bir hata oluştu:", err)
		return nil, err
	}

	db = db.Debug()

	// Veritabanını otomatik olarak oluştur
	err = db.Exec("CREATE DATABASE IF NOT EXISTS user").Error
	if err != nil {
		fmt.Println("Veritabanı oluşturulurken bir hata oluştu:", err)
		return nil, err
	}

	// Veritabanı bağlantısını kapat ve yeni veritabanı adıyla tekrar aç
	dbSQL, err := db.DB()
	if err != nil {
		fmt.Println("Veritabanı bağlantısı kapatılırken bir hata oluştu:", err)
		return nil, err

	}
	dbSQL.Close()

	dsn = "root:12345@tcp(127.0.0.1:3306)/user"
	fmt.Println(dsn, "ikinci dsn")

	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		fmt.Println("Yeni veritabanına bağlanırken bir hata oluştu:", err)
		return nil, err

	}

	// Gerekli tabloyu otomatik olarak oluştur
	err = db.AutoMigrate(&models.UserCreateEntity{})
	if err != nil {
		fmt.Println("Tablo oluşturulurken bir hata oluştu:", err)
		return nil, err
	}

	fmt.Println("Veritabanına başarıyla bağlandı ve gerekli tablo oluşturuldu.")

	return db, nil
}
