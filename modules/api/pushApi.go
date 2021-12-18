package api

import (
	"errors"
	"fmt"
	"strconv"

	ErrChecker "github.com/Can-U-Join-Us/CUJU-Backend/modules/errors"
	"github.com/Can-U-Join-Us/CUJU-Backend/modules/storage"
	"github.com/gin-gonic/gin"
)

func RefreshMsg(c *gin.Context) ([]Msg, error) {
	uid := c.Request.Header.Get("uid")
	fmt.Println("uid : ", uid)
	db := storage.DB()
	var pid, id, result int
	var name, email, title string
	query := `select p.pid, j.uid,u.name, u.email, j.result from join_queue j join project_post p on j.pid = p.pid and p.uid = ` + uid + ` left join user u on u.uid = j.uid where j.result = 0 and j.confirm = 0`
	rows, err := db.Query(query)
	if err != nil {
		return []Msg{}, err
	}
	defer rows.Close()
	msgList := make([]Msg, 0)
	res, _ := getAnnouncement(c)
	for i := 0; i < len(res); i++ {
		msgList = append(msgList, res[i])
	}
	var m Msg
	for rows.Next() {
		if err := rows.Scan(&pid, &id, &name, &email, &result); err != nil {
			return []Msg{}, err
		}
		uuid, _ := strconv.Atoi(uid)
		m = Msg{
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
		return []Msg{}, err
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&id, &title, &result); err != nil {
			return []Msg{}, err
		}
		var ctt string
		if result == 1 { // permission
			ctt = title + `(` + strconv.Itoa(pid) + `) 프로젝트 : 승인`
		} else { // deny
			ctt = title + `(` + strconv.Itoa(pid) + `) 프로젝트 : 거절`
		}
		m = Msg{
			TYPE:    2,
			TITLE:   strconv.Itoa(id) + ` 번 프로젝트 주최자의 답변이 도착했습니다 !`,
			CONTENT: ctt,
			PID:     0,
			UID:     0,
		}

		msgList = append(msgList, m)
	}

	if len(msgList) == 0 {
		return []Msg{}, errors.New("Nothing")
	}
	query = `update join_queue j join project_post p on j.pid= p.pid and p.uid = 7 left join user u on u.uid = j.uid  set j.confirm = 1 where j.result > 0`
	_, err = db.Exec(query) // 참여 결과 알림 메세지는 읽음 상태로 변경
	if err != nil {
		return []Msg{}, err
	}
	return msgList, nil
}
func getAnnouncement(c *gin.Context) ([]Msg, error) {
	annoList := make([]Msg, 0)
	db := storage.DB()
	query := `select title,content from announce where send = false`
	rows, err := db.Query(query)
	if err != nil {
		return []Msg{}, errors.New("Nothing")
	}
	defer rows.Close()

	for rows.Next() {
		var title, content string
		if err := rows.Scan(&title, &content); err != nil {
			return []Msg{}, err
		}
		m := Msg{
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
func Announcement(c *gin.Context) error {
	var reqBody Announce
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
