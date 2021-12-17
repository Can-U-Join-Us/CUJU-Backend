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
	if err := redisInit(); err != nil {
		panic(fmt.Errorf("Fatal error : redis is off status \n"))
	}
	api := r.Group("/api")
	api.Use(dummy)
	registerApiHandlers(api)
	viper.SetConfigName("config")
	viper.AddConfigPath(".")    // optionally look for config in the working directory
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %w \n", err))
	}
	ACCESS_SECRET = viper.GetString(`token.ACCESS_SECRET`)
	REFRESH_SECRET = viper.GetString(`token.REFRESH_SECRET`)

	r.Run(port)
}
func dummy(c *gin.Context) {
	fmt.Println("Access Token Check Stage")
}
func registerApiHandlers(api *gin.RouterGroup) {
	/*  Reply			200 -> token
	400 -> ID or PW incorrect
	*/
	api.POST("/User/login", func(c *gin.Context) {
		uid, token, err := loginUser(c)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
		} else {
			c.JSON(200, gin.H{"error": nil, "token": token, "uid": uid})
		}
	})
	/*  Reply			200 -> null
	400 -> Modify fail
	*/
	api.POST("/User/modify/pw", func(c *gin.Context) {
		err := modifyPW(c)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
		} else {
			c.JSON(200, gin.H{"error": nil})
		}
	})
	/*  Reply			200 -> null
	400 -> Modify fail
	*/
	api.POST("/User/modify/profile", func(c *gin.Context) {
		err := modifyProfile(c)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
		} else {
			c.JSON(200, gin.H{"error": nil})
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
	/*  Reply			200 -> find pw success
	400 -> DB Conn or Query err
	*/
	api.POST("/User/find/pw", func(c *gin.Context) {
		err := findUserPW(c)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
		} else {
			c.JSON(200, gin.H{"error": nil})
		}
	})
	/*  Reply			200 -> find id success
	400 -> DB Conn or Query err
	*/
	api.POST("/User/find/id", func(c *gin.Context) {
		id, err := findUserId(c)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
		} else {
			c.JSON(200, gin.H{"error": nil, "id": id})
		}
	})
	/*  Reply			200 -> Get List<post> success
	400 -> DB Conn or Query err
	*/
	api.GET("/Projects", func(c *gin.Context) {
		posts, err := getProjectList(c)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
		} else {
			c.JSON(200, gin.H{"error": nil, "projects": posts})
		}
	})
	/*  Reply			200 -> Get project obj success
	400 -> DB Conn or Query err
	*/
	api.GET("/Project/", func(c *gin.Context) {
		pd, err := getProjectDetail(c)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
		} else {
			c.JSON(200, gin.H{"error": nil, "projects": pd})
		}
	})
	/*  Reply			200 -> Get List<post> success
	400 -> DB Conn or Query err
	*/
	api.GET("/Projects/category", func(c *gin.Context) {
		posts, err := getCategory(c)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
		} else {
			c.JSON(200, gin.H{"error": nil, "projects": posts})
		}
	})
	/*  Reply			200 -> Add post success
	400 -> DB Conn or Query err
	*/
	api.POST("/Projects/add", func(c *gin.Context) {
		pid, err := addProject(c)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
		} else {
			c.JSON(200, gin.H{"error": nil, "pid": pid})
		}
	})
	/*  Reply			200 -> Permit post success
	400 -> DB Conn or Query err
	*/
	api.POST("/Projects/permit", func(c *gin.Context) {
		err := permitProject(c)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
		} else {
			c.JSON(200, gin.H{"error": nil})
		}
	})
	/*  Reply			200 -> Deny post success
	400 -> DB Conn or Query err
	*/
	api.POST("/Projects/deny", func(c *gin.Context) {
		err := denyProject(c)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
		} else {
			c.JSON(200, gin.H{"error": nil})
		}
	})
	/*  Reply			200 -> Join post success
	400 -> DB Conn or Query err
	*/
	api.POST("/Projects/join", func(c *gin.Context) {
		err := joinProject(c)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
		} else {
			c.JSON(200, gin.H{"error": nil})
		}
	})
	/*  Reply			200 -> Get List<msg> success
	400 -> DB Conn or Query err
	*/
	api.GET("/Refresh", func(c *gin.Context) {
		msgList, err := refreshMsg(c)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
		} else {
			c.JSON(200, gin.H{"error": nil, "Message": msgList})
		}
	})
	/*  Reply			200 -> POST announceMSG success
	400 -> DB Conn or Query err
	*/
	api.POST("/announce", func(c *gin.Context) {
		err := postAnnouncement(c)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
		} else {
			c.JSON(200, gin.H{"error": nil})
		}
	})
}
