package helpers

import (
	"errors"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

const secretKey = "projectMyGramHeruPurnama"

func GenerateJwtToken(id uint, email string, username string, age uint, pp string, created_at time.Time, updated_at time.Time) string {
	claims := jwt.MapClaims{
		"id":                id,
		"email":             email,
		"username":          username,
		"age":               age,
		"profile_image_url": pp,
		"created_at":        created_at,
		"updated_at":        updated_at,
	}
	parseToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, _ := parseToken.SignedString([]byte(secretKey))
	return signedToken
}

func VerifyToken(c *gin.Context) (interface{}, error) {
	err := errors.New("Sign in needed")
	auth := c.Request.Header.Get("Authorization")
	tokenBearer := strings.HasPrefix(auth, "Bearer")
	if !tokenBearer {
		return nil, err
	}
	stringToken := strings.Split(auth, " ")[1]

	token, err := jwt.Parse(stringToken, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, err
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, err
	}

	return claims, nil
}
