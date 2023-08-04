package boostrap

import (
	"fmt"
	"os"

	"github.com/Jhon-2801/courses-domain/domain"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func DBConnection() (*gorm.DB, error) {
	_ = godotenv.Load()

	user := os.Getenv("DATABASE_USER")
	password := os.Getenv("DATABASE_PASSWORD")
	host := os.Getenv("DATABASE_HOST")
	port := os.Getenv("DATABASE_PORT")
	name := os.Getenv("DATABASE_NAME")

	dsn := fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=utf8&parseTime=True",
		user,
		password,
		host,
		port,
		name)

	fmt.Println(dsn)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	if os.Getenv("DATABASE_MIGRATE") == "true" {
		if err := db.AutoMigrate(&domain.User{}); err != nil {
			return nil, err
		}
	}
	return db, nil
}
