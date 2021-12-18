package server

import (
	"fmt"
	"log"

	_ "github.com/Can-U-Join-Us/CUJU-Backend/modules/storage"
	"github.com/Can-U-Join-Us/CUJU-Backend/modules/token"
	"github.com/gin-gonic/gin"
)

const port = ":4000"

func Serve(mode int) { // local : 4000 호스팅 시작
	r := gin.Default()
	if err := token.RedisInit(); err != nil {
		panic(fmt.Errorf("Fatal error : redis is off status \n"))
	}
	api := r.Group("/api")
	if mode == 1 {
		api.Use(dummy)
	}
	RegisterApiHandlers(api)
	r.Run(port)
}

func dummy(c *gin.Context) {
	// access token of request header check stage
	log.Println("Access Token Check Stage")
}
func RegisterApiHandlers(api *gin.RouterGroup) {
	/*  Reply			200 -> token , uid
	400 -> ID or PW incorrect
	*/
	api.POST("/User/login", getLogin)

	/*  Reply			200 -> null
	400 -> Modify fail
	*/
	api.POST("/User/modify/pw", postModifyPW)

	/*  Reply			200 -> null
	400 -> Modify fail
	*/
	api.POST("/User/modify/profile", postModifyProfile)

	/*  Reply			200 -> token delete
	400 -> ID or PW incorrect
	*/
	api.POST("/User/logout", postLogout)

	/*  Reply			200 -> register success
	400 -> DB Conn or Query err
	*/
	api.POST("/User/register", postRegister)

	/*  Reply			200 -> null ( mail send )
	400 -> DB Conn or Query err
	*/
	api.POST("/User/find/pw", postFindPW)

	/*  Reply			200 -> id
	400 -> DB Conn or Query err
	*/
	api.POST("/User/find/id", postFindID)

	/*  Reply			200 -> List<post>
	400 -> DB Conn or Query err
	*/
	api.GET("/Projects", getProjects)

	/*  Reply			200 -> project obj
	400 -> DB Conn or Query err
	*/
	api.GET("/Project/", getProject)

	/*  Reply			200 -> Get List<post>
	400 -> DB Conn or Query err
	*/
	api.GET("/Projects/category", getCategory)

	/*  Reply			200 -> null
	400 -> DB Conn or Query err
	*/
	api.POST("/Projects/add", postAddProject)

	/*  Reply			200 -> null
	400 -> DB Conn or Query err
	*/
	api.POST("/Projects/permit", postPermitJoin)

	/*  Reply			200 -> null
	400 -> DB Conn or Query err
	*/
	api.POST("/Projects/deny", postDenyJoin)

	/*  Reply			200 -> null
	400 -> DB Conn or Query err
	*/
	api.POST("/Projects/join", postJoin)

	/*  Reply			200 -> List<msg>
	400 -> DB Conn or Query err
	*/
	api.GET("/Refresh", getRefresh)

	/*  Reply			200 -> null
	400 -> DB Conn or Query err
	*/
	api.POST("/announce", postAnnounce)
}
