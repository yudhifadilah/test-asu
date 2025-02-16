package routes

import (
	"influencer-golang/controllers"

	"github.com/gin-gonic/gin"
)

// SetupRoutes mengatur semua rute aplikasi
func SetupRoutes(r *gin.Engine) {

	// Article Endpoint
	articleGroup := r.Group("/api/articles")
	{
		articleGroup.GET("/", controllers.GetAllArticles)
		articleGroup.GET("/:id", controllers.GetArticleByID)
		articleGroup.POST("/", controllers.CreateArticle)
		articleGroup.PUT("/:id", controllers.UpdateArticle)
		articleGroup.DELETE("/:id", controllers.DeleteArticle)
	}

}
