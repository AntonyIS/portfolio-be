package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/AntonyIS/portfolio-be/internal/core/services"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

type middleware struct {
	svc *services.PortfolioService
}

func NewMiddleware(svc *services.PortfolioService) *middleware {
	return &middleware{
		svc: svc,
	}
}

func (m middleware) GenerateToken(id string) (string, error) {
	key := []byte(os.Getenv("SECRET_KEY"))
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	user, err := m.svc.ReadUser(id)
	if err != nil {
		return "", err
	}

	claims["user_id"] = user.Id
	claims["firstname"] = user.FirstName
	claims["lastname"] = user.LastName
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()

	tokenString, err := token.SignedString(key)

	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (m middleware) Authorize(c *gin.Context) {
	tokenString := c.GetHeader("token")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["sub"])
		}
		return []byte(os.Getenv("SECRET_KEY")), nil
	})

	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		email := fmt.Sprintf("%s", claims["email"])
		user_id := fmt.Sprintf("%s", claims["user_id"])
		firstname := fmt.Sprintf("%s", claims["firstname"])
		lastname := fmt.Sprintf("%s", claims["lastname"])
		c.Set("email", email)
		c.Set("user_id", user_id)
		c.Set("firstname", firstname)
		c.Set("lastname", lastname)
		c.Next()
	} else {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": errors.New("Request not authorized"),
		})
		return
	}
}
