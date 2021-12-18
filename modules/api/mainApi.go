package api

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"

	ErrChecker "github.com/Can-U-Join-Us/CUJU-Backend/modules/errors"
	"github.com/Can-U-Join-Us/CUJU-Backend/modules/storage"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func GetProjectList(c *gin.Context) ([]Project, error) {
	db := storage.DB()
	var length int
	_ = db.QueryRow(`select count(*) from project_post`).Scan(&length)
	if length == 0 {
		return []Project{}, errors.New("Nothing to show")
	}
	rows, err := db.Query(`select * from project_post`)
	if err := ErrChecker.Check(err); err != nil {
		return []Project{}, err
	}
	defer rows.Close()
	projects := make([]Project, 0)
	var pos Project
	for rows.Next() {
		err := rows.Scan(&pos.PID, &pos.UID, &pos.TITLE,
			&pos.TOTAL, &pos.DESCRIPTION, &pos.DUE, &pos.TERM, &pos.PATH)
		if err := ErrChecker.Check(err); err != nil {
			return []Project{}, err
		}
		projects = append(projects, pos)
	}
	return projects, nil
}
func GetProjectDetail(c *gin.Context) (Project, error) {
	var pd Project
	pid := c.Request.Header.Get("pid")
	pd.PID, _ = strconv.Atoi(pid)
	db := storage.DB()

	err := db.QueryRow(`select p.title, p.description, p.total, p.path, p.term, p.due, m.fe, m.be, m.aos, m.ios, m.pm, m.designer, m.devops, m.etc from project_post p join member m on p.pid = m.pid where p.pid =`+pid).Scan(&pd.TITLE, &pd.DESCRIPTION, &pd.TOTAL, &pd.PATH, &pd.TERM, &pd.DUE, &pd.FE, &pd.BE, &pd.AOS, &pd.IOS, &pd.PM, &pd.DESIGNER, &pd.DEVOPS, &pd.ETC)
	if err == nil {
		return Project{}, err
	}
	return pd, nil
}
func GetCategory(c *gin.Context) ([]Project, error) {
	db := storage.DB()
	category := c.Request.Header.Get("category")
	join := category + "_join"
	_ = join

	query := `select count(*) from project_post`
	var length int
	_ = db.QueryRow(query).Scan(&length)
	if length == 0 {
		return []Project{}, errors.New("Nothing to show")
	}
	projects := make([]Project, 0)
	var pos Project
	query = `select project_post.PID,TITLE,DESCRIPTION,TOTAL,TERM,DUE,PATH from project_post join member on project_post.pid = member.pid and  member.` + join + ` < member.` + category

	// category별 참여 인원이 덜 차있는 게시물만 리턴
	rows, err := db.Query(query)
	if err := ErrChecker.Check(err); err != nil {
		return []Project{}, err
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&pos.PID, &pos.TITLE, &pos.DESCRIPTION,
			&pos.TOTAL, &pos.TERM, &pos.DUE, &pos.PATH)
		if err := ErrChecker.Check(err); err != nil {
			return []Project{}, err
		}
		projects = append(projects, pos)
	}
	return projects, nil
}
func AddProject(c *gin.Context) (int, error) {
	var reqBody AddProjectForm
	err := c.ShouldBindWith(&reqBody, binding.FormMultipart)
	if err != nil {
		return -1, err
	}
	fmt.Println(reqBody)
	return -1, errors.New("test")
	c.Request.FormValue("uid")
	reqBody.UID, _ = strconv.Atoi(c.Request.Form.Get("uid"))
	reqBody.TOTAL, _ = strconv.Atoi(c.Request.Form.Get("total"))
	reqBody.TERM, _ = strconv.Atoi(c.Request.Form.Get("term"))
	reqBody.FE, _ = strconv.Atoi(c.Request.Form.Get("fe"))
	reqBody.BE, _ = strconv.Atoi(c.Request.Form.Get("be"))
	reqBody.AOS, _ = strconv.Atoi(c.Request.Form.Get("aos"))
	reqBody.IOS, _ = strconv.Atoi(c.Request.Form.Get("ios"))
	reqBody.PM, _ = strconv.Atoi(c.Request.Form.Get("pm"))
	reqBody.DESIGNER, _ = strconv.Atoi(c.Request.Form.Get("designer"))
	reqBody.DEVOPS, _ = strconv.Atoi(c.Request.Form.Get("devops"))
	reqBody.ETC, _ = strconv.Atoi(c.Request.Form.Get("etc"))
	reqBody.TITLE = c.Request.Form.Get("title")
	reqBody.DESC = c.Request.Form.Get("desc")
	reqBody.DUE = c.Request.Form.Get("due")
	reqBody.PATH = c.Request.Form.Get("path")
	reqBody.FE_desc = c.Request.Form.Get("fe_desc")
	reqBody.BE_desc = c.Request.Form.Get("be_desc")
	reqBody.AOS_desc = c.Request.Form.Get("aos_desc")
	reqBody.IOS_desc = c.Request.Form.Get("ios_desc")
	reqBody.PM_desc = c.Request.Form.Get("pm_desc")
	reqBody.DESIGNER_desc = c.Request.Form.Get("designer_desc")
	reqBody.DEVOPS_desc = c.Request.Form.Get("devops_desc")
	reqBody.ETC_desc = c.Request.Form.Get("etc_desc")

	file, _, err := c.Request.FormFile("hello")

	var pid int
	db := storage.DB()
	db.QueryRow(`select pid from project_post order by pid desc limit 1`).Scan(&pid)
	path := `http://192.168.0.7/project_img/` + strconv.Itoa(pid+1) + `.png`

	if err != nil {
		return -1, err
	}
	_, err = db.Exec(`Insert into project_post(UID,TITLE,TOTAL,DESCRIPTION, TERM, DUE, PATH) values(?,?,?,?,?,?,?)`, reqBody.UID, reqBody.TITLE, reqBody.TOTAL, reqBody.DESC, reqBody.TERM, reqBody.DUE, path)
	if err := ErrChecker.Check(err); err != nil {
		fmt.Println("err 1")
		return -1, err
	}

	fmt.Println("PATH : " + path)
	path = `/Users/macbook/Sites/project_img/` + strconv.Itoa(pid+1) + `.png`
	dst, err := os.Create(path)
	fmt.Println(dst)
	defer dst.Close()

	if err != nil {
		return -1, err
	}
	if _, err := io.Copy(dst, file); err != nil {
		return -1, err
	}
	fmt.Println(pid+1, "\n", reqBody)
	_, err = db.Exec(`Insert into member (pid,fe,be,aos,ios,pm,designer,devops,etc) values(?,?,?,?,?,?,?,?,?)`, pid+1,
		reqBody.FE, reqBody.BE, reqBody.AOS, reqBody.IOS, reqBody.PM, reqBody.DESIGNER,
		reqBody.DEVOPS, reqBody.ETC)
	if err := ErrChecker.Check(err); err != nil {
		return -1, err
	}
	return pid, nil
}
func DenyProject(c *gin.Context) error {
	var reqBody ReplyJoinForm
	err := c.ShouldBindJSON(&reqBody)
	if err := ErrChecker.Check(err); err != nil {
		return err
	}
	db := storage.DB()
	var count int
	_ = db.QueryRow(`select count(*) from join_queue where uid = ` + strconv.Itoa(reqBody.UID) + ` and pid = ` + strconv.Itoa(reqBody.PID) + ` and confirm = 0`).Scan(&count)
	if count == 0 {
		return errors.New("Nothing")
	}
	_, err = db.Exec(`update join_queue set result = 2, confirm = 1 where uid =  ` + strconv.Itoa(reqBody.UID))

	if err != nil {
		return err
	}
	return nil
}
func PermitProject(c *gin.Context) error {
	var reqBody ReplyJoinForm
	err := c.ShouldBindJSON(&reqBody)
	if err := ErrChecker.Check(err); err != nil {
		return err
	}
	db := storage.DB()
	var count int
	_ = db.QueryRow(`select count(*) from join_queue where uid = ` + strconv.Itoa(reqBody.UID) + ` and pid = ` + strconv.Itoa(reqBody.PID) + ` and confirm = 0`).Scan(&count)
	if count == 0 {
		return errors.New("Nothing")
	}
	_, err = db.Exec(`update join_queue set result = 1, confirm = 1 where uid =  ` + strconv.Itoa(reqBody.UID))

	if err != nil {
		return err
	}
	return nil
}
func JoinProject(c *gin.Context) error {
	var reqBody JoinForm
	err := c.ShouldBindJSON(&reqBody)
	if err := ErrChecker.Check(err); err != nil {
		return err
	}
	db := storage.DB()
	_, err = db.Exec("insert into join_queue(pid,uid,category) value(?,?,?)", reqBody.PID, reqBody.UID, reqBody.CATEGORY)
	if err != nil {
		return err
	}
	return nil
}
func GetNumProject(c *gin.Context) (int, error) {
	var reqBody JoinForm
	err := c.ShouldBindJSON(&reqBody)
	if err := ErrChecker.Check(err); err != nil {
		return -1, err
	}
	db := storage.DB()
	_, err = db.Exec("insert into join_queue(pid,uid,category) value(?,?,?)", reqBody.PID, reqBody.UID, reqBody.CATEGORY)
	if err != nil {
		return -1, err
	}
	return -1, nil
}
