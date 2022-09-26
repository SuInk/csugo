package models

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"github.com/PuerkitoBio/goquery"
	"github.com/astaxie/beego"
	"github.com/csuhan/csugo/utils"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
	"strings"
)

const JWC_UNIFIED_URL = "https://ca.csu.edu.cn/authserver/login?service=http%3A%2F%2Fcsujwc.its.csu.edu.cn%2Fsso.jsp"
const JWC_BASE_URL = "http://csujwc.its.csu.edu.cn/jsxsd/"
const JWC_LOGIN_URL = JWC_BASE_URL + "xk/LoginToXk"
const JWC_GRADE_URL = JWC_BASE_URL + "kscj/yscjcx_list"
const JWC_RANK_URL = JWC_BASE_URL + "kscj/zybm_cx"
const JWC_CLASS_URL = JWC_BASE_URL + "xskb/xskb_list.do"
const JWC_EXAMOPTION_URL = JWC_BASE_URL + "xsks/xsksap_query"
const JWC_EXAM_URL = JWC_BASE_URL + "xsks/xsksap_list"

type JwcUser struct {
	Id, Pwd, Name, College, Margin, Class string
}

type JwcGrade struct {
	ClassNo int
	FirstTerm, GottenTerm, ClassName,
	MiddleGrade, FinalGrade, Grade,
	ClassScore, ClassType, ClassProp string
}

type Rank struct {
	Term, TotalScore, ClassRank, AverScore string
}

type Exam struct {
	Term, ExamState, Round, Coden, Name, Time, Classroom, Seat, Others string
}

type JwcRank struct {
	User  JwcUser
	Ranks []Rank
}
type GPA struct {
	Algorithm string
	GPA       float64
}

type Class struct {
	ClassName, Teacher, Weeks, Place string
}

type Weeklist struct {
	WeekList string
}

type Jwc struct{}

//成绩查询
func (this *Jwc) Grade(user *JwcUser) ([]JwcGrade, error) {
	//登录系统
	cookies, err := this.Login(user)
	if err != nil {
		beego.Debug(err)
		return nil, err
	}
	response, err := this.LogedRequest(user, "GET", JWC_GRADE_URL, cookies, nil)
	if err != nil {
		return []JwcGrade{}, err
	}
	data, _ := ioutil.ReadAll(response.Body)
	//beego.Info(data)
	defer response.Body.Close()
	if !strings.Contains(string(data), "学生个人考试成绩") {
		return []JwcGrade{}, utils.ERROR_JWC
	}
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(data)))
	if err != nil {
		return []JwcGrade{}, utils.ERROR_SERVER
	}
	Grades := []JwcGrade{}
	doc.Find("table#dataList tr").Each(func(i int, selection *goquery.Selection) {
		if i != 0 {
			s := selection.Find("td")
			jwcgrade := JwcGrade{
				ClassNo:     i,
				FirstTerm:   s.Eq(1).Text(),
				GottenTerm:  s.Eq(2).Text(),
				ClassName:   s.Eq(3).Text(),
				MiddleGrade: s.Eq(4).Text(),
				FinalGrade:  s.Eq(5).Text(),
				Grade:       s.Eq(6).Text(),
				ClassScore:  s.Eq(7).Text(),
				ClassType:   s.Eq(8).Text(),
				ClassProp:   s.Eq(9).Text(),
			}
			Grades = append(Grades, jwcgrade)
		}
	})
	return Grades, nil
}

