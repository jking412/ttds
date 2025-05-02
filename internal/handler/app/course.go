package app

import (
	"awesomeProject/internal/usecase"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

var courseService usecase.CourseService

// initCourseService initializes the course service.
// It's called from InitApp in init.go
func initCourseService() {
	// Assuming NewCourseService exists and initializes the service
	// Make sure NewCourseService is implemented in the usecase package
	courseService = usecase.NewCourseService()
}

// GetAllCoursesHandler handles the request to get all courses.
// GET /courses
// Response format based on API spec:
//
//	{
//	  "status": "success",
//	  "data": [
//	    {
//	      "course_id": 1,
//	      "title": "Course Title",
//	      "description": "Course Description",
//	      "created_at": "2023-01-01T00:00:00Z"
//	    }
//	  ]
//	}
func GetAllCoursesHandler(c *gin.Context) {
	courses, err := courseService.GetAllCourses()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to retrieve courses: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   courses,
	})
}

// GetCourseByIDHandler handles the request to get a specific course by its ID.
// GET /courses/{course_id}
// Response format based on API spec:
//
//	{
//	  "status": "success",
//	  "data": {
//	    "course_id": 1,
//	    "title": "Course Title",
//	    "description": "Course Description",
//	    "created_at": "2023-01-01T00:00:00Z",
//	    "chapters": [
//	      {
//	        "chapter_id": 1,
//	        "title": "Chapter Title",
//	        "sections": [
//	          {
//	            "section_id": 1,
//	            "title": "Section Title"
//	          }
//	        ]
//	      }
//	    ]
//	  }
//	}
func GetCourseByIDHandler(c *gin.Context) {
	courseIDStr := c.Param("course_id")
	courseID, err := strconv.ParseUint(courseIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid course ID format"})
		return
	}

	course, err := courseService.GetCourseByID(uint(courseID))
	if err != nil {
		// Use a more robust error check, e.g., errors.Is(err, gorm.ErrRecordNotFound)
		if err.Error() == "record not found" { // Basic check, improve if possible
			c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "Course not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to retrieve course: " + err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   course,
	})
}

// GetCourseReferencesHandler handles the request to get reference books for a specific course.
// GET /courses/{course_id}/reference-books
// Response format based on API spec:
//
//	{
//	  "status": "success",
//	  "data": [
//	    {
//	      "reference_book_id": 1,
//	      "title": "Book Title",
//	      "type": "PDF",
//	      "description": "Book Description",
//	      "url": "http://example.com/book.pdf"
//	    }
//	  ]
//	}
func GetCourseReferencesHandler(c *gin.Context) {
	courseIDStr := c.Param("course_id")
	courseID, err := strconv.ParseUint(courseIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid course ID format"})
		return
	}

	referenceBooks, err := courseService.GetCourseReferences(uint(courseID))
	if err != nil {
		// Use a more robust error check
		if err.Error() == "record not found" { // Basic check, improve if possible
			c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "Reference books not found for this course"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to retrieve reference books: " + err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   referenceBooks,
	})
}

// GetCourseExperimentStatusHandler handles the request to get user's course experiment status.
// GET /api/v1/user/course-experiment-status?course_id={course_id}
// Response format based on API spec:
//
//	{
//	  "status": "success",
//	  "data": [
//	    {
//	      "section_id": 1,
//	      "section_title": "Section Title",
//	      "completed": true
//	    }
//	  ]
//	}
func GetCourseExperimentStatusHandler(c *gin.Context) {
	// Get user ID from context (assuming it's set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "Unauthorized"})
		return
	}

	courseIDStr := c.Query("course_id")
	courseID, err := strconv.ParseUint(courseIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid course ID format"})
		return
	}

	status, err := courseService.GetCourseStatus(uint(userID.(uint)), uint(courseID))
	if err != nil {
		if err.Error() == "record not found" {
			c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "Course experiment status not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to retrieve course experiment status: " + err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   status,
	})
}
