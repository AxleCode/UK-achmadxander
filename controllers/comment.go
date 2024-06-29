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

type CommentResponse struct {
    ID        uint          `json:"id"`
    Message   string        `json:"message"`
    PhotoID   uint          `json:"photo_id"`
    UserID    uint          `json:"user_id"`
    CreatedAt *time.Time    `json:"created_at,omitempty"`
    UpdatedAt *time.Time    `json:"updated_at,omitempty"`
    User      UserResponse  `json:"User"`
    Photo     PhotoResponse `json:"Photo"`
}

type UserResponse struct {
    ID       uint   `json:"id"`
    Email    string `json:"email"`
    Username string `json:"username"`
}

type PhotoResponse struct {
    ID       uint   `json:"id"`
    Title    string `json:"title"`
    Caption  string `json:"caption"`
    PhotoURL string `json:"photo_url"`
    UserID   uint   `json:"user_id"`
}

func AddComment(c *gin.Context) {
    var comment models.Comment
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

    // Bind JSON input
    if err := c.ShouldBindJSON(&comment); err != nil {
        helpers.RespondError(c, http.StatusBadRequest, "Invalid JSON input")
        return
    }

    // Set user ID from JWT token
    comment.UserID = userID

    // Validate photo existence
    var photo models.Photo
    if err := database.DB.First(&photo, comment.PhotoID).Error; err != nil {
        helpers.RespondError(c, http.StatusNotFound, "Photo not found")
        return
    }

    // Save comment to database
    if err := database.DB.Create(&comment).Error; err != nil {
        helpers.RespondError(c, http.StatusInternalServerError, "Failed to add comment")
        return
    }

    // Prepare response payload
    commentResponse := struct {
        ID        uint       `json:"id"`
        Message   string     `json:"message"`
        PhotoID   uint       `json:"photo_id"`
        UserID    uint       `json:"user_id"`
        CreatedAt *time.Time `json:"created_at"`
    }{
        ID:        comment.ID,
        Message:   comment.Message,
        PhotoID:   comment.PhotoID,
        UserID:    comment.UserID,
        CreatedAt: comment.CreatedAt,
    }

    helpers.RespondJSON(c, http.StatusCreated, commentResponse)
}

func GetAllComments(c *gin.Context) {
    var comments []models.Comment

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

    // Fetch comments related to the user
    if err := database.DB.Where("user_id = ?", userID).Preload("User").Preload("Photo").Find(&comments).Error; err != nil {
        helpers.RespondError(c, http.StatusInternalServerError, err.Error())
        return
    }
   
    var response []CommentResponse
    for _, comment := range comments {
        commentResponse := CommentResponse{
            ID:        comment.ID,
            Message:   comment.Message,
            PhotoID:   comment.PhotoID,
            UserID:    comment.UserID,
            CreatedAt: comment.CreatedAt,
            UpdatedAt: comment.UpdatedAt,
            User: UserResponse{
                ID:       comment.User.ID,
                Email:    comment.User.Email,
                Username: comment.User.Username,
            },
            Photo: PhotoResponse{
                ID:        comment.Photo.ID,
                Title:     comment.Photo.Title,
                Caption:   comment.Photo.Caption,
                PhotoURL:  comment.Photo.URL,
                UserID:    comment.Photo.UserID,
            },
        }
        response = append(response, commentResponse)
    }

    helpers.RespondJSON(c, http.StatusOK, response)
}

func UpdateComment(c *gin.Context) {
    var comment models.Comment
    commentID, err := strconv.Atoi(c.Param("commentId"))
    if err != nil {
        helpers.RespondError(c, http.StatusBadRequest, "Invalid comment ID")
        return
    }

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

    // Check if the comment exists and belongs to the authenticated user
    if err := database.DB.Where("id = ? AND user_id = ?", commentID, userID).First(&comment).Error; err != nil {
        helpers.RespondError(c, http.StatusNotFound, "Comment not found or unauthorized")
        return
    }

    // Bind JSON or form data based on content type
    contentType := helpers.GetContentType(c)
    if contentType == helpers.AppJSON {
        if err := c.ShouldBindJSON(&comment); err != nil {
            helpers.RespondError(c, http.StatusBadRequest, "Invalid JSON input: "+err.Error())
            return
        }
    } else {
        if err := c.ShouldBind(&comment); err != nil {
            helpers.RespondError(c, http.StatusBadRequest, "Invalid form data: "+err.Error())
            return
        }
    }

    // Save the updated comment
    if err := database.DB.Save(&comment).Error; err != nil {
        helpers.RespondError(c, http.StatusInternalServerError, "Failed to update comment: "+err.Error())
        return
    }

    // Prepare the response payload
    response := gin.H{
        "id":         comment.ID,
        "message":    comment.Message,
        "photo_id":   comment.PhotoID,
        "user_id":    comment.UserID,
        "updated_at": comment.UpdatedAt.Format(time.RFC3339Nano),
    }

    // Respond with the updated comment
    helpers.RespondJSON(c, http.StatusOK, response)
}




func DeleteComment(c *gin.Context) {
    var comment models.Comment
    commentID, _ := strconv.Atoi(c.Param("commentId"))
    userID, _ := c.Get("user_id")

    if err := database.DB.Where("id = ? AND user_id = ?", commentID, userID).First(&comment).Error; err != nil {
        helpers.RespondError(c, http.StatusNotFound, "Comment not found or unauthorized")
        return
    }

    if err := database.DB.Delete(&comment).Error; err != nil {
        helpers.RespondError(c, http.StatusInternalServerError, err.Error())
        return
    }

    helpers.RespondJSON(c, http.StatusOK, gin.H{"message": "Comment deleted successfully"})
}
