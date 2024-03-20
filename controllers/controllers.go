package controllers

import (
	"mygram/database"
	"mygram/helpers"
	"mygram/models"
	"net/http"
	"strconv"

	"github.com/asaskevich/govalidator"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func UserRegister(c *gin.Context) {
	var newUser models.Users

	if err := c.BindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := database.GetDB()
	if !govalidator.IsEmail(newUser.Email) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email format"})
		return
	}

	if !govalidator.StringLength(newUser.Password, "6", "1000") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Minimal password 6 karakter."})
		return
	}

	if err := db.Create(&newUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := struct {
		ID              uint   `json:"id"`
		Email           string `json:"email"`
		Username        string `json:"username"`
		Age             uint   `json:"age"`
		ProfileImageURL string `json:"profile_image_url"`
	}{
		ID:              newUser.ID,
		Email:           newUser.Email,
		Username:        newUser.Username,
		Age:             newUser.Age,
		ProfileImageURL: newUser.ProfileImageURL,
	}

	c.JSON(201, response)

}

func UserLogin(c *gin.Context) {
	var request models.LoginRequest
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := database.GetDB()
	var user models.Users
	if err := db.Where("email = ?", request.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Unauthorized",
			"message": "invalid email/password",
		})
		return
	}

	if !helpers.CompareDong([]byte(user.Password), []byte(request.Password)) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Unauthorized",
			"message": "invalid email/password",
		})
		return
	}
	//GenerateJwtToken(id uint, email string, username string, age int, pp string, created_at string, updated_at string)
	token := helpers.GenerateJwtToken(user.ID, user.Email, user.Username, user.Age, user.ProfileImageURL, user.CreatedAt, user.UpdatedAt)
	c.JSON(200, gin.H{"token": token})
}

func UserUpdate(c *gin.Context) {
	var request models.Users
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userData := c.MustGet("userData").(jwt.MapClaims)
	userID := uint(userData["id"].(float64))

	if request.Password != "" {
		request.Password = helpers.HashDong(request.Password)
	}

	db := database.GetDB()
	if err := db.Model(&models.Users{}).Where("id = ?", userID).Updates(request).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	response := struct {
		ID              uint   `json:"id"`
		Email           string `json:"email"`
		Username        string `json:"username"`
		Age             uint   `json:"age"`
		ProfileImageURL string `json:"profile_image_url"`
	}{
		ID:              request.ID,
		Email:           request.Email,
		Username:        request.Username,
		Age:             request.Age,
		ProfileImageURL: request.ProfileImageURL,
	}
	c.JSON(200, response)
}

func UserDelete(c *gin.Context) {
	userData, _ := c.Get("userData")
	userID := userData.(jwt.MapClaims)["id"].(float64)

	db := database.GetDB()
	var user models.Users
	if err := db.Where("id = ?", userID).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Data tidak ada!"})
		return
	}

	if err := db.Delete(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus data!"})
		return
	}

	c.JSON(200, "OK")
}

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

func PhotoList(c *gin.Context) {
	// isi photos dengan models photos
	var photos []models.Photos
	db := database.GetDB()
	//cari data photo di db
	db.Preload("User").Find(&photos)

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

func CommentPost(c *gin.Context) {
	userData, _ := c.Get("userData")
	userDataID := uint(userData.(jwt.MapClaims)["id"].(float64))

	var newComment models.Comments

	if err := c.BindJSON(&newComment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	newComment.UserID = userDataID
	db := database.GetDB()

	if err := db.Create(&newComment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := struct {
		ID      uint   `json:"id"`
		Message string `json:"message"`
		PhotoID uint   `json:"photo_id"`
		UserID  uint   `json:"user_id"`
	}{
		ID:      newComment.ID,
		Message: newComment.Message,
		PhotoID: newComment.PhotoID,
		UserID:  newComment.UserID,
	}

	c.JSON(201, response)

	// type Comments struct {
	// 	ID        uint      `json:"id" gorm:"primary_key;type:bigint"`
	// 	UserID    uint      `json:"user_id" gorm:"type:bigint;not null"`
	// 	PhotoID   uint      `json:"photo_id" gorm:"type:bigint;not null"`
	// 	Message   string    `json:"message" gorm:"type:varchar(200);"`
	// 	CreatedAt time.Time `json:"created_at" gorm:"type:timestamp"`
	// 	UpdatedAt time.Time `json:"updated_at" gorm:"type:timestamp"`
	// }

}

// jika tidak pakai validator

// func validasiEmail(email string) error {
// 	if email == "" {
// 		return errors.New("email tidak boleh kosong")
// 	}
// 	// cek jika tidak ada @ dan .
// 	if !strings.Contains(email, "@") || !strings.Contains(email, ".") {
// 		return errors.New("format email tidak lazim tolong gunakan xxx@xxx.xxx")
// 	}
// 	db := database.GetDB()
// 	var count int64
// 	db.Model(&models.Users{}).Where("email = ?", email).Count(&count)
// 	if count > 0 {
// 		return errors.New("email telah terdaftar")
// 	}

// 	return nil
// }

// func validasiUsername(username string) error {
// 	if username == "" {
// 		return errors.New("username tidak boleh kosong")
// 	}
// 	db := database.GetDB()
// 	var count int64
// 	db.Model(&models.Users{}).Where("username = ?", username).Count(&count)
// 	if count > 0 {
// 		return errors.New("username telah terdaftar")
// 	}

// 	return nil
// }

// func validasiPassword(password string) error {
// 	if password == "" {
// 		return errors.New("password cannot be empty")
// 	}

// 	if len(password) < 6 {
// 		return errors.New("password must be at least 6 characters long")
// 	}

// 	return nil
// }

// func validasiAge(age uint) error {
// 	if age <= 8 {
// 		return errors.New("umur minimal 9 tahun.")
// 	}

// 	return nil
// }

// func validasiProfileImage(pp string) error {
// 	u, err := url.Parse(pp)
// 	if err != nil {
// 		return errors.New("profile Image URL tidak lazim")
// 	}
// 	if u.Scheme != "http" && u.Scheme != "https" {
// 		return errors.New("sementara ini hanya menerima http atau https tidak menerima yang lain")
// 	}

// 	return nil
// }

// kalo gapake validator func UserRegister(c *gin.Context) { }

// if err := validasiEmail(newUser.Email); err != nil {
// 	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 	return
// }

// if err := validasiUsername(newUser.Username); err != nil {
// 	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 	return
// }

// if err := validasiPassword(newUser.Password); err != nil {
// 	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 	return
// }

// if err := validasiAge(newUser.Age); err != nil {
// 	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 	return
// }

// if err := validasiProfileImage(newUser.ProfileImageURL); err != nil {
// 	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 	return
// }
