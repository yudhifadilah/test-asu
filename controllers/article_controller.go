package controllers

import (
	"encoding/json"
	"fmt"
	"influencer-golang/config"
	"influencer-golang/models"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"golang.org/x/net/context"
)

var ctx = context.Background()

// GetAllArticles mengambil semua artikel dari Redis
func GetAllArticles(c *gin.Context) {
	keys, err := config.RedisClient.Keys(ctx, "article:*").Result()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve articles"})
		return
	}

	var articles []models.Article
	for _, key := range keys {
		articleJSON, err := config.RedisClient.Get(ctx, key).Result()
		if err != nil {
			continue
		}

		var article models.Article
		if err := json.Unmarshal([]byte(articleJSON), &article); err == nil {
			articles = append(articles, article)
		}
	}

	c.JSON(http.StatusOK, gin.H{"articles": articles})
}

// GetArticleByID mengambil artikel berdasarkan ID dari Redis
func GetArticleByID(c *gin.Context) {
	id := c.Param("id")
	key := fmt.Sprintf("article:%s", id)

	articleJSON, err := config.RedisClient.Get(ctx, key).Result()
	if err == redis.Nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Article not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve article"})
		return
	}

	var article models.Article
	if err := json.Unmarshal([]byte(articleJSON), &article); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse article data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"article": article})
}

// CreateArticle membuat artikel baru dengan menyimpan ke Redis
func CreateArticle(c *gin.Context) {
	title := c.PostForm("title")
	excerpt := c.PostForm("excerpt")
	content := c.PostForm("content")

	if title == "" || excerpt == "" || content == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "All fields are required"})
		return
	}

	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Image file is required"})
		return
	}

	imageDir := "uploads/image"
	if _, err := os.Stat(imageDir); os.IsNotExist(err) {
		os.MkdirAll(imageDir, os.ModePerm)
	}

	imagePath := filepath.Join(imageDir, file.Filename)
	if err := c.SaveUploadedFile(file, imagePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save image"})
		return
	}

	// Generate ID unik dari Redis
	id := fmt.Sprintf("%d", config.RedisClient.Incr(ctx, "article_id").Val())

	article := models.Article{
		ID:      id,
		Title:   title,
		Excerpt: excerpt,
		Content: content,
		Image:   imagePath,
	}

	articleJSON, err := json.Marshal(article)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encode article"})
		return
	}

	if err := config.RedisClient.Set(ctx, fmt.Sprintf("article:%s", id), articleJSON, 0).Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save article"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Article created successfully", "article": article})
}

// UpdateArticle memperbarui artikel yang ada
func UpdateArticle(c *gin.Context) {
	id := c.Param("id")
	key := fmt.Sprintf("article:%s", id)

	articleJSON, err := config.RedisClient.Get(ctx, key).Result()
	if err == redis.Nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Article not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve article"})
		return
	}

	var article models.Article
	if err := json.Unmarshal([]byte(articleJSON), &article); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse article data"})
		return
	}

	title := c.PostForm("title")
	excerpt := c.PostForm("excerpt")
	content := c.PostForm("content")

	if title != "" {
		article.Title = title
	}
	if excerpt != "" {
		article.Excerpt = excerpt
	}
	if content != "" {
		article.Content = content
	}

	// Jika ada file gambar baru yang diunggah
	file, err := c.FormFile("image")
	if err == nil {
		// Hapus gambar lama jika ada
		if article.Image != "" {
			os.Remove(article.Image)
		}

		// Simpan gambar baru
		imagePath := filepath.Join("uploads/image", file.Filename)
		if err := c.SaveUploadedFile(file, imagePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save new image"})
			return
		}
		article.Image = imagePath
	}

	// Simpan perubahan ke Redis
	updatedArticleJSON, err := json.Marshal(article)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encode updated article"})
		return
	}

	if err := config.RedisClient.Set(ctx, key, updatedArticleJSON, 0).Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update article"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Article updated successfully", "article": article})
}

// DeleteArticle menghapus artikel dari Redis
func DeleteArticle(c *gin.Context) {
	id := c.Param("id")
	key := fmt.Sprintf("article:%s", id)

	articleJSON, err := config.RedisClient.Get(ctx, key).Result()
	if err == redis.Nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Article not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve article"})
		return
	}

	var article models.Article
	if err := json.Unmarshal([]byte(articleJSON), &article); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse article data"})
		return
	}

	// Hapus gambar jika ada
	if article.Image != "" {
		os.Remove(article.Image)
	}

	// Hapus artikel dari Redis
	if err := config.RedisClient.Del(ctx, key).Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete article"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Article deleted successfully"})
}
