package services

import (
	"context"
	"errors"
	"job-portal/models"
	"job-portal/utils"

	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	Collection *mongo.Collection
}

func NewUserService(collection *mongo.Collection) *UserService {
	return &UserService{Collection: collection}
}

func (s *UserService) Register(user models.User) error {
	// Check if email already exists
	var existingUser models.User
	err := s.Collection.FindOne(context.TODO(), bson.M{"email": user.Email}).Decode(&existingUser)
	if err == nil {
		return errors.New("email already in use")
	}
	if user.Role == "" {
		user.Role = "user"
	}
	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)

	// Set timestamps
	user.ID = primitive.NewObjectID()
	user.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	user.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())

	// Insert user into database
	_, err = s.Collection.InsertOne(context.TODO(), user)
	return err
}

// Authenticate authenticates a user by email and password, and generates a JWT token
func (s *UserService) Authenticate(email, password string) (string, *models.User, error) {
	var user models.User

	// Find user by email
	err := s.Collection.FindOne(context.TODO(), bson.M{"email": email}).Decode(&user)
	if err != nil {
		return "", nil, errors.New(err.Error())
	}

	// Compare password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", nil, errors.New(err.Error())
	}

	// Generate JWT token
	token, err := utils.GenerateJWT(user)
	if err != nil {
		return "", nil, errors.New("failed to generate token")
	}

	return token, &user, nil
}
