package services

import (
	"contact/db"
	"contact/models"
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var usersCollection *mongo.Collection = db.GetCollection("users")

var (
	ErrEmailTaken         = errors.New("username already taken")
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func getUserByEmail(email string) (*models.User, error) {
	var user models.User

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"email": email}
	err := usersCollection.FindOne(ctx, filter).Decode(&user)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return &user, nil

}

func RegisterUser(user *models.User) error {

	existingUser, err := getUserByEmail(user.Email)
	if existingUser != nil {
		return ErrEmailTaken
	}

	if err != nil {
		return err
	}
	//hash password
	password, err := HashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = string(password)

	if err != nil {
		return err
	}

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	_, err = usersCollection.InsertOne(ctx, user)
	return err

}

func LoginUser(email, password string) (*models.User, error) {
	user, err := getUserByEmail(email)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, ErrUserNotFound
	}

	// Compare the passwords
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	return user, nil
}
