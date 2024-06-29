package middlewares

import (
    "net/http"
    "uk-achmadxander/helpers"
    "github.com/dgrijalva/jwt-go"
    "github.com/gin-gonic/gin"
    "strings"
)

func Authentication() gin.HandlerFunc {
    return func(c *gin.Context) {
        tokenString := c.GetHeader("Authorization")
        if tokenString == "" {
            helpers.RespondError(c, http.StatusUnauthorized, "Unauthorized")
            c.Abort()
            return
        }

        tokenString = strings.Replace(tokenString, "Bearer ", "", 1) // Remove "Bearer " prefix if present

        claims := &jwt.StandardClaims{}
        token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
            return helpers.JwtKey, nil
        })

        if err != nil || !token.Valid {
            helpers.RespondError(c, http.StatusUnauthorized, "Unauthorized")
            c.Abort()
            return
        }

        // Set user_id in the context for controllers to access
		c.Set("user_id", claims.Subject)
        c.Next()
    }
}
