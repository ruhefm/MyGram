package controllers

import (
	"mygram/database"
	"mygram/models"
	"net/http"
	"strconv"

	"github.com/asaskevich/govalidator"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func PhotoUpload(c *gin.Context) {
	userData, _ := c.Get("userData")
	userDataID := uint(userData.(jwt.MapClaims)["id"].(float64))

	var newPhoto models.Photos

	if err := c.BindJSON(&newPhoto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	newPhoto.UserID = userDataID
	db := database.GetDB()

	if err := db.Create(&newPhoto).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := struct {
		ID       uint   `json:"id"`
		Caption  string `json:"caption"`
		Title    string `json:"title"`
		PhotoURL string `json:"photo_url"`
		UserID   uint   `json:"user_id"`
	}{
		ID:       newPhoto.ID,
		Caption:  newPhoto.Caption,
		Title:    newPhoto.Title,
		PhotoURL: newPhoto.PhotoURL,
		UserID:   userDataID,
	}

	c.JSON(201, response)

	// type Photos struct {
	// 	ID        uint      `json:"id" gorm:"primary_key;type:bigint"`
	// 	Title     string    `json:"title" gorm:"type:varchar(100);not null"`
	// 	Caption   string    `json:"caption" gorm:"type:varchar(200);"`
	// 	PhotoURL  string    `json:"photo_url" valid:"url" gorm:"type:text;not null"`
	// 	UserID    uint      `json:"user_id" gorm:"type:bigint;not null"`
	// 	User      Users     `json:"user" gorm:"foreignkey:UserID;references:ID"`
	// 	CreatedAt time.Time `json:"created_at" gorm:"type:timestamp"`
	// 	UpdatedAt time.Time `json:"updated_at" gorm:"type:timestamp"`
	// }

}

type PhotoListResponse struct {
	ID       uint         `json:"id"`
	Caption  string       `json:"caption"`
	Title    string       `json:"title"`
	PhotoURL string       `json:"photo_url"`
	UserID   uint         `json:"user_id"`
	User     UserResponse `json:"user"`
}

type UserResponse struct {
	ID       uint   `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
}

type PhotoResponse struct {
	ID       uint   `json:"id"`
	Caption  string `json:"caption"`
	Title    string `json:"title"`
	PhotoURL string `json:"photo_url"`
	UserID   uint   `json:"user_id"`
}

func PhotoList(c *gin.Context) {
	// isi photos dengan models photos
	var photos []models.Photos
	db := database.GetDB()
	//cari data photo di db
	db.Preload("User").Find(&photos)
	if len(photos) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "Tidak ada data"})
		return
	}
	//custom resp
	var response []PhotoListResponse
	// untuk photo didalam photos ulang sebanyak photos
	for _, photo := range photos {
		resp := PhotoListResponse{
			ID:       photo.ID,
			Caption:  photo.Caption,
			Title:    photo.Title,
			PhotoURL: photo.PhotoURL,
			UserID:   photo.UserID,
			User:     UserResponse{ID: photo.User.ID, Email: photo.User.Email, Username: photo.User.Username},
		}
		response = append(response, resp)
	}

	c.JSON(200, response)
}

func PhotoListByID(c *gin.Context) {
	id := c.Param("id")
	var photos models.Photos
	db := database.GetDB()
	if err := db.Preload("User").First(&photos, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Photo not found"})
		return
	}
	response := PhotoListResponse{
		ID:       photos.ID,
		Caption:  photos.Caption,
		Title:    photos.Title,
		PhotoURL: photos.PhotoURL,
		UserID:   photos.UserID,
		User:     UserResponse{ID: photos.User.ID, Email: photos.User.Email, Username: photos.User.Username},
	}
	c.JSON(200, response)
}

func PhotoUpdate(c *gin.Context) {
	id := c.Param("id")
	idConvert, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var request models.Photos
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userData := c.MustGet("userData").(jwt.MapClaims)
	userID := uint(userData["id"].(float64))

	_, errCreate := govalidator.ValidateStruct(request)
	if errCreate != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": errCreate.Error()})
		return
	}

	db := database.GetDB()
	result := db.Model(&models.Photos{}).Where("user_id = ?", userID).Where("id = ?", idConvert).Updates(request)
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Photo tidak ditemukan"})
		return
	}
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	response := struct {
		ID       int    `json:"id"`
		Caption  string `json:"caption"`
		Title    string `json:"title"`
		PhotoURL string `json:"photo_url"`
		UserID   uint   `json:"user_id"`
	}{
		ID:       idConvert,
		Caption:  request.Caption,
		Title:    request.Title,
		PhotoURL: request.PhotoURL,
		UserID:   userID,
	}
	c.JSON(200, response)
}

func PhotoDelete(c *gin.Context) {
	id := c.Param("id")
	idConvert, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id"})
		return
	}
	userData, _ := c.Get("userData")
	userID := userData.(jwt.MapClaims)["id"].(float64)

	db := database.GetDB()
	var user models.Photos
	if err := db.Where("user_id = ?", userID).Where("id = ?", idConvert).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Data tidak ada!"})
		return
	}

	if err := db.Delete(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus data!"})
		return
	}

	c.JSON(200, "OK")
}
