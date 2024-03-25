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

type SocialListResponse struct {
	ID             uint               `json:"id"`
	Name           string             `json:"name"`
	SocialMediaURL string             `json:"social_media_url"`
	UserID         uint               `json:"user_id"`
	User           SocialUserResponse `json:"user"`
}

type SocialUserResponse struct {
	ID       uint   `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
}

func SocialPost(c *gin.Context) {
	userData, _ := c.Get("userData")
	userDataID := uint(userData.(jwt.MapClaims)["id"].(float64))

	var newComment models.Social_Medias

	if err := c.Bind(&newComment); err != nil {
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
		ID             uint   `json:"id"`
		Name           string `json:"name"`
		SocialMediaURL string `json:"social_media_url"`
		UserID         uint   `json:"user_id"`
	}{
		ID:             newComment.ID,
		Name:           newComment.Name,
		SocialMediaURL: newComment.SocialMediaURL,
		UserID:         newComment.UserID,
	}

	c.JSON(201, response)
}

func SocialList(c *gin.Context) {
	// isi photos dengan models photos
	var comments []models.Social_Medias
	db := database.GetDB()
	//cari data photo di db
	db.Preload("User").Find(&comments)
	if len(comments) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "Tidak ada data"})
		return
	}
	//custom resp
	var response []SocialListResponse
	// untuk photo didalam photos ulang sebanyak photos
	for _, comment := range comments {
		resp := SocialListResponse{
			ID:             comment.ID,
			Name:           comment.Name,
			SocialMediaURL: comment.SocialMediaURL,
			UserID:         comment.UserID,
			User:           SocialUserResponse{ID: comment.User.ID, Email: comment.User.Email, Username: comment.User.Username},
		}
		response = append(response, resp)
	}

	c.JSON(200, response)
}

func SocialListByID(c *gin.Context) {
	id := c.Param("id")
	var comment models.Social_Medias
	db := database.GetDB()
	if err := db.Preload("User").First(&comment, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Social Media not found"})
		return
	}
	response := SocialListResponse{
		ID:             comment.ID,
		Name:           comment.Name,
		SocialMediaURL: comment.SocialMediaURL,
		UserID:         comment.UserID,
		User:           SocialUserResponse{ID: comment.User.ID, Email: comment.User.Email, Username: comment.User.Username}}
	c.JSON(200, response)
}

func SocialUpdate(c *gin.Context) {
	id := c.Param("id")
	idConvert, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var comment models.Social_Medias
	if err := c.Bind(&comment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userData := c.MustGet("userData").(jwt.MapClaims)
	userID := uint(userData["id"].(float64))
	db := database.GetDB()
	_, errCreate := govalidator.ValidateStruct(comment)
	if errCreate != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": errCreate.Error()})
		return
	}
	result := db.Model(&models.Social_Medias{}).Where("user_id = ?", userID).Where("id = ?", idConvert).Updates(comment)

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Social Media tidak ada yang berubah"})
		return
	}
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	if err := db.Preload("User").Where("user_id = ?", userID).Where("id = ?", idConvert).First(&comment).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Social Media tidak ditemukan"})
		return
	}

	response := struct {
		ID             uint   `json:"id"`
		Name           string `json:"name"`
		SocialMediaURL string `json:"social_media_url"`
		UserID         uint   `json:"user_id"`
	}{
		ID:             uint(idConvert),
		Name:           comment.Name,
		SocialMediaURL: comment.SocialMediaURL,
		UserID:         comment.UserID,
	}
	c.JSON(200, response)
}

func SocialDelete(c *gin.Context) {
	id := c.Param("id")
	idConvert, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id"})
		return
	}
	userData, _ := c.Get("userData")
	userID := userData.(jwt.MapClaims)["id"].(float64)

	db := database.GetDB()
	var user models.Social_Medias
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
