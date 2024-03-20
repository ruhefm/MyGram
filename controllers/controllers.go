package controllers

import (
	"mygram/database"
	"mygram/helpers"
	"mygram/models"
	"net/http"

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
		c.JSON(http.StatusNotFound, gin.H{"error": "Record not found!"})
		return
	}

	if err := db.Delete(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete data!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Success delete"})
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
