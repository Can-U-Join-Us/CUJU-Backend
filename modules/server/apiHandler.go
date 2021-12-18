package server

import (
	. "github.com/Can-U-Join-Us/CUJU-Backend/modules/api"
	"github.com/gin-gonic/gin"
)

func getLogin(c *gin.Context) {
	uid, token, err := LoginUser(c)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
	} else {
		c.JSON(200, gin.H{"error": nil, "token": token, "uid": uid})
	}
}
func postModifyPW(c *gin.Context) {
	err := ModifyPW(c)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
	} else {
		c.JSON(200, gin.H{"error": nil})
	}
}
func postModifyProfile(c *gin.Context) {
	err := ModifyProfile(c)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
	} else {
		c.JSON(200, gin.H{"error": nil})
	}
}
func postLogout(c *gin.Context) {
	err := LogoutUser(c)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
	} else {
		c.JSON(200, gin.H{"error": nil})
	}
}
func postRegister(c *gin.Context) {
	err := RegisterUser(c)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
	} else {
		c.JSON(200, gin.H{"error": nil})
	}
}
func postFindPW(c *gin.Context) {
	err := FindUserPW(c)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
	} else {
		c.JSON(200, gin.H{"error": nil})
	}
}
func postFindID(c *gin.Context) {
	id, err := FindUserId(c)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
	} else {
		c.JSON(200, gin.H{"error": nil, "id": id})
	}
}
func getProjects(c *gin.Context) {
	posts, err := GetProjectList(c)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
	} else {
		c.JSON(200, gin.H{"error": nil, "projects": posts})
	}
}
func getProject(c *gin.Context) {
	pd, err := GetProjectDetail(c)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
	} else {
		c.JSON(200, gin.H{"error": nil, "projects": pd})
	}
}
func getCategory(c *gin.Context) {
	posts, err := GetCategory(c)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
	} else {
		c.JSON(200, gin.H{"error": nil, "projects": posts})
	}
}
func postAddProject(c *gin.Context) {
	pid, err := AddProject(c)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
	} else {
		c.JSON(200, gin.H{"error": nil, "pid": pid})
	}
}
func postPermitJoin(c *gin.Context) {
	err := PermitProject(c)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
	} else {
		c.JSON(200, gin.H{"error": nil})
	}
}
func postDenyJoin(c *gin.Context) {
	err := DenyProject(c)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
	} else {
		c.JSON(200, gin.H{"error": nil})
	}
}
func postJoin(c *gin.Context) {
	err := JoinProject(c)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
	} else {
		c.JSON(200, gin.H{"error": nil})
	}
}
func getRefresh(c *gin.Context) {
	msgList, err := RefreshMsg(c)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
	} else {
		c.JSON(200, gin.H{"error": nil, "Message": msgList})
	}
}
func postAnnounce(c *gin.Context) {
	err := Announcement(c)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
	} else {
		c.JSON(200, gin.H{"error": nil})
	}
}
