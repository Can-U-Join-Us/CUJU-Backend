package api

import (
	"errors"
	"math/rand"
	"net/smtp"
	"strconv"
	"time"

	ErrChecker "github.com/Can-U-Join-Us/CUJU-Backend/modules/errors"
	"github.com/Can-U-Join-Us/CUJU-Backend/modules/storage"
	"github.com/Can-U-Join-Us/CUJU-Backend/modules/token"
	"github.com/gin-gonic/gin"
)

func RegisterUser(c *gin.Context) error {
	var reqBody ResgisterForm
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
	_, err = db.Exec("Insert into user(email,PW,Name,Phone) values(?,?,?,?)", reqBody.Email, reqBody.PW, reqBody.Name, reqBody.Phone)
	if err := ErrChecker.Check(err); err != nil {
		return err
	}
	return nil
}

func LoginUser(c *gin.Context) (uint64, map[string]string, error) {
	var reqBody LoginForm
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
	ts, err := token.CreateToken(uid)
	if err := ErrChecker.Check(err); err != nil {
		return 0, map[string]string{}, err
	}
	err = token.CreateAuth(uid, ts) // Redis 토큰 메타데이터 저장
	if err := ErrChecker.Check(err); err != nil {
		return 0, map[string]string{}, err
	}
	tokens := map[string]string{
		"access_token":  ts.AccessToken,
		"refresh_token": ts.RefreshToken,
	}
	return uid, tokens, nil
}
func LogoutUser(c *gin.Context) error {
	// request header 에 담긴 access & refresh token을 검증 후 redis 에서 삭제
	au, ru, err := token.ExtractBothTokenMetadata(c.Request)
	if err := ErrChecker.Check(err); err != nil {
		return err
	}
	deleted, err := token.DeleteAuth(au.AccessUuid, ru.RefreshUuid)
	if err := ErrChecker.Check(err); err != nil || deleted == 0 {
		return err
	}
	return nil
}
func FindUserPW(c *gin.Context) error {
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
func FindUserId(c *gin.Context) (string, error) {
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
func ModifyPW(c *gin.Context) error {
	var reqBody ModifyForm
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
	_, err = db.Exec(`update user set pw ="` + reqBody.NEW + `" where uid = ` + uid)
	if err := ErrChecker.Check(err); err != nil {
		return err
	}
	return nil
}
func ModifyProfile(c *gin.Context) error {

	return nil
}
