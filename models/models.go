package models

import (
	"errors"
	"mygram/helpers"
	"strings"
	"time"

	"github.com/asaskevich/govalidator"
	"gorm.io/gorm"
)

type Social_Medias struct {
	ID             uint      `json:"id" gorm:"primary_key;type:bigint"`
	Name           string    `json:"name" gorm:"type:varchar(50);not null"`
	SocialMediaURL string    `json:"social_media_url" gorm:"type:varchar(50);not null"`
	UserID         uint      `json:"user_id" gorm:"type:bigint;not null"`
	User           Users     `json:"user" gorm:"foreignkey:UserID"`
	CreatedAt      time.Time `json:"created_at" gorm:"type:timestamp"`
	UpdatedAt      time.Time `json:"updated_at" gorm:"type:timestamp"`
}

type Users struct {
	ID              uint            `json:"id" gorm:"primary_key;type:bigint"`
	Username        string          `json:"username" gorm:"unique;type:varchar(50);not null"`
	Email           string          `json:"email" form:"email" valid:"required, email" gorm:"unique;type:varchar(150);not null;uniqueIndex"`
	Password        string          `json:"password" form:"password" valid:"required,minstringlength(6)" gorm:"type:text;not null"`
	Age             uint            `json:"age" gorm:"type:int;not null"`
	ProfileImageURL string          `json:"profile_image_url" valid:"url" gorm:"type:text"`
	CreatedAt       time.Time       `json:"created_at" gorm:"type:timestamp"`
	UpdatedAt       time.Time       `json:"updated_at" gorm:"type:timestamp"`
	SocialMedias    []Social_Medias `json:"social_medias" gorm:"foreignkey:UserID"`
	Photos          []Photos        `json:"photos" gorm:"foreignkey:UserID"`
	Comments        []Comments      `json:"comments" gorm:"foreignkey:UserID"`
}

func (user *Users) Validate(db *gorm.DB) {
	if !strings.Contains(user.Email, "@") {
		db.AddError(errors.New("format email tidak sesuai"))
	}

	if user.Age <= 8 {
		db.AddError(errors.New("Minimal 9 tahun."))
	}
}

func (user *Users) BeforeCreate(tx *gorm.DB) (err error) {
	user.Validate(tx)
	_, errCreate := govalidator.ValidateStruct(user)
	if errCreate != nil {
		err = errCreate
		return
	}
	user.Password = helpers.HashDong(user.Password)
	return nil
}

// Kurang berfungsi

// func (user *Users) BeforeUpdate(tx *gorm.DB) (err error) {
// 	if user.Password != "" {
// 		hash := helpers.HashDong(user.Password)
// 		user.Password = hash
// 	}
// 	return nil
// }

type LoginRequest struct {
	Email    string `json:"email" valid:"required, email" binding:"required"`
	Password string `json:"password" valid:"required,minstringlength(6)" binding:"required"`
}
type Photos struct {
	ID        uint       `json:"id" gorm:"primary_key;type:bigint"`
	Title     string     `json:"title" gorm:"type:varchar(100);not null"`
	Caption   string     `json:"caption" gorm:"type:varchar(200);"`
	PhotoURL  string     `json:"photo_url" valid:"url" gorm:"type:text;not null"`
	UserID    uint       `json:"user_id" gorm:"type:bigint;not null"`
	CreatedAt time.Time  `json:"created_at" gorm:"type:timestamp"`
	UpdatedAt time.Time  `json:"updated_at" gorm:"type:timestamp"`
	Comments  []Comments `json:"comments" gorm:"foreignkey:PhotoID"`
}

func (user *Photos) BeforeCreate(tx *gorm.DB) (err error) {
	_, errCreate := govalidator.ValidateStruct(user)
	if errCreate != nil {
		err = errCreate
		return
	}
	return nil
}

type Comments struct {
	ID        uint      `json:"id" gorm:"primary_key;type:bigint"`
	UserID    uint      `json:"user_id" gorm:"type:bigint;not null"`
	PhotoID   uint      `json:"photo_id" gorm:"type:bigint;not null"`
	Message   string    `json:"message" gorm:"type:varchar(200);"`
	CreatedAt time.Time `json:"created_at" gorm:"type:timestamp"`
	UpdatedAt time.Time `json:"updated_at" gorm:"type:timestamp"`
}

func (user *Comments) BeforeCreate(tx *gorm.DB) (err error) {
	_, errCreate := govalidator.ValidateStruct(user)
	if errCreate != nil {
		err = errCreate
		return
	}
	return nil
}
