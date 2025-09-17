package routes

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"rentx/models"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func uploadPostImages(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to read form"})
		return
	}

	files := form.File["files"] // expecting input name="files"
	if len(files) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "No files uploaded"})
		return
	}

	saveDir := "storage/posts"
	if err := os.MkdirAll(saveDir, os.ModePerm); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Could not create storage directory"})
		return
	}

	var urls []string
	for _, file := range files {
		ext := filepath.Ext(file.Filename)
		newFileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
		savePath := filepath.Join(saveDir, newFileName)
		if err := c.SaveUploadedFile(file, savePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Could not save file"})
			return
		}
		urls = append(urls, "/storage/posts/"+newFileName)
	}

	c.JSON(http.StatusOK, gin.H{
		"urls": urls,
	})
}

// ----------------- CREATE POST -----------------
func createPost(c *gin.Context) {
	var p models.Post
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input", "error": err.Error()})
		return
	}

	p.UserId = c.GetInt64("userId") // from auth middleware

	if err := p.Save(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Could not save post", "error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, p)
}

// ----------------- UPDATE POST -----------------
func updatePost(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid post ID"})
		return
	}

	var p models.Post
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input", "error": err.Error()})
		return
	}

	p.Id = id
	userId := c.GetInt64("userId")
	userRole := c.GetString("userRole")

	if err := p.Update(userId, userRole); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, p)
}

// ----------------- DELETE POST -----------------
func deletePost(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid post ID"})
		return
	}

	p := models.Post{Id: id}
	userId := c.GetInt64("userId")
	userRole := c.GetString("userRole")

	if err := p.Delete(userId, userRole); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Post deleted"})
}

// ----------------- GET POST BY ID -----------------
func getPostByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid post ID"})
		return
	}

	post, err := models.GetPostByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, post)
}

// ----------------- LIST ALL POSTS -----------------
func listPosts(c *gin.Context) {
	posts, err := models.ListPosts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Could not fetch posts", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, posts)
}