//专业排名查询
func (this *Jwc) Rank(user *JwcUser) ([]Rank, error) {
	//登录系统
	cookies, err := this.Login(user)
	if err != nil {
		beego.Debug(err)
		return nil, err
	}
	response, err := this.LogedRequest(user, "POST", JWC_RANK_URL, cookies, strings.NewReader("xqfw="+url.QueryEscape("入学以来")))
	if err != nil {
		return []Rank{}, err
	}
	data, _ := ioutil.ReadAll(response.Body)
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(data)))
	if err != nil {
		return []Rank{}, utils.ERROR_SERVER
	}
	terms := make([]string, 0)
	doc.Find("#xqfw option").Each(func(i int, s *goquery.Selection) {
		terms = append(terms, s.Text())
	})
	err = nil
	ranks := make([]Rank, len(terms))
	ch := make(chan map[int]Rank)
	chanRanks := []map[int]Rank{}
	for key, term := range terms {
		go func(key int, term string, ch chan map[int]Rank) {
			resp, _ := this.LogedRequest(user, "POST", JWC_RANK_URL, cookies, strings.NewReader("xqfw="+url.QueryEscape(term)))
			data, _ := ioutil.ReadAll(resp.Body)
			resp.Body.Close()
			doc, _ := goquery.NewDocumentFromReader(strings.NewReader(string(data)))
			td := doc.Find("#dataList tr").Eq(1).Find("td")
			rank := Rank{
				Term:       term,
				TotalScore: td.Eq(1).Text(),
				ClassRank:  td.Eq(2).Text(),
				AverScore:  td.Eq(3).Text(),
			}
			ch <- map[int]Rank{key: rank}
		}(key, term, ch)
	}
	for range terms {
		chanRanks = append(chanRanks, <-ch)
	}
	for i := 0; i < len(terms); i++ {
		for j := 0; j < len(chanRanks); j++ {
			if v, ok := chanRanks[j][i]; ok {
				ranks[i] = v
			}
		}
	}
	return ranks, err
}

//课表查询
func (this *Jwc) Class(user *JwcUser, Week, Term string) ([][]Class, string, error) {
	if Week == "0" {
		Week = ""
	}
	body := strings.NewReader("zc=" + url.QueryEscape(Week) + "&xnxq01id=" + url.QueryEscape(Term) + "&sfFD=1")
	//登录系统
	cookies, err := this.Login(user)
	if err != nil {
		beego.Debug(err)
		return [][]Class{}, "", err
	}
	response, err := this.LogedRequest(user, "POST", JWC_CLASS_URL, cookies, body)
	if err != nil {
		return [][]Class{}, "", err
	}
	data, _ := ioutil.ReadAll(response.Body)
	defer response.Body.Close()
	classes := make([][]Class, 0)
	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(string(data)))
	doc.Find("table#kbtable").Eq(0).Find("td div.kbcontent").Each(func(i int, s *goquery.Selection) {
		font := s.Find("font")
		if font.Size() == 3 || font.Size() == 2 { //一节课
			class := Class{
				ClassName: s.Nodes[0].FirstChild.Data,
				Teacher:   font.Eq(0).Text(),
				Weeks:     font.Eq(1).Text(),
				Place:     font.Eq(2).Text(),
			}
			classes = append(classes, []Class{class})
		} else if font.Size() == 6 || font.Size() == 5 || font.Size() == 4 { //两节课
			class := []Class{
				Class{
					ClassName: s.Nodes[0].FirstChild.Data,
					Teacher:   font.Eq(0).Text(),
					Weeks:     font.Eq(1).Text(),
					Place:     font.Eq(2).Text(),
				},
				Class{
					ClassName: font.Eq(3).Nodes[0].PrevSibling.PrevSibling.Data,
					Teacher:   font.Eq(3).Text(),
					Weeks:     font.Eq(4).Text(),
					Place:     font.Eq(5).Text(),
				},
			}
			classes = append(classes, class)
		} else {
			classes = append(classes, make([]Class, 1))
		}
	})

	classes = classes[1:]
	classes = append(classes, classes[40])

	//每学期开学时间
	temp := doc.Find("table#kbtable").Eq(1).Find("td").Eq(0).Text()
	startWeekDay := temp

	return classes, startWeekDay, nil
}

