package database

import (
	"fmt"
	"mygram/models"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	host     = os.Getenv("PGHOST")
	username = os.Getenv("PGUSER")
	password = os.Getenv("PGPASSWORD")
	dbName   = os.Getenv("PGDATABASE")
	port     = os.Getenv("PGPORT")
	db       *gorm.DB
	err      error
)

func StartDB() {
	config := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%v sslmode=disable TimeZone=Asia/Jakarta", host, username, password, dbName, port)
	db, err = gorm.Open(postgres.Open(config), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	if err != nil {
		panic(err)
	}
	err = db.AutoMigrate(&models.Users{}, &models.Social_Medias{}, &models.Photos{}, &models.Comments{})
	if err != nil {
		panic(err)
	}
	fmt.Println("DB LOG: DB Connected")
}

func GetDB() *gorm.DB {
	return db
}

func UserRegister(register *models.Users) error {
	db := GetDB()

	if err := db.Create(register).Error; err != nil {
		return err
	}

	return nil
}
