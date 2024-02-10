package token

import (
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SignedDetails struct {
	Name  string
	Id    primitive.ObjectID
	Email string
	jwt.StandardClaims
}

var secretKey = os.Getenv("SECRET_KEY")

func TokenGenerator(id primitive.ObjectID, name string, email string) (token string, refreshToken string, err error) {
	claims := SignedDetails{
		Email: email,
		Name:  name,
		Id:    id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
			Issuer:    "microservice",
			IssuedAt:  time.Now().Unix(),
		},
	}

	refreshClaims := SignedDetails{
		Email: email,
		Name:  name,
		Id:    id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 144).Unix(),
			Issuer:    "microservice",
			IssuedAt:  time.Now().Unix(),
		},
	}

	var tokenString string
	tokenString, err = jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secretKey))
	if err != nil {
		return "", "", err
	}

	var refreshTokenString string
	refreshTokenString, err = jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(secretKey))
	if err != nil {
		return "", "", err
	}

	return tokenString, refreshTokenString, err
}

func ValidateToken(signedToken string) (claims *SignedDetails, msg string) {
	token, err := jwt.ParseWithClaims(signedToken, &SignedDetails{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		msg = err.Error()
		return nil, msg
	}

	claims, ok := token.Claims.(*SignedDetails)

	if !ok {
		msg = "The Token is invalid"
		return
	}
	if claims.ExpiresAt < time.Now().Local().Unix() {
		msg = "token is expired"
		return
	}
	return claims, msg
}

// update token when refresh token is valid
func UpdateToken(refreshToken string) (token string, err error) {
	claims, msg := ValidateToken(refreshToken)
	if msg != "" {
		return "", err
	}
	if claims.ExpiresAt < time.Now().Unix() {
		claims.ExpiresAt = time.Now().Add(time.Hour * 72).Unix()
		token, err = jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secretKey))
	}
	return token, err
}
