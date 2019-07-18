package rules

/*
 * @Author: li.xiangyang
 * @Date: 2019-07-15 09:57:56
 * @Last Modified by: li.xiangyang
 * @Last Modified time: 2019-07-15 09:59:00
 * 第9影视爬虫规则
 */

import (
	"github.com/henrylee2cn/pholcus/app/downloader/request" //必需
	spider "github.com/henrylee2cn/pholcus/app/spider"      //必需
	"github.com/henrylee2cn/pholcus/common/goquery"         //DOM解析
	//信息输出
)

func init() {
	d9ys.Register()
}

const baseURL = "http://www.d9ys.com"

var d9ys = &spider.Spider{
	Name:        "d9ys",
	Description: "第9影视 [www.d9ys.com]",
	// Pausetime: 300,
	// Keyin:   KEYIN,
	// Limit:        LIMIT,
	EnableCookie: false,
	RuleTree: &spider.RuleTree{
		Root: func(ctx *spider.Context) {
			ctx.AddQueue(&request.Request{Url: baseURL, Rule: "videoTypes"})
		},

		Trunk: map[string]*spider.Rule{
			"videoTypes": {
				ItemFields: []string{
					"title",
					"url",
				},
				ParseFunc: func(ctx *spider.Context) {
					query := ctx.GetDom()
					lis := query.Find(".menu2 a")
					lis.Each(func(i int, s *goquery.Selection) {
						if i == 0 {
							return
						}
						if url, ok := s.Attr("href"); ok {
							title := s.Text()
							url = baseURL + url
							types := map[int]interface{}{
								0: title,
								1: url,
							}
							// fmt.Println(types)
							ctx.Output(types)
							ctx.AddQueue(&request.Request{Url: url, Rule: "videos", Temp: map[string]interface{}{"VideoType": title}})
						}
					})
				},
			},
			"videos": {
				ItemFields: []string{
					"title",
					"actor",
					"statusText",
					"area",
					"year",
					"kind",
					"url",
				},
				ParseFunc: func(ctx *spider.Context) {
					query := ctx.GetDom()

					//下一页
					pageo := query.Find(".page em")
					if pageo != nil {
						pageo = pageo.Last()
						no := pageo.Next()
						if no != nil {
							// nextPage := no.Text()
							url, _ := no.Attr("href")
							url = baseURL + url
							ctx.AddQueue(&request.Request{
								Url:  url,
								Rule: "videos",
								Temp: map[string]interface{}{"VideoType": ctx.GetTemp("VideoType", "")}})
						}
					}
					lis := query.Find(".mlist li")
					kind := ctx.GetTemp("VideoType", "")
					lis.Each(func(i int, s *goquery.Selection) {
						to := s.Find(".info h2 a")
						ps := s.Find(".info p")
						title := to.Text()            //视频名称
						url, _ := to.Attr("href")     //视频播放页地址
						actor := ps.Eq(0).Text()      //主演
						statusText := ps.Eq(1).Text() //状态
						area := ps.Eq(2).Text()       //地区
						year := ps.Eq(3).Text()       //年代

						detail := map[int]interface{}{
							0: title,
							1: actor,
							2: statusText,
							3: area,
							4: year,
							5: kind,
							6: baseURL + url,
						}
						// fmt.Println(detail)
						ctx.Output(detail)
					})
				},
			},
		},
	},
}
