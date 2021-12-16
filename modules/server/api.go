package server

import (
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/smtp"
	"os"
	"strconv"
	"strings"
	"time"

	ErrChecker "github.com/Can-U-Join-Us/CUJU-Backend/modules/errors"
	storage "github.com/Can-U-Join-Us/CUJU-Backend/modules/storage"
	"github.com/gin-gonic/gin"
)

func registerUser(c *gin.Context) error {
	var reqBody struct {
		Email string `json:"email"`
		PW    string `json:"pw"`
		Name  string `json:"name"`
		Phone string `json:"phone"`
	}
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
	var reqBody struct {
		ID string `json:"Email"`
		PW string `json:"pw"`
	}
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
		ID string `json:"Email"`
	}
	err := c.ShouldBindJSON(&reqBody)
	if err := ErrChecker.Check(err); err != nil {
		return err
	}
	var email string
	var name string
	db := storage.DB()
	query := `select email,name from user where Email = "` + reqBody.ID + `"`
	row := db.QueryRow(query)
	err = row.Scan(&email, &name)
	if err := ErrChecker.Check(err); err != nil {
		return err
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
	res, err := db.Exec(query)
	if err := ErrChecker.Check(err); err != nil {
		return err
	}
	nRow, err := res.RowsAffected()
	fmt.Println("update count : ", nRow)
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
	var email string
	var reqBody struct {
		PHONE string `json:"phone"`
	}
	err := c.ShouldBindJSON(&reqBody)
	if err := ErrChecker.Check(err); err != nil {
		return "", err
	}
	fmt.Println(reqBody)
	db := storage.DB()
	query := `select email from user where phone = "` + reqBody.PHONE + `"`
	row := db.QueryRow(query)
	err = row.Scan(&email)
	if err := ErrChecker.Check(err); err != nil {
		return "", err
	}
	return email, nil
}
func modifyPW(c *gin.Context) error {
	var reqBody struct {
		UID int    `json:"uid"`
		PW  string `json:"pw"`
		NEW string `json:"new"`
	}
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
	projectDetail := project{}

	return projectDetail, nil
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
	return projects, nil
}
func addProject(c *gin.Context) (int, error) {
	val := strings.Repeat("?,", 16)
	val += "?)"
	val = "(" + val
	var reqBody struct {
		UID           uint   `json:"uid"`
		TITLE         string `json:"title"`
		DESC          string `json:"desc"`
		TOTAL         uint   `json:"total"`
		TERM          uint   `json:"term"`
		DUE           string `json:"due"`
		PATH          string `json:"path"`
		FE            uint   `json:"fe"`
		BE            uint   `json:"be"`
		AOS           uint   `json:"aos"`
		IOS           uint   `json:"ios"`
		PM            uint   `json:"pm"`
		DESIGNER      uint   `json:"designer"`
		DEVOPS        uint   `json:"devops"`
		ETC           uint   `json:"etc"`
		FE_desc       string `json:"fe_desc"`
		BE_desc       string `json:"be_desc"`
		AOS_desc      string `json:"aos_desc"`
		IOS_desc      string `json:"ios_desc"`
		PM_desc       string `json:"pm_desc"`
		DESIGNER_desc string `json:"designer_desc"`
		DEVOPS_desc   string `json:"devops_desc"`
		ETC_desc      string `json:"etc_desc"`
	}
	file, handler, err := c.Request.FormFile("hello")
	if err != nil {
		return -1, err
	}
	fmt.Println(c.Request.FormValue("uid"))
	fmt.Println()
	fmt.Println(file)
	fmt.Println()
	fmt.Println(handler)
	fmt.Println()
	dst, err := os.Create(`~/Sites/project_img/` + "5" + `.png`)
	fmt.Println(dst)
	defer dst.Close()

	if err != nil {
		return -1, err
	}
	if _, err := io.Copy(dst, file); err != nil {
		return -1, err
	}
	err = c.ShouldBindJSON(&reqBody)

	if err := ErrChecker.Check(err); err != nil {
		return -1, err
	}

	db := storage.DB()
	_, err = db.Exec(`Insert into project_post(UID,TITLE,TOTAL,DESCRIPTION, TERM, DUE, PATH) values(?,?,?,?,?,?,?)`, reqBody.UID, reqBody.TITLE, reqBody.TOTAL, reqBody.DESC, reqBody.TERM, reqBody.DUE, reqBody.PATH)
	if err := ErrChecker.Check(err); err != nil {
		return -1, err
	}
	var pid int

	db.QueryRow(`select pid from project_post order by pid desc limit 1`).Scan(&pid)

	_, err = db.Exec(`Insert into member (pid,fe,be,aos,ios,pm,designer,devops,etc,fe_desc,be_desc,aos_desc,ios_desc,pm_desc,designer_desc,devops_desc,etc_desc) values`+val, pid,
		reqBody.FE, reqBody.BE, reqBody.AOS, reqBody.IOS, reqBody.PM, reqBody.DESIGNER,
		reqBody.DEVOPS, reqBody.ETC, `"`+reqBody.FE_desc+`"`, `"`+reqBody.BE_desc+`"`, `"`+reqBody.AOS_desc+`"`,
		`"`+reqBody.IOS_desc+`"`, `"`+reqBody.PM_desc+`"`, `"`+reqBody.DESIGNER_desc+`"`, `"`+reqBody.DEVOPS_desc+`"`, `"`+reqBody.ETC_desc+`"`)
	if err := ErrChecker.Check(err); err != nil {
		return -1, err
	}
	return pid, nil
}
func denyProject(c *gin.Context) error {
	var reqBody struct {
		PID int `json:"pid"`
		UID int `json:"uid"`
	}
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
	var reqBody struct {
		PID int `json:"pid"`
		UID int `json:"uid"`
	}
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
	var reqBody struct {
		PID      int    `json:"pid"`
		UID      int    `json:"uid"`
		CATEGORY string `json:"category"`
	}
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
	var reqBody struct {
		PID      int    `json:"pid"`
		UID      int    `json:"uid"`
		CATEGORY string `json:"category"`
	}
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
			TITLE:   strconv.Itoa(id) + ` 번 프로젝트 참여 신청입니다 !`,
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
	var reqBody struct {
		TITLE   string `json:"title"`
		CONTENT string `json:"content"`
	}
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

type msg struct {
	TYPE    int    `json:"type"`
	TITLE   string `json:"title"`
	CONTENT string `json:"content"`
	PID     int    `json:"pid"`
	UID     int    `json:"uid"`
}
type project struct {
	PID         uint   `json:"pid"`
	UID         uint   `json:"uid"`
	TITLE       string `json:"title"`
	DESCRIPTION string `json:"desc"`
	PATH        string `json:"path"`
	TOTAL       uint   `json:"total"`
	TERM        uint   `json:"term"`
	DUE         string `json:"due"`
}

type member struct {
	PID           uint
	FE            uint
	BE            uint
	AOS           uint
	IOS           uint
	PM            uint
	DESIGNER      uint
	DEVOPS        uint
	ETC           uint
	FE_desc       string
	BE_desc       string
	AOS_desc      string
	IOS_desc      string
	PM_desc       string
	DESIGNER_desc string
	DEVOPS_desc   string
	ETC_desc      string
}
