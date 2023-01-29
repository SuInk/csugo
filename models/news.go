package models

import (
	"encoding/json"
	"github.com/PuerkitoBio/goquery"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/httplib"
	"github.com/csuhan/csugo/utils"
	"io"
	"net/http"
	"strings"
)

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

const NewsContentUrl = "http://tz.its.csu.edu.cn/Home/Release_TZTG_zd/"
const NewsUnifiedLoginUrl = "https://oa.csu.edu.cn/con/ggtz"
const NewsListUrl = "https://oa.csu.edu.cn/con/xnbg/contentList"

type NewsItem struct {
	ID, Title, Dept, Time string
	Link, Content         string
	ViewCount             int
}

type NewsList struct {
	NowPage              string
	TotalPage, TotalNews int
	News                 []NewsItem
}

func GetNewsList(user *JwcUser, PageID string) (NewsList, error) {
	cookies, err := UnifiedLogin(user, NewsUnifiedLoginUrl)
	cookie := strings.Split(cookies, ";")
	cookies = cookie[2]
	beego.Info(cookies)
	if err != nil {
		return NewsList{}, err
	}
	req, _ := http.NewRequest("POST", NewsListUrl, strings.NewReader("params=%7B%22tableName%22%3A%22ZNDX_ZHBG_GGTZ%22%2C%22tjnr%22%3A%22%22%7D&pageSize=1&pageNo=20"))
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Cookie", cookies)
	resp, err := http.DefaultClient.Do(req)
	body, _ := io.ReadAll(resp.Body)
	beego.Info(string(body))
	var newsListJson NewsListJson
	err = json.Unmarshal(body, &newsListJson)
	if err != nil {
		return NewsList{}, err
	}

	var newsItems []NewsItem
	for _, data := range newsListJson.Data {
		newsItems = append(newsItems, NewsItem{
			ID:        data.DJSJP,
			Title:     data.WJBT,
			Dept:      data.QCBMMC,
			ViewCount: data.LLCS,
			Time:      data.QCSJ,
			Link:      data.YWMC,
		})
	}
	var news NewsList
	news.News = newsItems
	news.NowPage = PageID
	news.TotalNews = newsListJson.Count
	news.TotalPage = newsListJson.Count / 20
	return news, nil
}

func GetNewsContent(link string) (string, error) {
	req := httplib.Get(NewsContentUrl + link)
	req.Header("x-forwarded-for", "202.197.71.84") //模仿校内登录
	resp, err := req.String()
	if err != nil {
		return "", utils.ERROR_SERVER
	}
	res, err := htmldeparse(resp)
	if err != nil {
		return "", utils.ERROR_SERVER
	}
	return res, nil
}

func htmldeparse(resp string) (string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(resp))
	if err != nil {
		return "", utils.ERROR_SERVER
	}
	//找到文章区
	docContent := doc.Find("table").Eq(2).Find("tr").Eq(2).Find("td").Eq(0)
	//内容处理,去除多余内容
	docContent.Find("p.MsoNormal").Each(func(i int, s *goquery.Selection) {
		s.SetAttr("style", "text-indent: 32px;")
		temp := strings.Trim(s.Text(), "\u00a0")
		if temp == "" {
			s.Remove()
		} else {
			s.SetHtml(temp)
		}
	})
	res, err := docContent.Html()
	res = "<div style='margin:20px 10px;font-size:16px!important;'>" + res + "</div>"
	//o:p标签,特殊字符去除
	spestrs := []string{"<o:p></o:p>", "<o:p>", "</o:p>"}
	for _, spestr := range spestrs {
		res = strings.Replace(res, spestr, "", -1)
	}
	return res, nil
}
