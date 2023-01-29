package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

	beego.GlobalControllerRouter["github.com/csuhan/csugo/controllers:BusController"] = append(beego.GlobalControllerRouter["github.com/csuhan/csugo/controllers:BusController"],
		beego.ControllerComments{
			Method:           "Search",
			Router:           "/bus/search/:start/:end/:time",
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/csuhan/csugo/controllers:CetController"] = append(beego.GlobalControllerRouter["github.com/csuhan/csugo/controllers:CetController"],
		beego.ControllerComments{
			Method:           "GetHGrade",
			Router:           "/cet/hgrade/:id/:name",
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/csuhan/csugo/controllers:CetController"] = append(beego.GlobalControllerRouter["github.com/csuhan/csugo/controllers:CetController"],
		beego.ControllerComments{
			Method:           "GetZKZ",
			Router:           "/cet/zkz/:id/:type",
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/csuhan/csugo/controllers:ClassRoomController"] = append(beego.GlobalControllerRouter["github.com/csuhan/csugo/controllers:ClassRoomController"],
		beego.ControllerComments{
			Method:           "GetJXL",
			Router:           "/classroom/jxl/:xq",
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/csuhan/csugo/controllers:ClassRoomController"] = append(beego.GlobalControllerRouter["github.com/csuhan/csugo/controllers:ClassRoomController"],
		beego.ControllerComments{
			Method:           "GetJXLS",
			Router:           "/classroom/jxls",
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/csuhan/csugo/controllers:ClassRoomController"] = append(beego.GlobalControllerRouter["github.com/csuhan/csugo/controllers:ClassRoomController"],
		beego.ControllerComments{
			Method:           "GetFreeWeekTime",
			Router:           "/classroom/time/:date/:xq/:jxl",
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/csuhan/csugo/controllers:JobController"] = append(beego.GlobalControllerRouter["github.com/csuhan/csugo/controllers:JobController"],
		beego.ControllerComments{
			Method:           "List",
			Router:           "/job/:typeid/:pageindex/:pagesize/:hastime",
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/csuhan/csugo/controllers:JwcController"] = append(beego.GlobalControllerRouter["github.com/csuhan/csugo/controllers:JwcController"],
		beego.ControllerComments{
			Method:           "Class",
			Router:           "/jwc/:id/:pwd/class/:term/:week",
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/csuhan/csugo/controllers:JwcController"] = append(beego.GlobalControllerRouter["github.com/csuhan/csugo/controllers:JwcController"],
		beego.ControllerComments{
			Method:           "Exam",
			Router:           "/jwc/:id/:pwd/exam",
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/csuhan/csugo/controllers:JwcController"] = append(beego.GlobalControllerRouter["github.com/csuhan/csugo/controllers:JwcController"],
		beego.ControllerComments{
			Method:           "Grade",
			Router:           "/jwc/:id/:pwd/grade",
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/csuhan/csugo/controllers:JwcController"] = append(beego.GlobalControllerRouter["github.com/csuhan/csugo/controllers:JwcController"],
		beego.ControllerComments{
			Method:           "Rank",
			Router:           "/jwc/:id/:pwd/rank",
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})
	beego.GlobalControllerRouter["github.com/csuhan/csugo/controllers:JwcController"] = append(beego.GlobalControllerRouter["github.com/csuhan/csugo/controllers:JwcController"],
		beego.ControllerComments{
			Method:           "GPA",
			Router:           "/jwc/:id/:pwd/gpa",
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})
	beego.GlobalControllerRouter["github.com/csuhan/csugo/controllers:JwcController"] = append(beego.GlobalControllerRouter["github.com/csuhan/csugo/controllers:JwcController"],
		beego.ControllerComments{
			Method:           "WeekList",
			Router:           "/jwc/:id/:pwd/weeklist/:term",
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/csuhan/csugo/controllers:LibController"] = append(beego.GlobalControllerRouter["github.com/csuhan/csugo/controllers:LibController"],
		beego.ControllerComments{
			Method:           "List",
			Router:           "/lib/list/:id/:pwd",
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/csuhan/csugo/controllers:LibController"] = append(beego.GlobalControllerRouter["github.com/csuhan/csugo/controllers:LibController"],
		beego.ControllerComments{
			Method:           "Login",
			Router:           "/lib/login/:id/:pwd",
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/csuhan/csugo/controllers:LibController"] = append(beego.GlobalControllerRouter["github.com/csuhan/csugo/controllers:LibController"],
		beego.ControllerComments{
			Method:           "Reloan",
			Router:           "/lib/reloan/:id/:pwd/:books",
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/csuhan/csugo/controllers:LibController"] = append(beego.GlobalControllerRouter["github.com/csuhan/csugo/controllers:LibController"],
		beego.ControllerComments{
			Method:           "Search",
			Router:           "/lib/search/:keyword[get]",
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/csuhan/csugo/controllers:NewsController"] = append(beego.GlobalControllerRouter["github.com/csuhan/csugo/controllers:NewsController"],
		beego.ControllerComments{
			Method:           "GetNewsContent",
			Router:           "/news/article/:link",
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/csuhan/csugo/controllers:NewsController"] = append(beego.GlobalControllerRouter["github.com/csuhan/csugo/controllers:NewsController"],
		beego.ControllerComments{
			Method:           "GetNewsList",
			Router:           "/news/:id/:pwd/list/:pageid",
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/csuhan/csugo/controllers:StudentController"] = append(beego.GlobalControllerRouter["github.com/csuhan/csugo/controllers:StudentController"],
		beego.ControllerComments{
			Method:           "GetStudentInfo",
			Router:           "/student/:id",
			AllowHTTPMethods: []string{"get"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

	beego.GlobalControllerRouter["github.com/csuhan/csugo/controllers:WxUserController"] = append(beego.GlobalControllerRouter["github.com/csuhan/csugo/controllers:WxUserController"],
		beego.ControllerComments{
			Method:           "Login",
			Router:           "/login",
			AllowHTTPMethods: []string{"post"},
			MethodParams:     param.Make(),
			Filters:          nil,
			Params:           nil})

}
