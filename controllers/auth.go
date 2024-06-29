package controllers

import (
    "net/http"
    "uk-achmadxander/helpers"
    "uk-achmadxander/models"
    "uk-achmadxander/database"
    "github.com/gin-gonic/gin"
    "golang.org/x/crypto/bcrypt"
)

func Register(c *gin.Context) {
    contentType := helpers.GetContentType(c)
    user := models.User{}

    if contentType == helpers.AppJSON {
        if err := c.ShouldBindJSON(&user); err != nil {
            helpers.RespondError(c, http.StatusBadRequest, "Invalid JSON input")
            return
        }
    } else {
        if err := c.ShouldBind(&user); err != nil {
            helpers.RespondError(c, http.StatusBadRequest, "Invalid form data")
            return
        }
    }

    // Validasi manual untuk usia minimal
    if user.Age < 9 {
        helpers.RespondError(c, http.StatusBadRequest, "Age must be at least 9 years old")
        return
    }

    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
    if err != nil {
        helpers.RespondError(c, http.StatusInternalServerError, err.Error())
        return
    }
    user.Password = string(hashedPassword)

    if err := database.DB.Create(&user).Error; err != nil {
        helpers.RespondError(c, http.StatusBadRequest, err.Error())
        return
    }

    // Buat response JSON yang sesuai langsung di sini
    userResponse := gin.H{
        "id":       user.ID,
        "username": user.Username,
        "email":    user.Email,
        "age":      user.Age,
    }

    helpers.RespondJSON(c, http.StatusCreated, userResponse)
}

func Login(c *gin.Context) {
    contentType := helpers.GetContentType(c)
    user := models.User{}
    password := ""

    if contentType == helpers.AppJSON {
        c.ShouldBindJSON(&user)
    } else {
        c.ShouldBind(&user)
    }

    password = user.Password

    if err := database.DB.Where("email = ?", user.Email).First(&user).Error; err != nil {
        helpers.RespondError(c, http.StatusUnauthorized, "Invalid email or password")
        return
    }

    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
        helpers.RespondError(c, http.StatusUnauthorized, "Invalid email or password")
        return
    }

    token, _ := helpers.GenerateToken(user.ID, user.Email)

    helpers.RespondJSON(c, http.StatusOK, gin.H{"token": token})
}
