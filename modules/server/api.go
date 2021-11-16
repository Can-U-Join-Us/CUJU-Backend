package server

import (
	"errors"

	storage "github.com/Can-U-Join-Us/CUJU-Backend/modules/storage"
	"github.com/gin-gonic/gin"
)

func pingTest(c *gin.Context) string {
	return "PONG"
}
func dbConnectionTest(c *gin.Context) error {
	db := storage.DB()
	if db != nil {
		return nil
	}
	return errors.New("Not connected")
}
func registerUser(c *gin.Context) error {
	var reqBody struct {
		ID    string `json:"id"`
		PW    string `json:"pw"`
		Name  string `json:"name"`
		Email string `json:"email"`
	}
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		return err
	}
	db := storage.DB()
	stmt, err := db.Prepare("Insert into user(ID,PW,Name,Email) values(?,?,?,?)")
	if err != nil {
		return err
	}
	rs, err := stmt.Exec(reqBody.ID, reqBody.PW, reqBody.Name, reqBody.Email)
	if err != nil {
		return err
	}
	_ = rs
	defer stmt.Close()
	return nil
}

func loginUser(c *gin.Context) (map[string]string, error) {
	var reqBody struct {
		ID string `json:"id"`
		PW string `json:"pw"`
	}
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		return map[string]string{}, err
	}
	db := storage.DB()
	query := `select uid, PW from user where ID = "` + reqBody.ID + `"`
	var pw string
	row := db.QueryRow(query)
	var uid uint64
	err := row.Scan(&uid, &pw)
	if err != nil { // ID 가 없으면 ID 가 없다는 오류 반환
		return map[string]string{}, errors.New("ID")
	}
	if reqBody.PW != pw { // PW 가 다르면 PW 가 다르다는 오류 반환
		return map[string]string{}, errors.New("PW")
	}
	ts, err := createToken(uid)
	if err != nil {
		return map[string]string{}, err
	}
	saveErr := createAuth(uid, ts) // Redis 토큰 메타데이터 저장
	if saveErr != nil {
		return map[string]string{}, err
	}
	tokens := map[string]string{
		"access_token":  ts.AccessToken,
		"refresh_token": ts.RefreshToken,
	}
	return tokens, nil
}
func logoutUser(c *gin.Context) error {
	// request header 에 담긴 access & refresh token을 검증 후 redis 에서 삭제
	au, ru, err := ExtractTokenMetadata(c.Request)
	if err != nil {
		return err
	}
	deleted, delErr := DeleteAuth(au.AccessUuid, ru.RefreshUuid)
	if delErr != nil || deleted == 0 {
		return err
	}
	return nil

}
func modifyUser(c *gin.Context) error {
	return nil
}
func getPostList(c *gin.Context) ([]post, error) {
	db := storage.DB()
	query := `select count(*) from post`

	var length int
	_ = db.QueryRow(query).Scan(&length)
	if length == 0 {
		return []post{}, errors.New("Nothing to show")
	}
	query = `select post.PID,post.UId,Title,TotalMember,FE,BE,AOS,IOS,PM,Designer,More from post join member on post.pid = member.pid`
	rows, err := db.Query(query)
	if err != nil {
		return []post{}, err
	}
	defer rows.Close()

	posts := make([]post, 0)
	var pos post
	for rows.Next() {
		err := rows.Scan(&pos.PostID, &pos.UID, &pos.Title,
			&pos.TotalMember, &pos.FE, &pos.BE, &pos.AOS, &pos.IOS,
			&pos.PM, &pos.Designer, &pos.More)
		if err != nil {
			return []post{}, err
		}
		posts = append(posts, pos)
	}
	return posts, nil
}
func getPostDetail(c *gin.Context) (post, error) {
	postDetail := post{}

	return postDetail, nil
}
func addPost(c *gin.Context) error {

	var reqBody struct {
		UID           uint   `json:"uid"`
		Title         string `json:"title"`
		Desc          string `json:"desc"`
		TotMem        uint   `json:"totalMember"`
		FE            uint   `json:"fe"`
		BE            uint   `json:"be"`
		AOS           uint   `json:"aos"`
		IOS           uint   `json:"ios"`
		PM            uint   `json:"pm"`
		Designer      uint   `json:"designer"`
		More          uint   `json:"more"`
		FE_desc       string `json:"fe_desc"`
		BE_desc       string `json:"be_desc"`
		AOS_desc      string `json:"aos_desc"`
		IOS_desc      string `json:"ios_desc"`
		PM_desc       string `json:"pm_desc"`
		Designer_desc string `json:"designer_desc"`
		More_desc     string `json:"more_desc"`
	}
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		return err
	}

	db := storage.DB()
	_, err := db.Exec(`Insert into post(UID,Title,TotalMember) values(?,?,?)`, reqBody.UID, reqBody.Title, reqBody.TotMem)
	if err != nil {
		return err
	}
	var pid uint
	db.QueryRow(`select pid from post order by pid desc limit 1`).Scan(&pid)
	_, err = db.Exec(`Insert into member values(?,?,?,0,?,?,0,?,?,0,?,?,0,?,?,0,?,?,0,?,?,0)`, pid,
		reqBody.FE, reqBody.FE_desc, reqBody.BE, reqBody.BE_desc, reqBody.AOS, reqBody.AOS_desc,
		reqBody.IOS, reqBody.IOS_desc, reqBody.PM, reqBody.PM_desc, reqBody.Designer,
		reqBody.Designer_desc, reqBody.More, reqBody.More_desc)
	if err != nil {
		return err
	}
	db.Prepare(`Insert into member values(?)`)
	_, err = db.Exec(`Insert into postDetail values(?,?)`, pid, reqBody.Desc)
	if err != nil {
		return err
	}
	return nil
}

type user struct {
	UID   uint   `json:"uid"`
	ID    string `json:"id"`
	PW    string `json:"pw"`
	Name  string `json:"name"`
	Email string `json:"email"`
	PR    string `json:"pr"`
}
type post struct {
	PostID      uint   `json:"pid"`
	UID         uint   `json:"uid"`
	Title       string `json:"title"`
	Description string `json:"desc"`
	TotalMember uint   `json:"totalMember"`
	FE          uint   `json:"fe"`
	BE          uint   `json:"be"`
	AOS         uint   `json:"aos"`
	IOS         uint   `json:"ios"`
	PM          uint   `json:"pm"`
	Designer    uint   `json:"designer"`
	More        uint   `json:"more"`
}

type member struct {
	PostID        uint
	FE            uint
	BE            uint
	AOS           uint
	IOS           uint
	PM            uint
	Designer      uint
	More          uint
	FE_desc       string
	BE_desc       string
	AOS_desc      string
	IOS_desc      string
	PM_desc       string
	Designer_desc string
	More_desc     string
}
