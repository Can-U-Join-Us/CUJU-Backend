package server

import (
	_ "github.com/Can-U-Join-Us/CUJU-Backend/modules/storage"
	"github.com/gin-gonic/gin"
)

const port = ":4000"

func init() { // local : 4000 호스팅 시작
	r := gin.Default()

	api := r.Group("/ui/api")
	registerApiHandlers(api)
	r.Run(port)
}
func registerApiHandlers(api *gin.RouterGroup) {
	api.GET("/ping", func(c *gin.Context) {
		message := pingTest(c)
		c.JSON(200, gin.H{
			"message": message,
		})
	})
	api.GET("/db", func(c *gin.Context) {
		err := dbConnectionTest(c)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
		} else {
			c.JSON(200, gin.H{"error": nil})
		}
	})
	api.POST("/User/login", func(c *gin.Context) {
		err := loginUser(c)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
		} else {
			c.JSON(200, gin.H{"error": nil})
		}
	})
	api.POST("/User/register", func(c *gin.Context) {
		err := registerUser(c)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
		} else {
			c.JSON(200, gin.H{"error": nil})
		}
	})
	api.GET("/Posts", func(c *gin.Context) {
		posts, err := getPostList(c)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
		} else {
			c.JSON(200, gin.H{"error": nil, "post": posts})
		}
	})
	api.POST("/Posts/add", func(c *gin.Context) {
		err := addPost(c)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
		} else {
			c.JSON(200, gin.H{"error": nil})
		}
	})
}
