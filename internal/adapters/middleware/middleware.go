package middleware

import (
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

func (m middleware) GenerateToken(email string) (string, error) {
	key := []byte(os.Getenv("SECRET_KEY"))
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	user, err := m.svc.ReadUserWithEmail(email)
	if err != nil {
		return "", err
	}

	claims["email"] = email
	claims["user_id"] = user.Id
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()

	tokenString, err := token.SignedString(key)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (m middleware) Authorize(ctx *gin.Context) {
	tokenString := ctx.GetHeader("token")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["sub"])
		}
		return []byte(os.Getenv("SECRET_KEY")), nil
	})
	if err != nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		email := fmt.Sprintf("%s", claims["email"])
		user_id := fmt.Sprintf("%s", claims["user_id"])
		ctx.Set("email", email)
		ctx.Set("user_id", user_id)
		ctx.Next()
	} else {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
}
