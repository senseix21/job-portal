package services

import (
	"context"
	"errors"
	"job-portal/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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

// ListJobs retrieves all jobs with optional filters
func (s *JobService) ListJobs(filter bson.M) ([]models.Job, error) {
	cursor, err := s.Collection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}

	var jobs []models.Job
	if err = cursor.All(context.TODO(), &jobs); err != nil {
		return nil, err
	}

	return jobs, nil

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
