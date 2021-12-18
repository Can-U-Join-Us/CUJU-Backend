package api

type ResgisterForm struct {
	Email string `json:"email"`
	PW    string `json:"pw"`
	Name  string `json:"name"`
	Phone string `json:"phone"`
}
type LoginForm struct {
	ID string `json:"Email"`
	PW string `json:"pw"`
}
type ModifyForm struct {
	UID int    `json:"uid"`
	PW  string `json:"pw"`
	NEW string `json:"new"`
}
type Project struct {
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
type AddProjectForm struct {
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
type JoinForm struct {
	PID      int    `json:"pid"`
	UID      int    `json:"uid"`
	CATEGORY string `json:"category"`
}
type ReplyJoinForm struct {
	PID int `json:"pid"`
	UID int `json:"uid"`
}
type Announce struct {
	TITLE   string `json:"title"`
	CONTENT string `json:"content"`
}
type Msg struct {
	TYPE    int    `json:"type"`
	TITLE   string `json:"title"`
	CONTENT string `json:"content"`
	PID     int    `json:"pid"`
	UID     int    `json:"uid"`
}
