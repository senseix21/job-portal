package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Job represents a job posting
type Job struct {
	ID               primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Title            string             `json:"title" validate:"required"`
	Description      string             `json:"description" validate:"required"`
	Location         string             `json:"location" validate:"required"`
	MinSalary        float64            `json:"min_salary" bson:"min_salary" validate:"required"` // Minimum salary (numeric)
	MaxSalary        float64            `json:"max_salary" bson:"max_salary" validate:"required"` // Maximum salary (numeric)
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

	// Company Info (nested object)
	CompanyName      string `json:"company_name" bson:"company_name"`
	CompanyLogo      string `json:"company_logo" bson:"company_logo"`
	CompanyDescription string `json:"company_description" bson:"company_description"`

	// Work Location (On-site, Remote, Hybrid)
	WorkLocation     string `json:"work_location" bson:"work_location" validate:"required,oneof=on-site remote hybrid"`
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
	OnSite     = "on-site"
	Remote     = "remote"
	Hybrid     = "hybrid"
)