//考试查询
func (this *Jwc) Exam(user *JwcUser) ([]Exam, error) {
	//登录系统
	cookies, err := this.Login(user)
	if err != nil {
		beego.Debug(err)
		return nil, err
	}
	response, err := this.LogedRequest(user, "POST", JWC_EXAMOPTION_URL, cookies, strings.NewReader("xqfw="+url.QueryEscape("")))
	if err != nil {
		return []Exam{}, err
	}
	data, _ := ioutil.ReadAll(response.Body)
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(data)))
	if err != nil {
		return []Exam{}, utils.ERROR_SERVER
	}
	terms := make([]string, 0)
	doc.Find("#xnxqid option").Each(func(i int, s *goquery.Selection) {
		terms = append(terms, s.Text())
	})
	//beego.Info(terms)//打印cookie
	err = nil
	exams := []Exam{}
	for _, term := range terms {
		resp, _ := this.LogedRequest(user, "POST", JWC_EXAM_URL, cookies, strings.NewReader("xnxqid="+url.QueryEscape(term)))
		data, _ := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		doc, _ := goquery.NewDocumentFromReader(strings.NewReader(string(data)))
		doc.Find("#dataList tr").Each(func(i int, s *goquery.Selection) {
			tds := s.Find("td")
			exams = append(exams, Exam{
				Term:      term,
				ExamState: tds.Eq(1).Text(),
				Round:     tds.Eq(2).Text(),
				Coden:     tds.Eq(3).Text(),
				Name:      tds.Eq(4).Text(),
				Time:      tds.Eq(5).Text(),
				Classroom: tds.Eq(6).Text(),
				Seat:      tds.Eq(7).Text() + "号",
				Others:    tds.Eq(8).Text(),
			})

			//beego.Info(exams)
		})

	}
	exam_len := len(exams)
	var ret []Exam
	for i := 0; i < exam_len; i++ {
		if exams[i].Name == "" {
			continue
		}
		ret = append(ret, exams[i])
	}

	return ret, err
}

//教学周历查询
func (this *Jwc) WeekList(user *JwcUser, Term string) ([]Weeklist, error) {
	body := strings.NewReader("zc=" + url.QueryEscape("") + "&xnxq01id=" + url.QueryEscape(Term) + "&sfFD=1")
	//登录系统
	cookies, err := this.Login(user)
	if err != nil {
		beego.Debug(err)
		return nil, err
	}
	response, _ := this.LogedRequest(user, "POST", JWC_CLASS_URL, cookies, body)
	WeekLists := make([]Weeklist, 0)

	data, _ := ioutil.ReadAll(response.Body)
	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(string(data)))
	td2 := doc.Find("table#kbtable").Eq(1).Find("td")

	//temp :=td2.Eq(0).Text()
	WeekList := Weeklist{
		WeekList: td2.Eq(0).Text(),
	}
	for i := 0; i < 22; i++ {
		WeekList = Weeklist{
			WeekList: td2.Eq(i).Text(),
		}
		WeekLists = append(WeekLists, WeekList)
	}

	return WeekLists, nil
}
func ToScore(grade string) string {
	switch {
	case grade == "优":
		grade = "95"
	case grade == "良":
		grade = "85"
	case grade == "中":
		grade = "75"
	case grade == "及格" || grade == "合格":
		grade = "65"
	case grade == "不及格" || grade == "不合格":
		grade = "0"
	}
	return grade
}
func (this *Jwc) ConvertToGPA(user *JwcUser) ([]GPA, error) {
	grades, _ := this.Grade(user)
	var standard4, standard4_1 float64
	var credits = 0.0
	GPA := make([]GPA, len(grades))
	for _, v := range grades {
		credit, _ := strconv.ParseFloat(v.ClassScore, 64)
		credits += credit
		standard4 = standard4 + Standard4(v.Grade)*credit
		standard4_1 = standard4_1 + Standard4_1(v.Grade)*credit

	}
	GPA[0].GPA = standard4 / credits
	GPA[1].GPA = standard4_1 / credits
	return GPA, nil

}
func Standard4(grade string) float64 {
	GPA := 0.0
	grade = ToScore(grade)
	switch {
	case grade >= "90":
		GPA = 4.0
	case grade >= "80":
		GPA = 3.0
	case grade >= "70":
		GPA = 2.0
	case grade >= "60":
		GPA = 1.0
	default:
		GPA = 4.0
	}
	return GPA
}
func Standard4_1(grade string) float64 {
	GPA := 0.0
	grade = ToScore(grade)
	switch {
	case grade == "优" || grade >= "85":
		GPA = 4.0
	case grade == "良" || grade >= "70":
		GPA = 3.0
	case grade == "中" || grade >= "60" || grade == "及格" || grade == "合格":
		GPA = 2.0
	default:
		GPA = 0.0
	}
	return GPA
}

