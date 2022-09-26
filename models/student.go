package models

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/httplib"
	"github.com/csuhan/csugo/utils"
	"regexp"
	"strings"
)

const EVERYCLASS_URL = "https://everyclass.xyz/query?id="

type Student struct {
	Name, Deputy, School, Class string
}

type StudentList struct {
	ID, Name, TotalNum string
	Students     []Student
}


func GetMulStudentInfo(StudentID string) (StudentList, error) {
	req := httplib.Post(EVERYCLASS_URL + StudentID)

	resp, err := req.String()
	if err != nil {
		return StudentList{}, utils.ERROR_SERVER
	}
	students := StudentList{}
	//查找总页数,总信息数
	re := regexp.MustCompile("中南那么大，居然有(.*)个(.*)！</h1> ")
	res := re.FindStringSubmatch(resp)
	students.ID = StudentID
	if len(res) != 0 {
		students.TotalNum = res[1]
		students.Name = res[2]
	}
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(resp))
	if err != nil {
		return StudentList{}, utils.ERROR_SERVER
	}
	studentItems := []Student{}
	//查找每个tr
	doc.Find("tbody tr").Each(func(i int, s *goquery.Selection) {
		tds := s.Find("td")
		//temp := tds.Find("a").AttrOr("onclick", "")
		//link := regexp.MustCompile(`/Home/Release_TZTG_zd/(.*)', '', 'left=0`).FindStringSubmatch(temp)[1]
		studentItems = append(studentItems, Student{
			Name:       strings.Trim(tds.Eq(0).Text(), "\n "),
			Deputy:     strings.Trim(tds.Eq(1).Text(), "\n "),
			School:      strings.Trim(tds.Eq(2).Text(), "\n "),
			Class: strings.Trim(tds.Eq(3).Text(), "\n "),
			//Time:      strings.Trim(tds.Eq(6).Text(), "\n "),
			//Link:      link,
		})
	})
	students.Students = studentItems
	return students, nil
}

func GetSinStudentInfo(StudentID string) (StudentList, error) {
	req := httplib.Post(EVERYCLASS_URL + StudentID)

	resp, err := req.String()
	if err != nil {
		return StudentList{}, utils.ERROR_SERVER
	}
	students := StudentList{}
	//查找总页数,总信息数
	re := regexp.MustCompile("<h1 class=hero-header>(.*)</h1> ")
	res := re.FindStringSubmatch(resp)
	students.ID = StudentID
	students.TotalNum = "1"
	if len(res) != 0 {
		students.Name = res[1]
	}
	re2 := regexp.MustCompile("- (.*) - 每课")
	res2 := re2.FindStringSubmatch(resp)
	beego.Info(res2)
	students.ID = StudentID
	students.TotalNum = "1"
	if res2 == nil {
		res2 = append(res2, "0","查无此人")
		students.TotalNum = "0"
		return StudentList{}, utils.ERROR_NO_STUDENT
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(resp))
	if err != nil {
		return StudentList{}, utils.ERROR_SERVER
	}
	studentItems := []Student{}
	//查找每个tr
	doc.Find("div").Each(func(i int, s *goquery.Selection) {
		h1s := s.Find("h1")
		h4s := s.Find("h4")
		//temp := tds.Find("a").AttrOr("onclick", "")
		//link := regexp.MustCompile(`/Home/Release_TZTG_zd/(.*)', '', 'left=0`).FindStringSubmatch(temp)[1]
		studentItems = append(studentItems, Student{
			Name:       strings.Trim(h1s.Eq(0).Text(), "\n "),
			Deputy:     strings.Trim(res2[1], "\n "),
			School:      strings.Trim(h4s.Eq(0).Text(), "\n "),
			Class: "",
			//Time:      strings.Trim(jwc/8212190323/jwc456258/gradetds.Eq(6).Text(), "\n "),
			//Link:      link,
		})
		//对数组切片，取第一个元素
		studentItems = studentItems[:1]
	})
	students.Students = studentItems
	return students, err
}


