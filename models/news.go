package models

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/astaxie/beego"
	"github.com/csuhan/csugo/utils"
	"github.com/ledongthuc/pdf"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

const NewsContentUrl = "https://oa.csu.edu.cn/con/xnbg/loadContentPdf/"
const NewsUnifiedLoginUrl = "https://oa.csu.edu.cn/con/ggtz"
const NewsListUrl = "https://oa.csu.edu.cn/con/xnbg/contentList"

type NewsItem struct {
	ID                int
	Title, Dept, Time string
	Link, Content     string
	ViewCount         int
}

type NewsList struct {
	NowPage, Cookie      string
	TotalPage, TotalNews int
	News                 []NewsItem
}

type NewsListJson struct {
	Data []struct {
		DJSJP  string `json:"DJSJP"`
		JLNM   string `json:"JLNM"`
		LLCS   int    `json:"LLCS"`
		QCBMMC string `json:"QCBMMC"`
		PXBM   int    `json:"PXBM"`
		QCZHXM string `json:"QCZHXM"`
		WN     int    `json:"WN"`
		QCSJ   string `json:"QCSJ"`
		WJBT   string `json:"WJBT"`
		YWMC   string `json:"YWMC"`
		YWMS   string `json:"YWMS"`
		DJSJ   string `json:"DJSJ"`
	} `json:"data"`
	Count int `json:"count"`
}

// UnifiedLogin 统一认证登录的封装,返回cookie
func UnifiedLogin(user *JwcUser, unifiedUrl string) (string, error) {
	//尝试免登录
	if client, ok := ClientMap[*user]; ok {
		req, _ := http.NewRequest("GET", unifiedUrl, nil)
		resp, _ := client.Do(req)
		body, _ := io.ReadAll(resp.Body)
		if strings.Contains(string(body), "校内通知") {
			var cookies string
			for _, v := range resp.Cookies() {
				cookies += v.String()
			}
			if cookies == "" {
				cookies = client.Jar.Cookies(req.URL)[0].String()
			}
			return "empty;empty;" + cookies, nil
		}
	}
	password, _ := base64.StdEncoding.DecodeString(user.Pwd)
	//获取cookie
	var client http.Client
	jar, err := cookiejar.New(nil)
	if err != nil {
		panic(err)
	}
	client.Jar = jar
	resp, _ := client.Get(unifiedUrl)
	nowUrl := resp.Request.URL.String()
	beego.Info(nowUrl)
	// 校内网免登录，会出错
	req, _ := http.NewRequest("GET", nowUrl, nil)
	req.Header.Add("User-Agent", "csulite robot v1.0")
	response, err := client.Do(req)
	//response, err := http.Get(JWC_UNIFIED_URL)
	if err != nil || response.StatusCode != 200 {
		return "", utils.ERROR_UNIFIED
	}
	//body1, _ := ioutil.ReadAll(response.Body)
	//beego.Info(string(body1))
	doc, err := goquery.NewDocumentFromReader(response.Body)
	//beego.Info(doc.Find("#pwdEncryptSalt").AttrOr("value", ""))
	encodePwd := AES_CBC_Encrypt([]byte(password), []byte(doc.Find("#pwdEncryptSalt").AttrOr("value", "")))
	// 验证码识别
	captcha := "None"
	respIsNeed, err := client.Get(fmt.Sprintf("https://ca.csu.edu.cn/authserver/checkNeedCaptcha.htl?username=%s&_=%s", user.Id, strconv.FormatInt(time.Now().UnixNano()/1e6, 10)))
	if err != nil {
		return "", err
	}
	body, _ := io.ReadAll(respIsNeed.Body)
	if strings.Contains(string(body), "true") {
		//需要验证码
		log.Println(user.Id, "需要验证码")
		captcha, err = utils.GetCaptcha(&client)
		if err != nil {
			return "", err
		}
	}
	reqData := url.Values{
		"username":   {user.Id},
		"password":   {encodePwd},
		"captcha":    {captcha},
		"rememberMe": {"True"},
		"_eventId":   {"submit"},
		"cllt":       {"userNameLogin"},
		"dllt":       {"generalLogin"},
		"lt":         {"None"},
		"execution":  {doc.Find("#execution").AttrOr("value", "")},
	}
	response, err = client.Post(nowUrl, "application/x-www-form-urlencoded", strings.NewReader(reqData.Encode()))
	//body, _ = io.ReadAll(response.Body)
	//log.Println(string(body))
	// 统一认证错误处理
	doc, err = goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		return "", utils.ErrorServer
	}
	if strings.Contains(doc.Text(), "中南e行APP扫码登录") && response.StatusCode != 200 {
		switch doc.Find("span#showErrorTip").First().Text() {
		case "验证码错误":
			return "", utils.ErrorCaptcha
		case "您提供的用户名或者密码有误":
			return "", utils.ErrorIdPwd
		case "输入多次密码错误账号冻结，5-10分钟自动解冻":
			return "", utils.ErrorLocked
		default:
			return "", utils.ErrorFailLogin
		}
	}
	ClientMap[*user] = client
	return req.Header.Get("Cookie"), nil
}

