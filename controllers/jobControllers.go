package controllers

import (
	"job-portal/models"
	"job-portal/services"
	"job-portal/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

type JobController struct {
	JobService *services.JobService
}

func NewJobController(jobService *services.JobService) *JobController {
	return &JobController{JobService: jobService}
}

// CreateJobHandler handles job creation
func (jc *JobController) CreateJobHandler(c echo.Context) error {
	var job models.Job
	if err := c.Bind(&job); err != nil {
		return err // Pass to the custom error handler
	}

	if err := c.Validate(&job); err != nil {
		return err // Validation errors are handled by the custom error handler
	}

	if err := jc.JobService.CreateJob(&job); err != nil {
		return err // Pass business logic errors to the error handler
	}

	return utils.SendResponse(c, http.StatusCreated, "Job created successfully", job)
}

// GetJobHandler retrieves a job by ID
func (jc *JobController) GetJobHandler(c echo.Context) error {
	id := c.Param("id")
	job, err := jc.JobService.GetJob(id)
	if err != nil {
		return err // Pass errors to the custom error handler
	}

	return utils.SendResponse(c, http.StatusOK, "Job retrieved successfully", job)
}

// UpdateJobHandler updates a job by ID
func (jc *JobController) UpdateJobHandler(c echo.Context) error {
	id := c.Param("id")
	var updateData map[string]interface{}

	if err := c.Bind(&updateData); err != nil {
		return err // Pass binding errors to the custom error handler
	}

	updatedJob, err := jc.JobService.UpdateJob(id, updateData)
	if err != nil {
		return err // Pass service errors to the custom error handler
	}

	return utils.SendResponse(c, http.StatusOK, "Job updated successfully", updatedJob)
}

// DeleteJobHandler deletes a job by ID
func (jc *JobController) DeleteJobHandler(c echo.Context) error {
	id := c.Param("id")

	deletedJob, err := jc.JobService.DeleteJob(id)
	if err != nil {
		return err // Pass errors to the custom error handler
	}

	return utils.SendResponse(c, http.StatusOK, "Job deleted successfully", deletedJob)
}

// ListJobsHandler retrieves all jobs with optional filters
func (jc *JobController) ListJobsHandler(c echo.Context) error {
	filter := make(map[string]interface{}) // Add logic to parse query params into filters if needed
	jobs, err := jc.JobService.ListJobs(filter)
	if err != nil {
		return err // Pass errors to the custom error handler
	}

	return utils.SendResponse(c, http.StatusOK, "Jobs retrieved successfully", jobs)
}
