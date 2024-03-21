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

type CommentListResponse struct {
	ID      uint                 `json:"id"`
	Message string               `json:"message"`
	PhotoID uint                 `json:"photo_id"`
	UserID  uint                 `json:"user_id"`
	User    CommentUserResponse  `json:"user"`
	Photo   CommentPhotoResponse `json:"photo"`
}

type CommentUserResponse struct {
	ID       uint   `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
}

type CommentPhotoResponse struct {
	ID       uint   `json:"id"`
	Caption  string `json:"caption"`
	Title    string `json:"title"`
	PhotoURL string `json:"photo_url"`
	UserID   uint   `json:"user_id"`
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
}

func CommentList(c *gin.Context) {
	// isi photos dengan models photos
	var comments []models.Comments
	db := database.GetDB()
	//cari data photo di db
	db.Preload("User").Preload("Photo").Find(&comments)

	//custom resp
	var response []CommentListResponse
	// untuk photo didalam photos ulang sebanyak photos
	for _, comment := range comments {
		resp := CommentListResponse{
			ID:      comment.ID,
			Message: comment.Message,
			PhotoID: comment.PhotoID,
			UserID:  comment.UserID,
			User:    CommentUserResponse{ID: comment.User.ID, Email: comment.User.Email, Username: comment.User.Username},
			Photo:   CommentPhotoResponse{ID: comment.Photo.ID, Caption: comment.Photo.Caption, Title: comment.Photo.Title, PhotoURL: comment.Photo.PhotoURL, UserID: comment.Photo.UserID},
		}
		response = append(response, resp)
	}

	c.JSON(200, response)
}

func CommentListByID(c *gin.Context) {
	id := c.Param("id")
	var comment models.Comments
	db := database.GetDB()
	if err := db.Preload("User").Preload("Photo").First(&comment, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Comments not found"})
		return
	}
	response := CommentListResponse{
		ID:      comment.ID,
		Message: comment.Message,
		PhotoID: comment.PhotoID,
		UserID:  comment.UserID,
		User:    CommentUserResponse{ID: comment.User.ID, Email: comment.User.Email, Username: comment.User.Username},
		Photo:   CommentPhotoResponse{ID: comment.Photo.ID, Caption: comment.Photo.Caption, Title: comment.Photo.Title, PhotoURL: comment.Photo.PhotoURL, UserID: comment.Photo.UserID},
	}
	c.JSON(200, response)
}

func CommentUpdate(c *gin.Context) {
	id := c.Param("id")
	idConvert, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var comment models.Comments
	if err := c.BindJSON(&comment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userData := c.MustGet("userData").(jwt.MapClaims)
	userID := uint(userData["id"].(float64))

	_, errCreate := govalidator.ValidateStruct(comment)
	if errCreate != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": errCreate.Error()})
		return
	}

	db := database.GetDB()

	result := db.Model(&models.Comments{}).Where("user_id = ?", userID).Where("id = ?", idConvert).Updates(comment)

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Comment tidak ada yang berubah"})
		return
	}
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	if err := db.Preload("User").Preload("Photo").Where("user_id = ?", userID).Where("id = ?", idConvert).First(&comment).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Comment tidak ditemukan"})
		return
	}

	response := struct {
		ID      uint   `json:"id"`
		Message string `json:"message"`
		PhotoID uint   `json:"photo_id"`
		UserID  uint   `json:"user_id"`
	}{
		ID:      uint(idConvert),
		Message: comment.Message,
		PhotoID: comment.Photo.ID,
		UserID:  comment.User.ID,
	}
	c.JSON(200, response)
}

func CommentDelete(c *gin.Context) {
	id := c.Param("id")
	idConvert, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id"})
		return
	}
	userData, _ := c.Get("userData")
	userID := userData.(jwt.MapClaims)["id"].(float64)

	db := database.GetDB()
	var user models.Comments
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
