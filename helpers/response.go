package helpers

import "github.com/gin-gonic/gin"

const (
    AppJSON = "application/json"
)

func GetContentType(c *gin.Context) string {
    return c.Request.Header.Get("Content-Type")
}

func RespondJSON(c *gin.Context, status int, payload interface{}) {
    c.JSON(status, payload)
}

func RespondError(c *gin.Context, status int, message string) {
    c.JSON(status, gin.H{"error": message})
}
