package helpers

import (
    "time"
    "strconv"
    "github.com/dgrijalva/jwt-go"
)

var JwtKey = []byte("rahasia") // Make sure this matches the key used for signing tokens

func GenerateToken(userID uint, email string) (string, error) {
    expirationTime := time.Now().Add(24 * time.Hour)
    claims := &jwt.StandardClaims{
        ExpiresAt: expirationTime.Unix(),
        IssuedAt:  time.Now().Unix(),
        Subject:   strconv.FormatUint(uint64(userID), 10),
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenString, err := token.SignedString(JwtKey)
    if err != nil {
        return "", err
    }

    return tokenString, nil
}
