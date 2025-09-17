package routes

import (
	"net/http"
	"rentx/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Create a new review
func createReview(c *gin.Context) {
	var r models.Review
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
		return
	}

	r.UserId = c.GetInt64("userId") // from middleware

	if err := r.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, r)
}

// List all reviews for a post
func listReviewsByPost(c *gin.Context) {
	postId, _ := strconv.ParseInt(c.Param("postId"), 10, 64)

	reviews, err := models.ListReviews(postId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Could not fetch reviews"})
		return
	}

	c.JSON(http.StatusOK, reviews)
}

// Delete a review (owner or admin only)
func deleteReview(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	userId := c.GetInt64("userId")

	r := models.Review{Id: id}
	if err := r.Delete(userId); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Review deleted"})
}
