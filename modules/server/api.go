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
func loginUser(c *gin.Context) (uint, error) {
	var reqBody struct {
		ID string `json:"id"`
		PW string `json:"pw"`
	}
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		return 0, err
	}
	db := storage.DB()
	query := `select uid, PW from user where ID = "` + reqBody.ID + `"`
	var pw string
	row := db.QueryRow(query)
	var uid uint
	err := row.Scan(&uid, &pw)
	if err != nil { // ID 가 없으면 ID 가 없다는 오류 반환
		return 0, errors.New("잘못된 ID")
	}
	if reqBody.PW != pw { // PW 가 다르면 PW 가 다르다는 오류 반환
		return 0, errors.New("PW 불일치")
	}
	return uid, nil
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

/*
create table post (
	PID int primary key auto_increment,
	UID int not null,
	Title varchar(50) not null,
	Description varchar(255) not null,
	Link varchar(255),
	foreign key(UID) references user(UID)
);
create table image (
	PID int primary key,
	ImageLoc varchar(255),
	foreign key(PID) references post(PID)
);
create table member (
	PID int primary key,
	FE int default 0, FE_desc varchar(100), FE_join int check(FE_join <= FE),
	BE int default 0, BE_desc varchar(100), BE_join int check(BE_join <= BE),
	AOS int default 0, AOS_desc varchar(100), AOS_join int check(AOS_join <= AOS),
	IOS int default 0, IOS_desc varchar(100), IOS_join int check(IOS_join <= IOS),
	PM int default 0, PM_desc varchar(100), PM_join int check(PM_join <= PM),
	Designer int default 0, Designer_desc varchar(100), Designer_join int check(Designer_join <= Designer),
	More int default 0, More_desc varchar(100), More_join int check(More_join <= More),
	foreign key(PID) references post(PID)
	);
*/
