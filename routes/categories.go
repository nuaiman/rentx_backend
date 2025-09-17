package routes

import (
	"fmt"
	"net/http"
	"rentx/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

func createCategory(c *gin.Context) {
	var category models.Category
	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
		return
	}
	userId := c.GetInt64("userId")
	category.UserId = userId
	if err := category.Save(); err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Could not create category"})
		return
	}
	c.JSON(http.StatusCreated, category)
}

func updateCategory(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	userId := c.GetInt64("userId")
	var category models.Category
	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
		return
	}
	category.Id = id
	category.UserId = userId
	if err := category.Update(userId); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, category)
}

func deleteCategory(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	userId := c.GetInt64("userId")
	category := models.Category{Id: id, UserId: userId}
	if err := category.Delete(userId); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Deleted"})
}

func listCategories(c *gin.Context) {
	categories, _ := models.GetCategories()
	c.JSON(http.StatusOK, categories)
}
