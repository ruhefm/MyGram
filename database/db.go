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
	// 	db.Exec(`
	// CREATE TABLE users (
	//     id SERIAL PRIMARY KEY,
	//     username VARCHAR(50) UNIQUE NOT NULL,
	//     email VARCHAR(150) UNIQUE NOT NULL,
	//     password TEXT NOT NULL,
	//     age INT NOT NULL CHECK (age >= 9),
	//     profile_image_url TEXT,
	//     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	//     updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	// );

	// INSERT INTO users (username, email, password, age, profile_image_url, created_at, updated_at)
	// VALUES ('john_doe', 'john@example.com', 'password123', 25, 'http://example.com/profile.jpg', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
	//        ('jane_doe', 'jane@example.com', 'password456', 30, 'http://example.com/profile.jpg', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

	// CREATE TABLE social_medias (
	//     id SERIAL PRIMARY KEY,
	//     name VARCHAR(50) NOT NULL,
	//     social_media_url VARCHAR(50) NOT NULL,
	//     user_id INT NOT NULL,
	//     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	//     updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	//     CONSTRAINT fk_users_social_medias FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	// );

	// INSERT INTO social_medias (name, social_media_url, user_id, created_at, updated_at)
	// VALUES ('Facebook', 'http://facebook.com/johndoe', 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
	//        ('Twitter', 'http://twitter.com/johndoe', 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

	// CREATE TABLE photos (
	//     id SERIAL PRIMARY KEY,
	//     title VARCHAR(100) NOT NULL,
	//     caption VARCHAR(200),
	//     photo_url TEXT NOT NULL,
	//     user_id INT NOT NULL,
	//     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	//     updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	//     CONSTRAINT fk_users_photos FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	// );

	// INSERT INTO photos (title, caption, photo_url, user_id, created_at, updated_at)
	// VALUES ('Photo 1', 'Caption 1', 'http://example.com/photo1.jpg', 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
	//        ('Photo 2', 'Caption 2', 'http://example.com/photo2.jpg', 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

	// CREATE TABLE comments (
	//     id SERIAL PRIMARY KEY,
	//     user_id INT NOT NULL,
	//     photo_id INT NOT NULL,
	//     message VARCHAR(200),
	//     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	//     updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	//     CONSTRAINT fk_users_comments FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
	//     CONSTRAINT fk_photos_comments FOREIGN KEY (photo_id) REFERENCES photos(id) ON DELETE CASCADE
	// );

	// INSERT INTO comments (user_id, photo_id, message, created_at, updated_at)
	// VALUES (1, 1, 'Comment 1', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
	//        (1, 2, 'Comment 2', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);
	// 	`)
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