//登录后请求
func (this *Jwc) LogedRequest(user *JwcUser, Method, Url string, client http.Client, Params io.Reader) (*http.Response, error) {

	//beego.Info(cookies)//打印cookies
	//查询
	Req, err := http.NewRequest(Method, Url, Params)
	Req.Header.Add("content-type", "application/x-www-form-urlencoded")
	if err != nil {
		return nil, utils.ERROR_SERVER
	}
	// beego.Info(Req)
	Cookiesreq, err := client.Do(Req)
	// beego.Info(Cookiesreq)

	return Cookiesreq, err
}

//随机字符串
func GetRandomString(n int) []byte {
	str := "ABCDEFGHJKMNPQRSTWXYZabcdefhijkmnprstwxyz2345678"
	bytes := []byte(str)
	var result []byte
	for i := 0; i < n; i++ {
		result = append(result, bytes[rand.Intn(len(bytes))])
	}
	return result
}

//对明文进行填充
func Padding(plainText []byte, blockSize int) []byte {
	//计算要填充的长度
	n := blockSize - len(plainText)%blockSize
	//对原来的明文填充n个n
	temp := bytes.Repeat([]byte{byte(n)}, n)
	plainText = append(plainText, temp...)
	return plainText
}

//AEC加密（CBC模式）
func AES_CBC_Encrypt(plainText []byte, key []byte) string {
	//指定加密算法，返回一个AES算法的Block接口对象
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	//进行填充
	plainText = append(GetRandomString(64), plainText...)
	plainText = Padding(plainText, block.BlockSize())
	//指定初始向量vi,长度和block的块尺寸一致
	iv := GetRandomString(16)
	//beego.Info("key=" + string(key))
	//beego.Info("vi=" + string(iv))
	//指定分组模式，返回一个BlockMode接口对象
	blockMode := cipher.NewCBCEncrypter(block, iv)
	//加密连续数据库
	cipherText := make([]byte, len(plainText))
	blockMode.CryptBlocks(cipherText, plainText)
	//返回密文
	return base64.StdEncoding.EncodeToString(cipherText)
}

//教务系统登录
func (this *Jwc) Login(user *JwcUser) (http.Client, error) {
	password, _ := base64.StdEncoding.DecodeString(user.Pwd)
	//获取cookie
	var client http.Client
	jar, err := cookiejar.New(nil)
	if err != nil {
		panic(err)
	}
	client.Jar = jar
	req, _ := http.NewRequest("GET", JWC_UNIFIED_URL, nil)
	// req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.198 Safari/537.36")
	response, err := client.Do(req)
	//response, err := http.Get(JWC_UNIFIED_URL)
	if err != nil {
		return client, utils.ERROR_SERVER
	}
	//body1, _ := ioutil.ReadAll(response.Body)
	//beego.Info(string(body1))
	doc, err := goquery.NewDocumentFromReader(response.Body)
	//beego.Info(doc.Find("#pwdEncryptSalt").AttrOr("value", ""))
	encoedepwd := AES_CBC_Encrypt([]byte(password), []byte(doc.Find("#pwdEncryptSalt").AttrOr("value", "")))

	//beego.Info(doc.Find("#execution").AttrOr("value", ""))
	reqData := url.Values{
		"username":   {user.Id},
		"password":   {encoedepwd},
		"captcha":    {"None"},
		"rememberMe": {"True"},
		"_eventId":   {"submit"},
		"cllt":       {"userNameLogin"},
		"dllt":       {"generalLogin"},
		"lt":         {"None"},
		"execution":  {doc.Find("#execution").AttrOr("value", "")},
	}
	//beego.Info(reqData)

	req, _ = http.NewRequest("POST", JWC_UNIFIED_URL, strings.NewReader(reqData.Encode()))
	req.Header.Add("content-type", "application/x-www-form-urlencoded")
	// req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.198 Safari/537.36")
	response1, err := client.Do(req)
	//beego.Info(response1.Cookies())
	if err != nil {
		return client, utils.ERROR_SERVER
	}
	body, _ := ioutil.ReadAll(response1.Body)
	defer response.Body.Close()
	//登陆成功
	if strings.Contains(string(body), "我的桌面") {
		return client, nil
	}
	//账号或密码错误
	return client, utils.ERROR_ID_PWD
}
