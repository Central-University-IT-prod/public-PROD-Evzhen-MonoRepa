package database

import (
	"github.com/Central-University-IT-prod/PROD-Evzhen-MonoRepa/internal/entities"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

var DB *gorm.DB

func Connect() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error load .env file")
	}
	dsn := os.Getenv("DSN")

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second, // Slow SQL threshold
			LogLevel:      logger.Info, // Log level
			Colorful:      true,        // Disable color
		},
	)

	log.Println(dsn)
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		panic("Не получилось подключиться к бд")
	} else {
		log.Println("Успешно подключились к базе данных")
	}
	DB = database
	err = database.AutoMigrate(
		&entities.Admin{},
		&entities.User{},
		&entities.Contest{},
		&entities.Profile{},
		&entities.Command{},
		&entities.Invite{})
}
