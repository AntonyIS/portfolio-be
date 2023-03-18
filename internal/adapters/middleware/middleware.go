package middleware

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
)

func GenerateToken(email string) (string, error) {
	tokenLifeSpan, err := strconv.Atoi(os.Getenv("TOKENLIFESPAN"))

	if err != nil {
		return "", fmt.Errorf("token lifespan error: %s", err.Error())
	}

	claims := jwt.MapClaims{}
	claims["authoization"] = true
	claims["email"] = email
	claims["exp"] = time.Now().Add(time.Minute * time.Duration(tokenLifeSpan)).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)

	return token.SignedString([]byte(os.Getenv("SECRET_KEY")))
}
