package controllers

import (
    "net/http"
    "strconv"
    "uk-achmadxander/database"
    "uk-achmadxander/helpers"
    "uk-achmadxander/models"
    "github.com/gin-gonic/gin"
    "time"
)

func AddPhoto(c *gin.Context) {
    var photo models.Photo
    contentType := helpers.GetContentType(c)
    userIDInterface, exists := c.Get("user_id")

    if !exists {
        helpers.RespondError(c, http.StatusBadRequest, "User ID not found in context")
        return
    }

    // Convert userIDInterface to uint
    var userID uint
    switch v := userIDInterface.(type) {
    case string:
        userIDUint, err := strconv.ParseUint(v, 10, 64)
        if err != nil {
            helpers.RespondError(c, http.StatusBadRequest, "Invalid user ID format")
            return
        }
        userID = uint(userIDUint)
    default:
        helpers.RespondError(c, http.StatusBadRequest, "Invalid user ID format")
        return
    }

    // Bind JSON or form data based on content type
    if contentType == helpers.AppJSON {
        if err := c.ShouldBindJSON(&photo); err != nil {
            helpers.RespondError(c, http.StatusBadRequest, "Invalid JSON input: "+err.Error())
            return
        }
    } else {
        if err := c.ShouldBind(&photo); err != nil {
            helpers.RespondError(c, http.StatusBadRequest, "Invalid form data: "+err.Error())
            return
        }
    }

    // Set the UserID field in Photo struct
    photo.UserID = userID

    // Validate photo URL format using gorm tags
    if err := database.DB.Create(&photo).Error; err != nil {
        helpers.RespondError(c, http.StatusBadRequest, "Failed to create photo: "+err.Error())
        return
    }

    // Prepare response payload
    photoResponse := struct {
        ID        uint      `json:"id"`
        Title     string    `json:"title"`
        Caption   string    `json:"caption"`
        PhotoURL  string    `json:"photo_url"`
        UserID    uint      `json:"user_id"`
        CreatedAt time.Time `json:"created_at"`
    }{
        ID:        photo.ID,
        Title:     photo.Title,
        Caption:   photo.Caption,
        PhotoURL:  photo.URL,
        UserID:    photo.UserID,
        CreatedAt: *photo.CreatedAt,
    }

    helpers.RespondJSON(c, http.StatusCreated, photoResponse)
}

func GetAllPhotos(c *gin.Context) {
    var photos []models.Photo
    userIDInterface, exists := c.Get("user_id")

    if !exists {
        helpers.RespondError(c, http.StatusBadRequest, "User ID not found in context")
        return
    }

    // Convert userIDInterface to uint
    var userID uint
    switch v := userIDInterface.(type) {
    case string:
        userIDUint, err := strconv.ParseUint(v, 10, 64)
        if err != nil {
            helpers.RespondError(c, http.StatusBadRequest, "Invalid user ID format")
            return
        }
        userID = uint(userIDUint)
    default:
        helpers.RespondError(c, http.StatusBadRequest, "Invalid user ID format")
        return
    }

    // Query photos with User information preloaded
    if err := database.DB.Where("user_id = ?", userID).Preload("User").Find(&photos).Error; err != nil {
        helpers.RespondError(c, http.StatusInternalServerError, err.Error())
        return
    }

    type PhotoResponseAll struct {
        ID        uint      `json:"id"`
        Title     string    `json:"title"`
        Caption   string    `json:"caption"`
        PhotoURL  string    `json:"photo_url"`
        UserID    uint      `json:"user_id"`
        CreatedAt time.Time `json:"created_at"`
        UpdatedAt time.Time `json:"updated_at"`
        User      struct {
            Email    string `json:"email"`
            Username string `json:"username"`
        } `json:"user"`
    }

    var photosResponse []PhotoResponseAll
    for _, photo := range photos {
        photoResponse := PhotoResponseAll{
            ID:        photo.ID,
            Title:     photo.Title,
            Caption:   photo.Caption,
            PhotoURL:  photo.URL,
            UserID:    photo.UserID,
            CreatedAt: *photo.CreatedAt,
            UpdatedAt: *photo.UpdatedAt,
        }
        // Copy user information
        photoResponse.User.Email = photo.User.Email
        photoResponse.User.Username = photo.User.Username

        photosResponse = append(photosResponse, photoResponse)
    }

    helpers.RespondJSON(c, http.StatusOK, photosResponse)
}



func UpdatePhoto(c *gin.Context) {
    var photo models.Photo
    photoID, _ := strconv.Atoi(c.Param("photoId"))
    userID, _ := c.Get("user_id")

    if err := database.DB.Where("id = ? AND user_id = ?", photoID, userID).First(&photo).Error; err != nil {
        helpers.RespondError(c, http.StatusNotFound, "Photo not found or unauthorized")
        return
    }

    contentType := helpers.GetContentType(c)
    if contentType == helpers.AppJSON {
        c.ShouldBindJSON(&photo)
    } else {
        c.ShouldBind(&photo)
    }

    if err := database.DB.Save(&photo).Error; err != nil {
        helpers.RespondError(c, http.StatusInternalServerError, err.Error())
        return
    }

    type PhotoResponseUpdate struct {
        ID        uint       `json:"id"`
        Title     string     `json:"title"`
        Caption   string     `json:"caption"`
        PhotoURL  string     `json:"photo_url"`
        UserID    uint       `json:"user_id"`
        UpdatedAt *time.Time `json:"updated_at"`
    }

    photoResponse := PhotoResponseUpdate{
        ID:        photo.ID,
        Title:     photo.Title,
        Caption:   photo.Caption,
        PhotoURL:  photo.URL,
        UserID:    photo.UserID,
        UpdatedAt: photo.UpdatedAt,
    }

    helpers.RespondJSON(c, http.StatusOK, photoResponse)
}

func DeletePhoto(c *gin.Context) {
    var photo models.Photo
    photoID, _ := strconv.Atoi(c.Param("photoId"))
    userID, _ := c.Get("user_id")

    if err := database.DB.Where("id = ? AND user_id = ?", photoID, userID).First(&photo).Error; err != nil {
        helpers.RespondError(c, http.StatusNotFound, "Photo not found or unauthorized")
        return
    }

    if err := database.DB.Delete(&photo).Error; err != nil {
        helpers.RespondError(c, http.StatusInternalServerError, err.Error())
        return
    }

    helpers.RespondJSON(c, http.StatusOK, gin.H{"message": "Your photo has been successfully deleted"})
}

