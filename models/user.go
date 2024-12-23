package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Role Constants
const (
	RoleAdmin     = "admin"
	RoleUser      = "user"
	RoleRecruiter = "recruiter"
)

// User struct with oneof validation for roles
type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name      string             `json:"name" validate:"required"`
	Email     string             `json:"email" validate:"required,email" bson:"email"`
	Password  string             `json:"password" validate:"required"`
	Role      string             `json:"role" validate:"required,oneof=admin user recruiter"` // Oneof validation for roles
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}
