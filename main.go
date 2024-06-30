package main

import (
    "os"
    "uk-achmadxander/controllers"
    "uk-achmadxander/database"
    "uk-achmadxander/middlewares"
    "github.com/gin-gonic/gin"
)

func main() {
    r := gin.Default()
    database.Init()

    userRoutes := r.Group("/users")
    {
        userRoutes.POST("/register", controllers.Register)
        userRoutes.POST("/login", controllers.Login)
        userRoutes.PUT("/", middlewares.Authentication(), controllers.UpdateUser)
        userRoutes.DELETE("/", middlewares.Authentication(), controllers.DeleteUser)
    }

    photoRoutes := r.Group("/photos")
    {
        photoRoutes.POST("/", middlewares.Authentication(), controllers.AddPhoto)
        photoRoutes.GET("/", middlewares.Authentication(), controllers.GetAllPhotos)
        photoRoutes.PUT("/:photoId", middlewares.Authentication(), controllers.UpdatePhoto)
        photoRoutes.DELETE("/:photoId", middlewares.Authentication(), controllers.DeletePhoto)
    }

    commentRoutes := r.Group("/comments")
    {
        commentRoutes.POST("/", middlewares.Authentication(), controllers.AddComment)
        commentRoutes.GET("/", middlewares.Authentication(), controllers.GetAllComments)
        commentRoutes.PUT("/:commentId", middlewares.Authentication(), controllers.UpdateComment)
        commentRoutes.DELETE("/:commentId", middlewares.Authentication(), controllers.DeleteComment)
    }

    socialMediaRoutes := r.Group("/socialmedias")
    {
        socialMediaRoutes.POST("/", middlewares.Authentication(), controllers.AddSocialMedia)
        socialMediaRoutes.GET("/", middlewares.Authentication(), controllers.GetAllSocialMedia)
        socialMediaRoutes.PUT("/:socialMediaId", middlewares.Authentication(), controllers.UpdateSocialMedia)
        socialMediaRoutes.DELETE("/:socialMediaId", middlewares.Authentication(), controllers.DeleteSocialMedia)
    }

    var PORT := os.Getenv("PORT")

    routers.StartServer().r.Run(":" + PORT)
}
