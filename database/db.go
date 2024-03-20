package database

import (
	"fmt"
	"mygram/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	host     = "localhost"
	username = "postgres"
	password = "test123456"
	dbName   = "db-go-sql"
	port     = 5432
)

var (
	db  *gorm.DB
	err error
)

func StartDB() {
	config := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable", host, username, password, dbName, port)
	db, err = gorm.Open(postgres.Open(config), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&models.Users{}, &models.Social_Medias{}, &models.Photos{}, &models.Comments{})
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

// func GetUserByEmail(items string) (*models.Users, error) {
// 	var request models.Users
// 	db := GetDB()

// 	if err := db.First(&request, "email = ?", items).Error; err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			return nil, errors.New("email tidak ditemukan")
// 		}
// 		return nil, err
// 	}
// 	return &request, nil
// }

// func GetUserByUsername(items string) (*models.Users, error) {
// 	var request models.Users
// 	db := GetDB()

// 	if err := db.First(&request, "username = ?", items).Error; err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			return nil, errors.New("user tidak ditemukan")
// 		}
// 		return nil, err
// 	}
// 	return &request, nil
// }
