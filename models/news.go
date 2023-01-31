package models

import (
	"bytes"
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/ledongthuc/pdf"
	"io"
	"net/http"
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

func GetNewsList(user *JwcUser, PageID string) (NewsList, error) {
	cookie, err := UnifiedLogin(user, NewsUnifiedLoginUrl)
	cookies := strings.Split(cookie, ";")
	cookie = cookies[2]
	// beego.Info(cookie)
	if err != nil {
		return NewsList{}, err
	}
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

func GetNewsContent(link, cookie string) (string, error) {

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
	f, r, err := pdf.Open("./news/" + link + ".pdf")
	if err != nil {
		beego.Info(err)
		return "", err
	}
	// remember close file
	defer f.Close()
	var buf bytes.Buffer
	b, err := r.GetPlainText()
	if err != nil {
		beego.Info(err)
		return "", err
	}
	buf.ReadFrom(b)
	return buf.String(), nil
}
