package models

import (
	"errors"
	"mygram/helpers"
	"time"

	"github.com/asaskevich/govalidator"
	"gorm.io/gorm"
)

type Users struct {
	ID              uint            `json:"id" gorm:"primary_key;type:bigint"`
	Username        string          `json:"username" valid:"alphanum,minstringlength(3)" gorm:"unique;type:varchar(50);not null; default:null"`
	Email           string          `json:"email" form:"email" valid:"email"  gorm:"unique;type:varchar(150);not null; default:null;uniqueIndex"`
	Password        string          `json:"password" form:"password" valid:"minstringlength(6)" gorm:"type:text;not null; default:null"`
	Age             uint            `json:"age" gorm: valid:"min:9" "type:int;not null; default:null"`
	ProfileImageURL string          `json:"profile_image_url" valid:"url" gorm:"type:text"`
	CreatedAt       time.Time       `json:"created_at" gorm:"type:timestamp"`
	UpdatedAt       time.Time       `json:"updated_at" gorm:"type:timestamp"`
	SocialMedias    []Social_Medias `json:"social_medias" gorm:"foreignkey:UserID; constraint:OnDelete:CASCADE"`
	Photos          []Photos        `json:"photos" gorm:"foreignkey:UserID; constraint:OnDelete:CASCADE"`
	Comments        []Comments      `json:"comments" gorm:"foreignkey:UserID; constraint:OnDelete:CASCADE"`
}

func (user *Users) BeforeCreate(tx *gorm.DB) (err error) {
	if user.Age <= 8 {
		return errors.New("minimal sembilan tahun")
	}
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
	Title     string     `json:"title" gorm:"type:varchar(100);not null; default:null"`
	Caption   string     `json:"caption" gorm:"type:varchar(200);"`
	PhotoURL  string     `json:"photo_url" valid:"url" gorm:"type:text;not null; default:null"`
	UserID    uint       `json:"user_id" gorm:"type:bigint;not null; default:null"`
	User      Users      `json:"user" gorm:"foreignKey:ID; references:UserID"`
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
	UserID    uint      `json:"user_id" gorm:"type:bigint;not null; default:null"`
	PhotoID   uint      `json:"photo_id" gorm:"type:bigint;not null; default:null"`
	Message   string    `json:"message" gorm:"type:varchar(200);"`
	CreatedAt time.Time `json:"created_at" gorm:"type:timestamp"`
	UpdatedAt time.Time `json:"updated_at" gorm:"type:timestamp"`
	User      Users     `json:"user" gorm:"foreignKey:ID; references:UserID"`
	Photo     Photos    `json:"photo" gorm:"foreignKey:ID; references:PhotoID"`
}

func (user *Comments) BeforeCreate(tx *gorm.DB) (err error) {
	_, errCreate := govalidator.ValidateStruct(user)
	if errCreate != nil {
		err = errCreate
		return
	}
	return nil
}

type Social_Medias struct {
	ID             uint      `json:"id" gorm:"primary_key;type:bigint"`
	Name           string    `json:"name" gorm:"type:varchar(50);not null; default:null"`
	SocialMediaURL string    `json:"social_media_url" valid:"url" gorm:"type:varchar(50);not null; default:null"`
	UserID         uint      `json:"user_id" gorm:"type:bigint;not null; default:null"`
	User           Users     `json:"user" gorm:"foreignkey:ID; references:UserID" `
	CreatedAt      time.Time `json:"created_at" gorm:"type:timestamp"`
	UpdatedAt      time.Time `json:"updated_at" gorm:"type:timestamp"`
}

func (user *Social_Medias) BeforeCreate(tx *gorm.DB) (err error) {
	_, errCreate := govalidator.ValidateStruct(user)
	if errCreate != nil {
		err = errCreate
		return
	}
	return nil
}
