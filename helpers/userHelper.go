package helpers

import (
	"awesomeProject2/model/userModel"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"os"
	"strconv"
	"time"
)

var jwtKey = []byte(os.Getenv("JWT_KEY"))

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		panic(err)
	}
	return string(bytes)
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GenerateJwt(jwtPayload *userModel.JWTPayload) (tokenString string, err error) {
	expirationDate := time.Now().Add(1 * time.Hour)
	subject := uuid.New()
	claims := &userModel.JWTClaim{
		Username:   jwtPayload.Username,
		Email:      jwtPayload.Email,
		Id:         jwtPayload.Id,
		CustomData: map[string]interface{}{"subject": subject.String()},
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationDate.Unix(),
			Audience:  os.Getenv("FRONTEND_URL"),
			Issuer:    os.Getenv("BACKEND_URL"),
			Subject:   subject.String(),
			Id:        strconv.Itoa(jwtPayload.Id),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err = token.SignedString(jwtKey)
	return
}

func ValidateToken(signedToken string) (*userModel.JWTClaim, error) {
	tokenString, err := jwt.ParseWithClaims(
		signedToken,
		&userModel.JWTClaim{},
		func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
	if err != nil {
		return nil, err
	}
	claims, ok := tokenString.Claims.(*userModel.JWTClaim)
	if !ok || tokenString.Valid == false {
		return nil, err
	}
	if claims.ExpiresAt < time.Now().Unix() {
		return nil, errors.New("token expired")
	}
	return claims, nil
}
