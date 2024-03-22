package controllers

import (
	"mygram/database"
	"mygram/helpers"
	"mygram/models"
	"net/http"

	"github.com/asaskevich/govalidator"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type UsersListResponse struct {
	ID              uint                `json:"id" gorm:"primary_key;type:bigint"`
	Username        string              `json:"username" valid:"alphanum,minstringlength(3)"`
	Email           string              `json:"email" form:"email" valid:"email"`
	Age             uint                `json:"age"`
	ProfileImageURL string              `json:"profile_image_url"`
	SocialMedias    SocialMediaResponse `json:"social_medias"`
	Photos          PhotosResponse      `json:"photos" `
	Comments        CommentsResponse    `json:"comments" `
}

type SocialMediaResponse struct {
	ID             uint   `json:"id"`
	Name           string `json:"name"`
	SocialMediaURL string `json:"social_media_url"`
	UserID         uint   `json:"user_id"`
}

type PhotosResponse struct {
	ID       uint         `json:"id"`
	Caption  string       `json:"caption"`
	Title    string       `json:"title"`
	PhotoURL string       `json:"photo_url"`
	UserID   uint         `json:"user_id"`
	User     UserResponse `json:"user"`
}

type CommentsResponse struct {
	ID      uint                 `json:"id"`
	Message string               `json:"message"`
	PhotoID uint                 `json:"photo_id"`
	UserID  uint                 `json:"user_id"`
	User    CommentUserResponse  `json:"user"`
	Photo   CommentPhotoResponse `json:"photo"`
}

func UserRegister(c *gin.Context) {
	var newUser models.Users

	if err := c.Bind(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := database.GetDB()
	// untuk custom message
	// if !govalidator.IsEmail(newUser.Email) {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email format"})
	// 	return
	// }

	// if !govalidator.StringLength(newUser.Password, "6", "1000") {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "Minimal password 6 karakter."})
	// 	return
	// }
	if newUser.Age <= 8 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "minimal sembilan tahun."})
		return
	}
	_, errCreate := govalidator.ValidateStruct(newUser)
	if errCreate != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": errCreate.Error()})
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
	if err := c.Bind(&request); err != nil {
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
	userPassword := user.Password
	requestPassword := request.Password
	if !helpers.CompareDong([]byte(userPassword), []byte(requestPassword)) {
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
	if err := c.Bind(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userData := c.MustGet("userData").(jwt.MapClaims)
	userID := uint(userData["id"].(float64))
	db := database.GetDB()
	_, errCreate := govalidator.ValidateStruct(request)
	if errCreate != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": errCreate.Error()})
		return
	}
	if request.Age > 0 && request.Age <= 8 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "minimal sembilan tahun."})
		return
	}
	if request.Password != "" {
		request.Password = helpers.HashDong(request.Password)
	}
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
		ID:              userID,
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

	tx := db.Begin()

	if err := tx.Where("user_id = ?", user.ID).Delete(&[]models.Comments{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus komentar!"})
		return
	}

	if err := tx.Where("user_id = ?", user.ID).Delete(&[]models.Photos{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus foto!"})
		return
	}

	if err := tx.Where("user_id = ?", user.ID).Delete(&[]models.Social_Medias{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus media sosial!"})
		return
	}

	if err := tx.Delete(&user).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus pengguna!"})
		return
	}

	tx.Commit()

	c.JSON(200, "OK")
}

// type Comments struct {
// 	ID        uint      `json:"id" gorm:"primary_key;type:bigint"`
// 	UserID    uint      `json:"user_id" gorm:"type:bigint;not null"`
// 	PhotoID   uint      `json:"photo_id" gorm:"type:bigint;not null"`
// 	Message   string    `json:"message" gorm:"type:varchar(200);"`
// 	CreatedAt time.Time `json:"created_at" gorm:"type:timestamp"`
// 	UpdatedAt time.Time `json:"updated_at" gorm:"type:timestamp"`
// }

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
