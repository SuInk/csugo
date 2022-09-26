package controllers

import (
	"github.com/astaxie/beego"
	"github.com/csuhan/csugo/models"
)

type StudentController struct {
	beego.Controller
}

// @router /student/:id
func (this *StudentController) GetStudentInfo() {
	pageid := this.Ctx.Input.Param(":id")

	students, err := models.GetMulStudentInfo(pageid)
	if students.TotalNum == "" {

		students, err = models.GetSinStudentInfo(pageid)

	}
	stateCode := 1
	errorstr := ""
	if err != nil {
		stateCode = -1
		errorstr = err.Error()
	}
	this.Data["json"] = struct {
		StateCode int
		Error     string
		Students      models.StudentList
	}{
		StateCode: stateCode,
		Error:     errorstr,
		Students:  students,
	}
	this.ServeJSON()
}

