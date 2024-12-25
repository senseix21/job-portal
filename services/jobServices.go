package services

import (
	"context"
	"errors"
	"fmt"
	"job-portal/models"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type JobService struct {
	Collection *mongo.Collection
}

// NewJobService creates a new instance of JobService
func NewJobService(collection *mongo.Collection) *JobService {
	return &JobService{Collection: collection}
}

// CreateJob adds a new job to the database
func (s *JobService) CreateJob(job *models.Job) error {
	job.ID = primitive.NewObjectID()
	job.CreatedAt = time.Now()
	job.UpdatedAt = job.CreatedAt
	job.PostedAt = time.Now() // Assume posted immediately
	_, err := s.Collection.InsertOne(context.TODO(), job)
	return err
}

// GetJob retrieves a job by its ID
func (s *JobService) GetJob(id string) (*models.Job, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid job ID format")
	}

	var job models.Job
	err = s.Collection.FindOne(context.TODO(), bson.M{"_id": objID}).Decode(&job)
	if err != nil {
		return nil, errors.New("job not found")
	}

	return &job, nil
}



// UpdateJob updates an existing job and returns the updated job
func (s *JobService) UpdateJob(id string, updateData map[string]interface{}) (*models.Job, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid job ID format")
	}

	// Retrieve the current job before updating
	var job models.Job
	err = s.Collection.FindOne(context.TODO(), bson.M{"_id": objID}).Decode(&job)
	if err != nil {
		return nil, errors.New("job not found")
	}

	// Update the job
	update := bson.M{
		"$set": updateData,
		"$currentDate": bson.M{
			"updated_at": true,
		},
	}
	_, err = s.Collection.UpdateOne(context.TODO(), bson.M{"_id": objID}, update)
	if err != nil {
		return nil, err
	}

	// Retrieve the updated job
	err = s.Collection.FindOne(context.TODO(), bson.M{"_id": objID}).Decode(&job)
	if err != nil {
		return nil, errors.New("failed to retrieve updated job")
	}

	return &job, nil
}

// DeleteJob removes a job from the database and returns the deleted job
func (s *JobService) DeleteJob(id string) (*models.Job, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid job ID format")
	}

	// Retrieve the job before deleting
	var job models.Job
	err = s.Collection.FindOne(context.TODO(), bson.M{"_id": objID}).Decode(&job)
	if err != nil {
		return nil, errors.New("job not found")
	}

	// Delete the job
	_, err = s.Collection.DeleteOne(context.TODO(), bson.M{"_id": objID})
	if err != nil {
		return nil, err
	}

	return &job, nil
}

// ListJobs retrieves jobs based on filters, search parameters, and pagination
func (s *JobService) ListJobs(filter bson.M, page, pageSize int, search, datePosted, jobType, salaryRange, workLocation string) ([]models.Job, map[string]interface{}, error) {
	// Apply the filters based on the provided parameters
	var err error
	filter, err = ApplyFilters(datePosted, jobType, salaryRange, workLocation, filter)
	if err != nil {
		return nil, nil, err
	}

	// Apply search filter if provided
	if search != "" {
		// Search in job title and description
		searchRegex := bson.M{
			"$regex":   search,
			"$options": "i", // Case-insensitive search
		}

		// Combine search condition with other filters using $or
		if existingFilter, ok := filter["$or"]; ok {
			// If there's already an existing $or filter, append the search conditions
			filter["$or"] = append(existingFilter.([]bson.M),
				bson.M{"title": searchRegex},
				bson.M{"description": searchRegex},
			)
		} else {
			// If no $or filter exists, create a new $or filter for search
			filter["$or"] = []bson.M{
				{"title": searchRegex},
				{"description": searchRegex},
			}
		}
	}

	// Set pagination options
	options := options.Find()
	options.SetSkip(int64((page - 1) * pageSize)) // Skip for pagination
	options.SetLimit(int64(pageSize))             // Limit results per page

	// Fetch jobs from the collection
	cursor, err := s.Collection.Find(context.TODO(), filter, options)
	if err != nil {
		return nil, nil, err
	}
	defer cursor.Close(context.TODO())

	var jobs []models.Job
	if err := cursor.All(context.TODO(), &jobs); err != nil {
		return nil, nil, err
	}

	// Count total items matching the filter
	totalItems, err := s.Collection.CountDocuments(context.TODO(), filter)
	if err != nil {
		return nil, nil, err
	}

	// Calculate total pages
	totalPages := int(math.Ceil(float64(totalItems) / float64(pageSize)))

	// Prepare pagination data
	pagination := map[string]interface{}{
		"totalItems":  totalItems,
		"totalPages":  totalPages,
		"currentPage": page,
		"pageSize":    pageSize,
	}

	return jobs, pagination, nil
}

// ApplyFilters constructs the filter based on the provided parameters
func ApplyFilters(datePosted, jobType, salaryRange, workLocation string, existingFilter bson.M) (bson.M, error) {
	// Initialize the filter if it doesn't exist
	if existingFilter == nil {
		existingFilter = bson.M{}
	}

	// Date filter (e.g., last 30 days)
	// Date filter (e.g., last 30 days, last 7 days, last 24 hours)
	if datePosted == "anytime" {
	// No filter
	} else if datePosted != "" {
	// Handle date range logic (last 7 days, last 30 days, last 24 hours, etc.)
	var postedDate time.Time
	switch datePosted {
	case "last_7_days":
		postedDate = time.Now().AddDate(0, 0, -7)
	case "last_24_hours":
		postedDate = time.Now().Add(-24 * time.Hour) // Subtract 24 hours from the current time
	default:
		postedDate = time.Now().AddDate(0, 0, -30) // Default to last 30 days
	}
	existingFilter["posted_at"] = bson.M{"$gte": postedDate}
}


	// Job Type filter (full-time, part-time, etc.)
	if jobType != "" {
		existingFilter["type"] = jobType
	}

		// Salary Range filter (e.g., "$50,000 - $90,000")
	if salaryRange != "" {
    // Remove currency symbols, commas, and spaces to extract numeric values
    re := regexp.MustCompile(`[^\d.-]`) // Match everything except digits, period, and minus sign
    cleanedSalaryRange := re.ReplaceAllString(salaryRange, "")

    // Split the salary range by "-"
    salaryRangeSplit := strings.Split(cleanedSalaryRange, "-")
    if len(salaryRangeSplit) == 2 {
        // Convert the salary values to float64
        minSalary, err := strconv.ParseFloat(salaryRangeSplit[0], 64)
        if err != nil {
            return nil, fmt.Errorf("invalid min salary: %v", err)
        }
        maxSalary, err := strconv.ParseFloat(salaryRangeSplit[1], 64)
        if err != nil {
            return nil, fmt.Errorf("invalid max salary: %v", err)
        }

        // Add the salary filter to the existing filter
        existingFilter["min_salary"] = bson.M{
            "$gte": minSalary,
        }
        existingFilter["max_salary"] = bson.M{
            "$lte": maxSalary,
        }
    } else {
        return nil, fmt.Errorf("invalid salary range format: %s", salaryRange)
    }
}




	// Work Location filter (on-site, remote, hybrid)
	if workLocation != "" {
		existingFilter["work_location"] = workLocation
	}

	// Return the updated filter
	return existingFilter, nil
}
