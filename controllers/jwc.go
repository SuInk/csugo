package controllers

import (
	//"github.com/PuerkitoBio/goquery"
	"github.com/astaxie/beego"
	"github.com/csuhan/csugo/models"
	//"io/ioutil"
	//"net/url"
	//"strings"
)

type JwcController struct {
	beego.Controller
}

// @router /jwc/:id/:pwd/grade [get]
func (this *JwcController) Grade() {
	user := &models.JwcUser{
		Id:  this.Ctx.Input.Param(":id"),
		Pwd: this.Ctx.Input.Param(":pwd")}
	jwc := &models.Jwc{}
	grade, err := jwc.Grade(user)
	stateCode := 1
	errorstr := ""
	if err != nil {
		stateCode = -1
		errorstr = err.Error()
	}
	this.Data["json"] = struct {
		StateCode int
		Error     string
		Grades    []models.JwcGrade
	}{
		StateCode: stateCode,
		Error:     errorstr,
		Grades:    grade,
	}
	this.ServeJSON()
}

// @router /jwc/:id/:pwd/grade [get]
func (this *JwcController) GPA() {
	user := &models.JwcUser{
		Id:  this.Ctx.Input.Param(":id"),
		Pwd: this.Ctx.Input.Param(":pwd")}
	jwc := &models.Jwc{}
	GPA, err := jwc.ConvertToGPA(user)
	stateCode := 1
	errorstr := ""
	if err != nil {
		stateCode = -1
		errorstr = err.Error()
	}
	this.Data["json"] = struct {
		StateCode int
		Error     string
		Grades    []models.GPA
	}{
		StateCode: stateCode,
		Error:     errorstr,
		Grades:    GPA,
	}
	this.ServeJSON()
}

// @router /jwc/:id/:pwd/rank [get]
func (this *JwcController) Rank() {
	user := &models.JwcUser{
		Id:  this.Ctx.Input.Param(":id"),
		Pwd: this.Ctx.Input.Param(":pwd")}
	jwc := &models.Jwc{}
	rank, err := jwc.Rank(user)
	stateCode := 1
	errorstr := ""
	if err != nil {
		stateCode = -1
		errorstr = err.Error()
	}
	this.Data["json"] = struct {
		StateCode int
		Error     string
		Rank      []models.Rank
	}{
		StateCode: stateCode,
		Error:     errorstr,
		Rank:      rank,
	}
	this.ServeJSON()
}

// @router /jwc/:id/:pwd/exam [get]
func (this *JwcController) Exam() {
	user := &models.JwcUser{
		Id:  this.Ctx.Input.Param(":id"),
		Pwd: this.Ctx.Input.Param(":pwd")}
	jwc := &models.Jwc{}
	exam, err := jwc.Exam(user)
	stateCode := 1
	errorstr := ""
	if err != nil {
		stateCode = -1
		errorstr = err.Error()
	}
	this.Data["json"] = struct {
		StateCode int
		Error     string
		Exam      []models.Exam
	}{
		StateCode: stateCode,
		Error:     errorstr,
		Exam:      exam,
	}
	this.ServeJSON()
}

// @router /jwc/:id/:pwd/class/:term/:week [get]
func (this *JwcController) Class() {
	user := &models.JwcUser{
		Id:  this.Ctx.Input.Param(":id"),
		Pwd: this.Ctx.Input.Param(":pwd")}
	week := this.Ctx.Input.Param(":week")
	term := this.Ctx.Input.Param(":term")
	jwc := &models.Jwc{}
	class, startWeekDay, err := jwc.Class(user, week, term)
	stateCode := 1
	errorstr := ""
	if err != nil {
		stateCode = -1
		errorstr = err.Error()
	}
	this.Data["json"] = struct {
		StateCode    int
		Error        string
		Class        [][]models.Class
		StartWeekDay string
	}{
		StateCode:    stateCode,
		Error:        errorstr,
		Class:        class,
		StartWeekDay: startWeekDay,
	}
	this.ServeJSON()
}

//@router /jwc/:id/:pwd/weeklist/:term [get]
func (this *JwcController) WeekList() {
	user := &models.JwcUser{
		Id:  this.Ctx.Input.Param(":id"),
		Pwd: this.Ctx.Input.Param(":pwd")}
	Term := this.Ctx.Input.Param(":term")
	jwc := &models.Jwc{}
	weeklists, err := jwc.WeekList(user, Term)
	stateCode := 1
	errorstr := ""
	if err != nil {
		stateCode = -1
		errorstr = err.Error()
	}

	this.Data["json"] = struct {
		StateCode int
		Error     string
		WeekList  []models.Weeklist
	}{
		StateCode: stateCode,
		Error:     errorstr,
		WeekList:  weeklists,
	}
	this.ServeJSON()

}
