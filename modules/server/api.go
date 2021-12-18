package server

import (
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/smtp"
	"os"
	"strconv"
	"time"

	ErrChecker "github.com/Can-U-Join-Us/CUJU-Backend/modules/errors"
	storage "github.com/Can-U-Join-Us/CUJU-Backend/modules/storage"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type resgisterForm struct {
	Email string `json:"email"`
	PW    string `json:"pw"`
	Name  string `json:"name"`
	Phone string `json:"phone"`
}
type loginForm struct {
	ID string `json:"Email"`
	PW string `json:"pw"`
}
type modifyForm struct {
	UID int    `json:"uid"`
	PW  string `json:"pw"`
	NEW string `json:"new"`
}
type project struct {
	PID         int    `json:"pid"`
	UID         int    `json:"uid"`
	TITLE       string `json:"title"`
	DESCRIPTION string `json:"desc"`
	TOTAL       int    `json:"total"`
	TERM        int    `json:"term"`
	DUE         string `json:"due"`
	PATH        string `json:"path"`
	FE          int    `json:"fe"`
	BE          int    `json:"be"`
	AOS         int    `json:"aos"`
	IOS         int    `json:"ios"`
	PM          int    `json:"pm"`
	DESIGNER    int    `json:"designer"`
	DEVOPS      int    `json:"devops"`
	ETC         int    `json:"etc"`
}
type addProjectForm struct {
	UID           int    `json:"uid"`
	TITLE         string `json:"title"`
	DESC          string `json:"desc"`
	TOTAL         int    `json:"total"`
	TERM          int    `json:"term"`
	DUE           string `json:"due"`
	PATH          string `json:"path"`
	FE            int    `json:"fe"`
	BE            int    `json:"be"`
	AOS           int    `json:"aos"`
	IOS           int    `json:"ios"`
	PM            int    `json:"pm"`
	DESIGNER      int    `json:"designer"`
	DEVOPS        int    `json:"devops"`
	ETC           int    `json:"etc"`
	FE_desc       string `json:"fe_desc"`
	BE_desc       string `json:"be_desc"`
	AOS_desc      string `json:"aos_desc"`
	IOS_desc      string `json:"ios_desc"`
	PM_desc       string `json:"pm_desc"`
	DESIGNER_desc string `json:"designer_desc"`
	DEVOPS_desc   string `json:"devops_desc"`
	ETC_desc      string `json:"etc_desc"`
}
type joinForm struct {
	PID      int    `json:"pid"`
	UID      int    `json:"uid"`
	CATEGORY string `json:"category"`
}
type replyJoinForm struct {
	PID int `json:"pid"`
	UID int `json:"uid"`
}
type announce struct {
	TITLE   string `json:"title"`
	CONTENT string `json:"content"`
}
type msg struct {
	TYPE    int    `json:"type"`
	TITLE   string `json:"title"`
	CONTENT string `json:"content"`
	PID     int    `json:"pid"`
	UID     int    `json:"uid"`
}

func registerUser(c *gin.Context) error {
	var reqBody resgisterForm
	err := c.ShouldBindJSON(&reqBody)
	if err := ErrChecker.Check(err); err != nil {
		return err
	}
	db := storage.DB()
	var count int
	_ = db.QueryRow(`Select count(*) from user where email = "` + reqBody.Email + `"`).Scan(&count)

	if count > 0 {
		return errors.New("ID Duplicate")
	}
	stmt, err := db.Prepare("Insert into user(email,PW,Name,Phone) values(?,?,?,?)")
	if err := ErrChecker.Check(err); err != nil {
		return err
	}
	rs, err := stmt.Exec(reqBody.Email, reqBody.PW, reqBody.Name, reqBody.Phone)
	if err := ErrChecker.Check(err); err != nil {
		return err
	}
	_ = rs
	defer stmt.Close()
	return nil
}

func loginUser(c *gin.Context) (uint64, map[string]string, error) {
	var reqBody loginForm
	err := c.ShouldBindJSON(&reqBody)
	if err := ErrChecker.Check(err); err != nil {
		return 0, map[string]string{}, err
	}
	db := storage.DB()
	query := `select uid, PW from user where Email = "` + reqBody.ID + `"`
	var pw string
	row := db.QueryRow(query)
	var uid uint64
	err = row.Scan(&uid, &pw)
	if err := ErrChecker.Check(err); err != nil {
		return 0, map[string]string{}, errors.New("ID")
	}
	if reqBody.PW != pw { // PW 가 다르면 PW 가 다르다는 오류 반환
		return 0, map[string]string{}, errors.New("PW")
	}
	ts, err := createToken(uid)
	if err := ErrChecker.Check(err); err != nil {
		return 0, map[string]string{}, err
	}
	err = createAuth(uid, ts) // Redis 토큰 메타데이터 저장
	if err := ErrChecker.Check(err); err != nil {
		return 0, map[string]string{}, err
	}
	tokens := map[string]string{
		"access_token":  ts.AccessToken,
		"refresh_token": ts.RefreshToken,
	}
	return uid, tokens, nil
}
func logoutUser(c *gin.Context) error {
	// request header 에 담긴 access & refresh token을 검증 후 redis 에서 삭제
	au, ru, err := ExtractBothTokenMetadata(c.Request)
	if err := ErrChecker.Check(err); err != nil {
		return err
	}
	deleted, err := DeleteAuth(au.AccessUuid, ru.RefreshUuid)
	if err := ErrChecker.Check(err); err != nil || deleted == 0 {
		return err
	}
	return nil
}
func findUserPW(c *gin.Context) error {
	var reqBody struct {
		ID string `json:"email"`
	}
	err := c.ShouldBindJSON(&reqBody)
	if err := ErrChecker.Check(err); err != nil {
		return err
	}
	var email string
	var name string
	var count int
	db := storage.DB()
	query := `select count(*),email,name from user where Email = "` + reqBody.ID + `"`
	row := db.QueryRow(query)
	err = row.Scan(&count, &email, &name)
	if err := ErrChecker.Check(err); err != nil {
		return err
	}
	if count == 0 {
		return errors.New("Invalid id")
	}
	pwByte := []byte{}
	for i := 0; i < 10; i++ {
		rand.Seed(time.Now().UnixNano())
		if a := rand.Intn(5); a < 4 {
			pwByte = append(pwByte, byte(rand.Intn(25)+97))
		} else {
			pwByte = append(pwByte, byte(rand.Intn(10)+48))
		}

	}
	pw := string(pwByte)
	query = `update user set pw ="` + pw + `" where email = "` + reqBody.ID + `"`
	_, err = db.Exec(query)
	if err := ErrChecker.Check(err); err != nil {
		return err
	}
	auth := smtp.PlainAuth("", "cujuserver@gmail.com", "cujuroot1!", "smtp.gmail.com")
	from := "cujuserver@gmail.com"
	to := []string{reqBody.ID}
	headerSubject := "Subject: 같이할래 임시 PW 발급\r\n"
	headerBlank := "\r\n"

	body :=
		`안녕하세요 
	
같이할래 플랫폼을 이용해주셔서 감사합니다.

` + name + `님의 임시 PW입니다.

PW:` + pw
	msg := []byte(headerSubject + headerBlank + body)
	err = smtp.SendMail("smtp.gmail.com:587", auth, from, to, msg)
	if err != nil {
		panic(err)
	}
	return nil
}
func findUserId(c *gin.Context) (string, error) {
	var reqBody struct {
		PHONE string `json:"phone"`
	}
	err := c.ShouldBindJSON(&reqBody)
	if err := ErrChecker.Check(err); err != nil {
		return "", err
	}
	db := storage.DB()
	query := `select email from user where phone = "` + reqBody.PHONE + `"`
	row := db.QueryRow(query)
	var email string
	err = row.Scan(&email)
	if err := ErrChecker.Check(err); err != nil {
		return "", err
	}
	return email, nil
}
func modifyPW(c *gin.Context) error {
	var reqBody modifyForm
	err := c.ShouldBindJSON(&reqBody)
	if err := ErrChecker.Check(err); err != nil {
		return err
	}
	db := storage.DB()
	var count int
	uid := strconv.Itoa(reqBody.UID)
	query := `select count(*) from user where uid = ` + uid + ` and pw = "` + reqBody.PW + `"`
	_ = db.QueryRow(query).Scan(&count)
	if count == 0 {
		return errors.New("Invalid pw")
	}
	query = `update user set pw ="` + reqBody.NEW + `" where uid = ` + uid
	fmt.Println(query)
	res, err := db.Exec(query)
	if err := ErrChecker.Check(err); err != nil {
		return err
	}
	nRow, _ := res.RowsAffected()
	fmt.Println("update count : ", nRow)
	return nil
}
func modifyProfile(c *gin.Context) error {
	au, _, err := ExtractBothTokenMetadata(c.Request)
	if err := ErrChecker.Check(err); err != nil {
		return err
	}
	fmt.Println("This user id is ", au.UserId)
	return nil
}
func getProjectList(c *gin.Context) ([]project, error) {
	db := storage.DB()
	query := `select count(*) from project_post`

	var length int
	_ = db.QueryRow(query).Scan(&length)
	if length == 0 {
		return []project{}, errors.New("Nothing to show")
	}
	query = `select * from project_post`
	rows, err := db.Query(query)
	if err := ErrChecker.Check(err); err != nil {
		return []project{}, err
	}
	defer rows.Close()

	projects := make([]project, 0)
	var pos project
	for rows.Next() {
		err := rows.Scan(&pos.PID, &pos.UID, &pos.TITLE,
			&pos.TOTAL, &pos.DESCRIPTION, &pos.DUE, &pos.TERM, &pos.PATH)
		if err := ErrChecker.Check(err); err != nil {
			return []project{}, err
		}
		projects = append(projects, pos)
	}
	return projects, nil
}
func getProjectDetail(c *gin.Context) (project, error) {
	var pd project
	pid := c.Request.Header.Get("pid")
	pd.PID, _ = strconv.Atoi(pid)
	db := storage.DB()

	err := db.QueryRow(`select p.title, p.description, p.total, p.path, p.term, p.due, m.fe, m.be, m.aos, m.ios, m.pm, m.designer, m.devops, m.etc from project_post p join member m on p.pid = m.pid where p.pid =`+pid).Scan(&pd.TITLE, &pd.DESCRIPTION, &pd.TOTAL, &pd.PATH, &pd.TERM, &pd.DUE, &pd.FE, &pd.BE, &pd.AOS, &pd.IOS, &pd.PM, &pd.DESIGNER, &pd.DEVOPS, &pd.ETC)
	if err == nil {
		return project{}, err
	}
	return pd, nil
}
func getCategory(c *gin.Context) ([]project, error) {
	db := storage.DB()
	category := c.Request.Header.Get("category")
	join := category + "_join"
	_ = join

	query := `select count(*) from project_post`
	var length int
	_ = db.QueryRow(query).Scan(&length)
	if length == 0 {
		return []project{}, errors.New("Nothing to show")
	}
	projects := make([]project, 0)
	var pos project
	query = `select project_post.PID,TITLE,DESCRIPTION,TOTAL,TERM,DUE,PATH from project_post join member on project_post.pid = member.pid and  member.` + join + ` < member.` + category

	// category별 참여 인원이 덜 차있는 게시물만 리턴
	rows, err := db.Query(query)
	if err := ErrChecker.Check(err); err != nil {
		return []project{}, err
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&pos.PID, &pos.TITLE, &pos.DESCRIPTION,
			&pos.TOTAL, &pos.TERM, &pos.DUE, &pos.PATH)
		if err := ErrChecker.Check(err); err != nil {
			return []project{}, err
		}
		projects = append(projects, pos)
	}
	fmt.Println(projects)
	return projects, nil
}
func addProject(c *gin.Context) (int, error) {
	var reqBody addProjectForm
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
func denyProject(c *gin.Context) error {
	var reqBody replyJoinForm
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
func permitProject(c *gin.Context) error {
	var reqBody replyJoinForm
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
func joinProject(c *gin.Context) error {
	var reqBody joinForm
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
func getNumProject(c *gin.Context) (int, error) {
	var reqBody joinForm
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
func refreshMsg(c *gin.Context) ([]msg, error) {
	uid := c.Request.Header.Get("uid")
	fmt.Println("uid : ", uid)
	db := storage.DB()
	var pid, id, result int
	var name, email, title string
	query := `select p.pid, j.uid,u.name, u.email, j.result from join_queue j join project_post p on j.pid = p.pid and p.uid = ` + uid + ` left join user u on u.uid = j.uid where j.result = 0 and j.confirm = 0`
	rows, err := db.Query(query)
	if err != nil {
		return []msg{}, err
	}
	defer rows.Close()
	msgList := make([]msg, 0)
	res, _ := getAnnouncement(c)
	for i := 0; i < len(res); i++ {
		msgList = append(msgList, res[i])
	}
	var m msg
	for rows.Next() {
		if err := rows.Scan(&pid, &id, &name, &email, &result); err != nil {
			return []msg{}, err
		}
		uuid, _ := strconv.Atoi(uid)
		m = msg{
			TYPE:    1,
			TITLE:   strconv.Itoa(pid) + ` 번 프로젝트 참여 신청입니다 !`,
			CONTENT: email,
			PID:     pid,
			UID:     uuid,
		}
		msgList = append(msgList, m)
	}
	query = `select p.pid,p.title,j.result from join_queue j join project_post p on j.pid = p.pid and p.uid = ` + uid + ` left join user u on u.uid = j.uid where j.result != 0 and confirm = 0`
	rows, err = db.Query(query)
	if err != nil {
		return []msg{}, err
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&id, &title, &result); err != nil {
			return []msg{}, err
		}
		var ctt string
		if result == 1 { // permission
			ctt = title + `(` + strconv.Itoa(pid) + `) 프로젝트 : 승인`
		} else { // deny
			ctt = title + `(` + strconv.Itoa(pid) + `) 프로젝트 : 거절`
		}
		m = msg{
			TYPE:    2,
			TITLE:   strconv.Itoa(id) + ` 번 프로젝트 주최자의 답변이 도착했습니다 !`,
			CONTENT: ctt,
			PID:     0,
			UID:     0,
		}

		msgList = append(msgList, m)
	}

	if len(msgList) == 0 {
		return []msg{}, errors.New("Nothing")
	}
	query = `update join_queue j join project_post p on j.pid= p.pid and p.uid = 7 left join user u on u.uid = j.uid  set j.confirm = 1 where j.result > 0`
	_, err = db.Exec(query) // 참여 결과 알림 메세지는 읽음 상태로 변경
	if err != nil {
		return []msg{}, err
	}
	return msgList, nil
}
func getAnnouncement(c *gin.Context) ([]msg, error) {
	annoList := make([]msg, 0)
	db := storage.DB()
	query := `select title,content from announce where send = false`
	rows, err := db.Query(query)
	if err != nil {
		return []msg{}, errors.New("Nothing")
	}
	defer rows.Close()

	for rows.Next() {
		var title, content string
		if err := rows.Scan(&title, &content); err != nil {
			return []msg{}, err
		}
		m := msg{
			TYPE:    0,
			TITLE:   title,
			CONTENT: content,
			PID:     0,
			UID:     0,
		}
		annoList = append(annoList, m)
	}

	return annoList, nil
}
func postAnnouncement(c *gin.Context) error {
	var reqBody announce
	err := c.ShouldBindJSON(&reqBody)
	if err := ErrChecker.Check(err); err != nil {
		return err
	}
	db := storage.DB()
	_, err = db.Exec(`insert into announce values(?,?,?)`, reqBody.TITLE, reqBody.CONTENT, 0)
	if err != nil {
		return err
	}
	return nil
}
