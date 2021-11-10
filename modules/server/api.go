package server

import (
	"errors"
	"fmt"

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
	fmt.Println(reqBody)
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
	fmt.Println("User add complete")

	return nil
}
func loginUser(c *gin.Context) error {
	var reqBody struct {
		ID string `json:"id"`
		PW string `json:"pw"`
	}
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		return err
	}
	fmt.Println(reqBody)
	db := storage.DB()
	query := `select PW from user where ID = "` + reqBody.ID + `"`
	var pw string
	fmt.Println(query)
	row := db.QueryRow(query)
	err := row.Scan(&pw)
	if err != nil { // ID 가 없으면 ID 가 없다는 오류 반환

		return errors.New("잘못된 ID")
	}
	if reqBody.PW != pw { // PW 가 다르면 PW 가 다르다는 오류 반환

		return errors.New("PW 불일치")
	}
	fmt.Println("User add complete")

	return nil
}
func modifyUser(c *gin.Context) error {
	return nil
}
func getPostList(c *gin.Context) ([]post, error) {
	db := storage.DB()
	query := `select * from post`
	rows, err := db.Query(query)
	fmt.Println("\n\n", rows)
	if err != nil {
		return []post{}, err
	}
	return []post{}, nil
}
func getPostDetail(c *gin.Context) error {

	return nil
}
func addPost(c *gin.Context) error {

	var reqBody post
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		return err
	}
	db := storage.DB()
	stmt, err := db.Prepare("Insert into Post values(?,?,?,?,?,?,?)")
	if err != nil {
		return err
	}
	rs, err := stmt.Exec(reqBody.PostID, reqBody.PartyType, reqBody.ItemType, reqBody.TotalPrice, reqBody.Title, reqBody.Description, reqBody.Link)
	if err != nil {
		return err
	}
	_ = rs
	defer stmt.Close()
	fmt.Println("Post add complete")
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
	PostID      uint     `json:"pid"`
	PartyType   uint     `json:"party_type"`
	ItemType    uint     `json:"item_type"`
	TotalPrice  uint     `json:"total_price"`
	Title       string   `json:"title"`
	Description string   `json:"desc"`
	Link        string   `json:"link"`
	ImageLoc    []string `json:"image_list"`
}
