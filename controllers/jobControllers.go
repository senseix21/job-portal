package controllers

import (
	"job-portal/models"
	"job-portal/services"
	"job-portal/utils"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
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


// ListJobsHandler handles the GET request for fetching job listings with filters and search
func (jc *JobController) ListJobsHandler(c echo.Context) error {
	// Parse query parameters
	datePosted := c.QueryParam("datePosted")     // e.g., 'anytime' or a specific date range
	jobType := c.QueryParam("jobType")           // e.g., 'full-time', 'part-time', etc.
	salaryRange := c.QueryParam("salaryRange")   // e.g., '0-2.5k'
	workLocation := c.QueryParam("workLocation") // 'on-site', 'remote', 'hybrid'
	search := c.QueryParam("search")             // Search term (e.g., job title or description)

	// Pagination (default to page 1 and 10 items per page)
	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil || page <= 0 {
		page = 1
	}
	pageSize, err := strconv.Atoi(c.QueryParam("pageSize"))
	if err != nil || pageSize <= 0 {
		pageSize = 10
	}

	// Apply filters to construct the query filter
	// Here, we need to pass an additional bson.M for the existing filters
	filter := bson.M{} // Initialize the filter as an empty bson.M
	filter, err = services.ApplyFilters(datePosted, jobType, salaryRange, workLocation, filter)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  http.StatusInternalServerError,
			"message": "Failed to apply filters",
			"error":   err.Error(),
		})
	}

	// Fetch filtered jobs with pagination and search
	// Now passing all the necessary arguments to ListJobs
	jobs, pagination, err := jc.JobService.ListJobs(filter, page, pageSize, search, datePosted, jobType, salaryRange, workLocation)
	if err != nil {
		// Handle error and return a custom error response
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  http.StatusInternalServerError,
			"message": "Failed to retrieve jobs",
			"error":   err.Error(),
		})
	}

	// Prepare response data with pagination
	response := map[string]interface{}{
		"status":      http.StatusOK,
		"message":     "Jobs retrieved successfully",
		"totalItems":  pagination["totalItems"],
		"totalPages":  pagination["totalPages"],
		"currentPage": pagination["currentPage"],
		"pageSize":    pagination["pageSize"],
		"data": map[string]interface{}{
			"jobs": jobs,
		},
	}

	// Send response
	return c.JSON(http.StatusOK, response)
}
