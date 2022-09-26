package controllers

import (
	"github.com/astaxie/beego"
	"github.com/csuhan/csugo/models"
)

type ClassRoomController struct {
	beego.Controller
}

// @router /classroom/time/:date/:xq/:jxl [get]
func (this *ClassRoomController) GetFreeWeekTime() {
	date := this.Ctx.Input.Param(":date")
	xq := this.Ctx.Input.Param(":xq")
	jxl := this.Ctx.Input.Param(":jxl")

	cls, err := models.GetFreeWeekTime(date, xq, jxl)
	stateCode := 1
	errorstr := ""
	if err != nil {
		stateCode = -1
		errorstr = err.Error()
	}
	this.Data["json"] = struct {
		StateCode int
		Error     string
		CLS       []models.ClassRoom
	}{
		StateCode: stateCode,
		Error:     errorstr,
		CLS:       cls,
	}
	this.ServeJSON()

}

// @router /classroom/jxl/:xq [get]
func (this *ClassRoomController) GetJXL() {
	xq := this.Ctx.Input.Param(":xq")
	jxls, err := models.GetBuildingsByXQ(xq)
	stateCode := 1
	errorstr := ""
	if err != nil {
		stateCode = -1
		errorstr = err.Error()
	}
	this.Data["json"] = struct {
		StateCode int
		Error     string
		JXLS      []models.JXL
	}{
		StateCode: stateCode,
		Error:     errorstr,
		JXLS:      jxls,
	}
	this.ServeJSON()
}

// @router /classroom/jxls [get]
func (this *ClassRoomController) GetJXLS() {
	xqs := models.GetXQS()
	this.Data["json"] = struct {
		StateCode int
		Error     string
		XQS      []models.XQ
	}{
		StateCode: 1,
		Error:     "",
		XQS:      xqs,
	}
	this.ServeJSON()
}
