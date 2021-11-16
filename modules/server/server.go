package server

import (
	"fmt"

	_ "github.com/Can-U-Join-Us/CUJU-Backend/modules/storage"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

const port = ":4000"

var ACCESS_SECRET string
var REFRESH_SECRET string

func init() { // local : 4000 호스팅 시작
	r := gin.Default()
	redisInit()
	api := r.Group("/api")
	registerApiHandlers(api)
	viper.SetConfigName("config")
	viper.AddConfigPath(".")    // optionally look for config in the working directory
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %w \n", err))
	}
	ACCESS_SECRET = viper.GetString(`token.ACCESS_SECRET`)
	REFRESH_SECRET = viper.GetString(`token.REFRESH_SECRET`)

	fmt.Println("\n\n", ACCESS_SECRET, REFRESH_SECRET, "\n\n")
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
	/*  Reply			200 -> token
	400 -> ID or PW incorrect
	*/
	api.POST("/User/login", func(c *gin.Context) {
		token, err := loginUser(c)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
		} else {
			c.JSON(200, gin.H{"error": nil, "token": token})
		}
	})
	/*  Reply			200 -> token delete
	400 -> ID or PW incorrect
	*/
	api.POST("/User/logout", func(c *gin.Context) {
		err := logoutUser(c)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
		} else {
			c.JSON(200, gin.H{"error": nil})
		}
	})
	/*  Reply			200 -> register success
	400 -> DB Conn or Query err
	*/
	api.POST("/User/register", func(c *gin.Context) {
		err := registerUser(c)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
		} else {
			c.JSON(200, gin.H{"error": nil})
		}
	})
	/*  Reply			200 -> Get List<post> success
	400 -> DB Conn or Query err
	*/
	api.GET("/Posts", func(c *gin.Context) {
		posts, err := getPostList(c)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
		} else {
			c.JSON(200, gin.H{"error": nil, "post": posts})
		}
	})
	/*  Reply			200 -> Add post success
	400 -> DB Conn or Query err
	*/
	api.POST("/Posts/add", func(c *gin.Context) {
		err := addPost(c)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
		} else {
			c.JSON(200, gin.H{"error": nil})
		}
	})
}
