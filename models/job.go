package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Job struct {
	ID               primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Title            string             `json:"title" validate:"required"`
	Description      string             `json:"description" validate:"required"`
	Location         string             `json:"location" validate:"required"`
	SalaryRange      string             `json:"salary_range" validate:"required"`
	Type             string             `json:"type" validate:"required,oneof=full-time part-time contract"`
	Experience       string             `json:"experience" validate:"required,oneof=entry mid senior"`
	Education        string             `json:"education" validate:"required,oneof=bachelor master phd"`
	Skills           []string           `json:"skills" validate:"required"`           // List of required skills
	Responsibilities []string           `json:"responsibilities" validate:"required"` // List of job responsibilities
	Benefits         []string           `json:"benefits" validate:"required"`         // List of benefits provided
	ApplyLink        string             `json:"apply_link" validate:"required,url"`   // URL where users can apply
	PostedAt         time.Time          `json:"posted_at" bson:"posted_at"`
	ApplyBy          time.Time          `json:"apply_by" bson:"apply_by"`
	CreatedAt        time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt        time.Time          `json:"updated_at" bson:"updated_at"`
}

// Enums for Job Fields
const (
	FullTime   = "full-time"
	PartTime   = "part-time"
	Contract   = "contract"
	EntryLevel = "entry"
	MidLevel   = "mid"
	Senior     = "senior"
	Bachelor   = "bachelor"
	Master     = "master"
	PhD        = "phd"
)
