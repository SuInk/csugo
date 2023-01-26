package models

import (
	"encoding/json"
	"github.com/PuerkitoBio/goquery"
	"github.com/astaxie/beego"
	"github.com/csuhan/csugo/utils"
	"io"
	"net/http"
	"strings"
)

// CHEER_TIMETABLE_URL 琦课 网址: https://cheer-timetable.vercel.app/
const CHEER_TIMETABLE_URL = "http://47.114.89.18:3000/search/"

type RawJson struct {
	Props struct {
		PageProps struct {
			Name string `json:"name"`
			Data [][]struct {
				Id             string      `json:"id"`
				Seq            json.Number `json:"seq"`
				Grade          string      `json:"grade"`
				Name           string      `json:"name"`
				Faculty        string      `json:"facultyName"`
				ProfessionName string      `json:"professionName"`
				ClassName      string      `json:"className"`
				Sex            string      `json:"sex"`
				CreatedAt      string      `json:"createdAt"`
				UpdatedAt      string      `json:"updatedAt"`
			} `json:"data"`
		} `json:"pageProps"`
		__N_SSG bool `json:"__N_SSG"`
	} `json:"props"`
	Page  string `json:"page"`
	Query struct {
		Name string `json:"name"`
	} `json:"query"`
	BuildId      string   `json:"buildId"`
	IsFallback   bool     `json:"isFallback"`
	Gsp          bool     `json:"gsp"`
	ScriptLoader []string `json:"scriptLoader"`
}
type Student struct {
	Name,
	Deputy,
	School,
	Class string
}

type StudentList struct {
	ID, Name, TotalNum string
	Students           []Student
}

func GetStudentInfo(StudentID string) (StudentList, error) {
	resp, err := http.Get(CHEER_TIMETABLE_URL + StudentID)
	if err != nil {
		return StudentList{}, utils.ERROR_SERVER
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	// beego.Info(string(body))
	students := StudentList{}
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(body)))
	if err != nil {
		return StudentList{}, utils.ERROR_SERVER
	}
	var rawJson RawJson
	// Using goquery.Single, only the first match is selected
	singleSel := doc.FindMatcher(goquery.Single("script#__NEXT_DATA__"))
	// beego.Info(singleSel.Text())
	err = json.Unmarshal([]byte(singleSel.Text()), &rawJson)
	// beego.Info(rawJson)
	if err != nil {
		beego.Info(err)
		return StudentList{}, err
	}
	if rawJson.Props.PageProps.Data[0] == nil && rawJson.Props.PageProps.Data[1] == nil {
		return StudentList{}, utils.ERROR_STUDENT_NOT_FOUND
	}
	// 学生
	for _, item := range rawJson.Props.PageProps.Data[0] {
		var student Student
		student.Name = item.Name
		student.Deputy = "学生"
		student.School = item.Faculty
		student.Class = item.ClassName
		students.Students = append(students.Students, student)
	}
	// 教师
	for _, item := range rawJson.Props.PageProps.Data[1] {
		var teacher Student
		teacher.Name = item.Name
		teacher.Deputy = "教师"
		teacher.School = item.Faculty
		teacher.Class = item.Faculty
		students.Students = append(students.Students, teacher)
	}
	students.ID = StudentID
	students.Name = rawJson.Props.PageProps.Name
	students.TotalNum = string(len(students.Students))

	return students, nil
}
