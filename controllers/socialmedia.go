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

func AddSocialMedia(c *gin.Context) {
    var socialMedia models.SocialMedia

    // Bind JSON data
    if err := c.ShouldBindJSON(&socialMedia); err != nil {
        helpers.RespondError(c, http.StatusBadRequest, "Invalid JSON input: "+err.Error())
        return
    }

    // Get user ID from JWT token
    userIDInterface, exists := c.Get("user_id")
    if !exists {
        helpers.RespondError(c, http.StatusBadRequest, "User ID not found in context")
        return
    }

    var userID uint
    switch v := userIDInterface.(type) {
    case float64:
        userID = uint(v)
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

    // Set user ID to the social media entry
    socialMedia.UserID = userID

    // Create the social media entry in the database
    if err := database.DB.Create(&socialMedia).Error; err != nil {
        helpers.RespondError(c, http.StatusInternalServerError, "Failed to create social media: "+err.Error())
        return
    }

    // Prepare the response payload
    response := gin.H{
        "id":              socialMedia.ID,
        "name":            socialMedia.Name,
        "socialmedia_url": socialMedia.SocialMediaURL,
        "user_id":         socialMedia.UserID,
        "updated_at":      socialMedia.UpdatedAt.Format(time.RFC3339Nano),
    }

    // Respond with the created social media entry
    helpers.RespondJSON(c, http.StatusCreated, response)
}

func GetAllSocialMedia(c *gin.Context) {
    var socialMedias []models.SocialMedia

    // Get user ID from JWT token
    userIDInterface, exists := c.Get("user_id")
    if !exists {
        helpers.RespondError(c, http.StatusBadRequest, "User ID not found in context")
        return
    }

    // Convert userIDInterface to uint
    var userID uint
    switch v := userIDInterface.(type) {
    case float64:
        userID = uint(v)
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

    // Fetch social medias related to the user
    if err := database.DB.Where("user_id = ?", userID).Find(&socialMedias).Error; err != nil {
        helpers.RespondError(c, http.StatusInternalServerError, err.Error())
        return
    }

    // Prepare the response payload
    type SocialMediaResponse struct {
        ID             uint      `json:"id"`
        Name           string    `json:"name"`
        SocialMediaURL string    `json:"social_media_url"`
        UserID         uint      `json:"user_id"`
        CreatedAt      time.Time `json:"created_at"`
        UpdatedAt      time.Time `json:"updated_at"`
        User           struct {
            ID       uint   `json:"id"`
            Username string `json:"username"`
        } `json:"User"`
    }

    var response = struct {
        SocialMedias []SocialMediaResponse `json:"social_medias"`
    }{}

    for _, socialMedia := range socialMedias {
        // Fetch the user associated with this social media
        var user models.User
        if err := database.DB.Model(&models.User{}).Where("id = ?", socialMedia.UserID).First(&user).Error; err != nil {
            helpers.RespondError(c, http.StatusInternalServerError, "Failed to find user: "+err.Error())
            return
        }

        // Populate the response payload
        response.SocialMedias = append(response.SocialMedias, SocialMediaResponse{
            ID:             socialMedia.ID,
            Name:           socialMedia.Name,
            SocialMediaURL: socialMedia.SocialMediaURL,
            UserID:         socialMedia.UserID,
            CreatedAt:      socialMedia.CreatedAt,
            UpdatedAt:      socialMedia.UpdatedAt,
            User: struct {
                ID       uint   `json:"id"`
                Username string `json:"username"`
            }{
                ID:       user.ID,
                Username: user.Username,
            },
        })
    }

    // Respond with the social medias
    helpers.RespondJSON(c, http.StatusOK, response)
}

func UpdateSocialMedia(c *gin.Context) {
    var socialMedia models.SocialMedia

    // Convert socialMediaId to integer
    socialMediaID, err := strconv.Atoi(c.Param("socialMediaId"))
    if err != nil {
        helpers.RespondError(c, http.StatusBadRequest, "Invalid social media ID")
        return
    }

    // Get user ID from context
    userIDInterface, exists := c.Get("user_id")
    if !exists {
        helpers.RespondError(c, http.StatusBadRequest, "User ID not found in context")
        return
    }

    var userID uint
    switch v := userIDInterface.(type) {
    case float64:
        userID = uint(v)
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

    // Find the social media entry in the database
    if err := database.DB.Where("id = ? AND user_id = ?", socialMediaID, userID).First(&socialMedia).Error; err != nil {
        helpers.RespondError(c, http.StatusNotFound, "Social media not found or unauthorized")
        return
    }

    // Bind JSON data
    if err := c.ShouldBindJSON(&socialMedia); err != nil {
        helpers.RespondError(c, http.StatusBadRequest, "Invalid JSON input: "+err.Error())
        return
    }

    // Update the social media entry in the database
    if err := database.DB.Save(&socialMedia).Error; err != nil {
        helpers.RespondError(c, http.StatusInternalServerError, "Failed to update social media: "+err.Error())
        return
    }

    // Prepare the response payload in the desired order
    type SocialMediaResponse struct {
        ID             uint      `json:"id"`
        Name           string    `json:"name"`
        SocialMediaURL string    `json:"social_media_url"`
        UserID         uint      `json:"user_id"`
        UpdatedAt      time.Time `json:"updated_at"`
    }

    socialMediaResponse := SocialMediaResponse{
        ID:             socialMedia.ID,
        Name:           socialMedia.Name,
        SocialMediaURL: socialMedia.SocialMediaURL,
        UserID:         socialMedia.UserID,
        UpdatedAt:      socialMedia.UpdatedAt,
    }

    // Respond with the updated social media entry
    helpers.RespondJSON(c, http.StatusOK, socialMediaResponse)
}



func DeleteSocialMedia(c *gin.Context) {
    var socialMedia models.SocialMedia
    socialMediaID, _ := strconv.Atoi(c.Param("socialMediaId"))
    userID, _ := c.Get("user_id")

    if err := database.DB.Where("id = ? AND user_id = ?", socialMediaID, userID).First(&socialMedia).Error; err != nil {
        helpers.RespondError(c, http.StatusNotFound, "Social media not found or unauthorized")
        return
    }

    if err := database.DB.Delete(&socialMedia).Error; err != nil {
        helpers.RespondError(c, http.StatusInternalServerError, err.Error())
        return
    }

    helpers.RespondJSON(c, http.StatusOK, gin.H{"message": "Your social media has been successfully deleted"})
}
