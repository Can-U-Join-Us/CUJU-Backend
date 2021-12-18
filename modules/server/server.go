package server

import (
	"fmt"

	_ "github.com/Can-U-Join-Us/CUJU-Backend/modules/storage"
	. "github.com/Can-U-Join-Us/CUJU-Backend/modules/token"
	"github.com/gin-gonic/gin"
)

const port = ":4000"

func init() { // local : 4000 호스팅 시작
	r := gin.Default()
	if err := RedisInit(); err != nil {
		panic(fmt.Errorf("Fatal error : redis is off status \n"))
	}
	api := r.Group("/api")
	api.Use(dummy)
	registerApiHandlers(api)

	r.Run(port)
}
func dummy(c *gin.Context) {
	fmt.Println("Access Token Check Stage")
}
func registerApiHandlers(api *gin.RouterGroup) {
	/*  Reply			200 -> token
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
	/*  Reply			200 -> find pw success
	400 -> DB Conn or Query err
	*/
	api.POST("/User/find/pw", postFindPW)
	/*  Reply			200 -> find id success
	400 -> DB Conn or Query err
	*/
	api.POST("/User/find/id", postFindID)
	/*  Reply			200 -> Get List<post> success
	400 -> DB Conn or Query err
	*/
	api.GET("/Projects", getProjects)
	/*  Reply			200 -> Get project obj success
	400 -> DB Conn or Query err
	*/
	api.GET("/Project/", getProject)
	/*  Reply			200 -> Get List<post> success
	400 -> DB Conn or Query err
	*/
	api.GET("/Projects/category", getCategory)
	/*  Reply			200 -> Add post success
	400 -> DB Conn or Query err
	*/
	api.POST("/Projects/add", postAddProject)
	/*  Reply			200 -> Permit post success
	400 -> DB Conn or Query err
	*/
	api.POST("/Projects/permit", postPermitJoin)
	/*  Reply			200 -> Deny post success
	400 -> DB Conn or Query err
	*/
	api.POST("/Projects/deny", postDenyJoin)
	/*  Reply			200 -> Join post success
	400 -> DB Conn or Query err
	*/
	api.POST("/Projects/join", postJoin)
	/*  Reply			200 -> Get List<msg> success
	400 -> DB Conn or Query err
	*/
	api.GET("/Refresh", getRefresh)
	/*  Reply			200 -> POST announceMSG success
	400 -> DB Conn or Query err
	*/
	api.POST("/announce", postAnnounce)
}