// GetNewsList 获取新闻列表
func GetNewsList(user *JwcUser, PageID string) (NewsList, error) {
	cookie, err := UnifiedLogin(user, NewsUnifiedLoginUrl)
	if err != nil {
		return NewsList{}, err
	}
	cookies := strings.Split(cookie, ";")
	cookie = cookies[2]
	req, _ := http.NewRequest("POST", NewsListUrl, strings.NewReader("params=%7B%22tableName%22%3A%22ZNDX_ZHBG_GGTZ%22%2C%22tjnr%22%3A%22%22%7D&pageSize="+PageID+"&pageNo=20"))
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Cookie", cookie)
	resp, err := http.DefaultClient.Do(req)
	body, _ := io.ReadAll(resp.Body)
	// beego.Info(string(body))
	var newsListJson NewsListJson
	err = json.Unmarshal(body, &newsListJson)
	if err != nil {
		return NewsList{}, err
	}

	var newsItems []NewsItem
	page, err := strconv.Atoi(PageID)
	for i, data := range newsListJson.Data {
		formattedTime, _ := time.Parse("Jan 2, 2006 3:04:05 PM", data.DJSJP)
		newsItems = append(newsItems, NewsItem{
			ID:        page*20 - 19 + i,                          //
			Title:     data.WJBT,                                 // 文件标题
			Dept:      data.QCBMMC,                               // 起草部门名称
			ViewCount: data.LLCS,                                 // 浏览次数
			Time:      formattedTime.Format("2006-01-02 15:04 "), // 登记时间
			Link:      data.JLNM,
		})
	}
	var news NewsList
	news.News = newsItems
	news.NowPage = PageID
	news.Cookie = cookie[strings.Index(cookie, "=")+1:]
	news.TotalNews = newsListJson.Count
	news.TotalPage = newsListJson.Count / 20
	return news, nil
}

func GetNewsContent(link, cookie string) ([]string, error) {

	if f, _, err := pdf.Open("./news/" + link + ".pdf"); err != nil {
		f.Close()
		file, err := os.Create("./news/" + link + ".pdf")
		beego.Info(link)
		req, _ := http.NewRequest("GET", NewsContentUrl+link, nil)
		req.Header.Add("Cookie", cookie)
		resp, err := http.DefaultClient.Do(req)
		_, err = io.Copy(file, resp.Body)
		if err != nil {
			beego.Info(err)
		}
	}
	return pdfParser(link)
}
func pdfParser(link string) ([]string, error) {
	f, r, err := pdf.Open("./news/" + link + ".pdf")
	// remember close file
	defer f.Close()
	if err != nil {
		return []string{}, err
	}
	totalPage := r.NumPage()
	var fullText []string
	for pageIndex := 1; pageIndex <= totalPage; pageIndex++ {
		p := r.Page(pageIndex)
		if p.V.IsNull() {
			continue
		}
		var lastTextStyle pdf.Text
		texts := p.Content().Text
		for _, text := range texts {
			if isSameSentence(lastTextStyle, text) {
				lastTextStyle.X = text.X
				lastTextStyle.Y = text.Y
				lastTextStyle.Font = text.Font
				lastTextStyle.FontSize = text.FontSize
				lastTextStyle.S = lastTextStyle.S + text.S
			} else {
				// fmt.Printf("Font: %s, Font-size: %f, x: %f, y: %f, content: %s \n", lastTextStyle.Font, lastTextStyle.FontSize, lastTextStyle.X, lastTextStyle.Y, lastTextStyle.S)
				fullText = append(fullText, lastTextStyle.S)
				lastTextStyle = text
			}
		}
		fullText = append(fullText, lastTextStyle.S)
	}
	for i, text := range fullText {
		if text == "" && i == 0 {
			fullText = append(fullText[:i], fullText[i+1:]...)
		}
		if text == "" && i != 0 && i < len(fullText)-2 {
			fullText[i-1] += fullText[i+1]
			fullText = append(fullText[:i], fullText[i+2:]...)
		}
	}
	return fullText, nil
}
func isSameSentence(text1, text2 pdf.Text) bool {
	if text2.FontSize-text1.FontSize > 2 || text2.FontSize-text1.FontSize < -2 {
		return false
	}
	if (text2.Y-text1.Y < 10 && text2.Y-text1.Y > -10) || text2.X < 100 || text2.FontSize > 20 {
		return true
	}
	return false
}
