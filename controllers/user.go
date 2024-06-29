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

func UpdateUser(c *gin.Context) {
    var user models.User
    userID := c.GetString("user_id") // Get user_id from context

    userIDUint64, err := strconv.ParseUint(userID, 10, 64)
    if err != nil {
        helpers.RespondError(c, http.StatusUnauthorized, "Invalid user ID")
        return
    }

    if err := database.DB.First(&user, userIDUint64).Error; err != nil {
        helpers.RespondError(c, http.StatusNotFound, "User not found")
        return
    }

    // Ensure user is allowed to update their own profile
    if strconv.FormatUint(userIDUint64, 10) != strconv.FormatUint(uint64(user.ID), 10) {
        helpers.RespondError(c, http.StatusUnauthorized, "Unauthorized")
        return
    }

    // Bind request body to user model
    if err := c.ShouldBindJSON(&user); err != nil {
        helpers.RespondError(c, http.StatusBadRequest, "Invalid request body")
        return
    }

    // Save updated user to database
    if err := database.DB.Save(&user).Error; err != nil {
        helpers.RespondError(c, http.StatusInternalServerError, err.Error())
        return
    }

    type UserResponse struct {
        ID        uint      `json:"id"`
        Email     string    `json:"email"`
        Username  string    `json:"username"`
        Age       uint      `json:"age"`
        UpdatedAt time.Time `json:"updated_at"`
    }

    userResponse := UserResponse{
        ID:        user.ID,
        Email:     user.Email,
        Username:  user.Username,
        Age:       user.Age,
        UpdatedAt: time.Now(),
    }

    // Respond with updated user information
    helpers.RespondJSON(c, http.StatusOK, userResponse)
}

func DeleteUser(c *gin.Context) {
    var user models.User

    // Retrieve user_id from context and convert it to uint64
    currentUserID := c.GetString("user_id")
    currentUserIDUint64, err := strconv.ParseUint(currentUserID, 10, 64)
    if err != nil {
        helpers.RespondError(c, http.StatusUnauthorized, "Invalid user ID")
        return
    }

    if err := database.DB.First(&user, currentUserIDUint64).Error; err != nil {
        helpers.RespondError(c, http.StatusNotFound, "User not found")
        return
    }

    if err := database.DB.Delete(&user).Error; err != nil {
        helpers.RespondError(c, http.StatusInternalServerError, err.Error())
        return
    }

    helpers.RespondJSON(c, http.StatusOK, gin.H{"message": "Your account has been successfully deleted"})
}
