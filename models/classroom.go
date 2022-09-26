package models

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/boltdb/bolt"
	"strconv"
	"io/ioutil"
	"net/http"
)

var ClassTime = []string{"0102","0304","0506","0708","0910"}
var FreeTime = []bool{false,false,false,false,false}


const API_URL = "http://csujwc.its.csu.edu.cn/app.do"

type UserAuth struct {
	Token string `json:"token"`
}
type XQ struct {
	ID   string `json:"xqid"`
	Name string `json:"xqmc"`
}


type JXL struct {
	ID   string `json:"jzwid"`
	Name string `json:"jzwmc"`
	XQ   XQ     `json:"XQ"`
}

type ClassRoom struct {
	JSID         string `json:"jsid"`
	RoomName  string `json:"jsmc"`
	ClassTime int
	FreeWeekTime []bool
}

type RoomList struct {
	JSList []ClassRoom `json:"jsList"`
}

var XQS = []XQ{
	{ID: "1", Name: "校本部"}, {ID: "2", Name: "南校区"}, {ID: "3", Name: "铁道校区"},
	{ID: "4", Name: "湘雅新校区"}, {ID: "5", Name: "湘雅老校区"}, {ID: "6", Name: "湘雅医院"},
	{ID: "7", Name: "湘雅二医院"}, {ID: "8", Name: "湘雅三医院"}, {ID: "9", Name: "新校区"},
}

func getDB() (*bolt.DB, error) {
	dbName := beego.AppConfig.String("DB::ClassesDB")
	db, err := bolt.Open(dbName, 777, nil)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func GetFreeWeekTime(Date, XQ, JXL string) ([]ClassRoom, error) {
	token := GetToken()
	roomlist := []RoomList{}

	jsList := [5][]ClassRoom{}
	jsList2 := []ClassRoom{}
	//classrooms := []ClassRoom{}

	//获取全天空闲教室信息
	for jc,time := range ClassTime{

	FREE_ROOM_API_URL := API_URL + "?method=getKxJscx&time=" +Date+ "&xqid=" +XQ+ "&jxlid=" +JXL+"&idleTime=" +time
	req, _:= http.NewRequest("GET", FREE_ROOM_API_URL, nil)
	req.Header.Add("token", token)

	response, _ := http.DefaultClient.Do(req)
	body, _ := ioutil.ReadAll(response.Body)
	defer response.Body.Close()


	err := json.Unmarshal([]byte(body), &roomlist)
	if err != nil{

	}

	jsList[jc] = roomlist[0].JSList

	for i , _ := range jsList[jc]{

		jsList[jc][i].ClassTime = jc + 1
		jsList[jc][i].FreeWeekTime = FreeTime

	}

	beego.Info(jsList[jc])
	xqid, err := strconv.Atoi(XQ)
	for i , _ := range roomlist{
		//classrooms[i].FreeWeekTime
		beego.Info(i,xqid)

	}

		jsList2 = append(jsList2, jsList[jc]...)

	}
	values := jsList2
	//冒泡排序
	for i := 0; i < len(values)-1; i++ {
		for j := i+1; j < len(values); j++ {
			if  values[i].JSID>values[j].JSID{
				values[i],values[j] = values[j],values[i]
			}
		}
	}
	//合并空闲信息
	left, right := 0, 1
	for ; right < len(values); right++ {

		if values[left].JSID == values[right].JSID {
			FreeTime[values[right].ClassTime - 1] = true
			continue
		}
		left++
		values[left-1].FreeWeekTime = FreeTime
		values[left-1].FreeWeekTime[values[left-1].ClassTime - 1] = true
		FreeTime = []bool{false,false,false,false,false}
		values[left] = values[right]
	}

	return values, nil
}

//根据校区获取教学楼
func GetBuildingsByXQ(XQ string) ([]JXL, error) {
	token := GetToken()

	JXL_API_URL := API_URL + "?method=getJxlcx&xqid=" +XQ
	req, _:= http.NewRequest("GET", JXL_API_URL, nil)
	req.Header.Add("token", token)

	response, _ := http.DefaultClient.Do(req)
	body, _ := ioutil.ReadAll(response.Body)
	beego.Info(string(body))

	defer response.Body.Close()

	jxls := []JXL{}

	err := json.Unmarshal([]byte(body), &jxls)
	if err != nil{
		return jxls, err
	}
	xqid, err := strconv.Atoi(XQ)
	for i , _ := range jxls{
		jxls[i].XQ = XQS[xqid - 1]

	}

	return jxls, nil
}

//获取所有校区
func GetXQS() []XQ {
	token := GetToken()

	XQ_API_URL := API_URL + "?method=getXqcx"
	req, _:= http.NewRequest("GET", XQ_API_URL, nil)
	req.Header.Add("token", token)
	//beego.Info(session.Body)
	response, _ := http.DefaultClient.Do(req)
	body, _ := ioutil.ReadAll(response.Body)
	beego.Info(string(body))

	defer response.Body.Close()

	xqs := []XQ{}
	err := json.Unmarshal([]byte(body), &xqs)
	if err != nil{

	}

	return xqs
}

func GetToken() string{

	LOGIN_API_URL := API_URL + "?method=authUser&xh=" + "8212190323"+ "&pwd=" +"jwc456258"

	session, _:= http.NewRequest("GET", LOGIN_API_URL, nil)
	session.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	//beego.Info(session.Body)
	response, _ := http.DefaultClient.Do(session)
	body, _ := ioutil.ReadAll(response.Body)

	defer response.Body.Close()
	var user UserAuth
	err := json.Unmarshal([]byte(string(body)), &user)
	if err != nil {

	}

	beego.Info(user.Token)
	return user.Token
}
