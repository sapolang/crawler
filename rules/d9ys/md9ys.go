package d9ys

/*
 * @Author: li.xiangyang
 * @Date: 2019-07-19 09:10:29
 * @Last Modified by: li.xiangyang
 * @Last Modified time: 2019-07-19 10:07:30
 * 第9影视爬虫规则m站
 */

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/henrylee2cn/pholcus/common/goquery"

	//必需

	"github.com/henrylee2cn/pholcus/app/downloader/request"
	spider "github.com/henrylee2cn/pholcus/app/spider" //必需
	//DOM解析
	//信息输出
)

func init() {
	md9ys.Register()
}

var md9ys = &spider.Spider{
	Name:        "md9ys",
	Description: "第9影视 [m.d9ys.com]",
	// Pausetime: 300,
	// Keyin:   KEYIN,
	// Limit:        LIMIT,
	EnableCookie: false,
	RuleTree: &spider.RuleTree{
		Root: func(ctx *spider.Context) {
			// ctx.AddQueue(&request.Request{Url: "http://m.d9ys.com/qvod/70694/", Rule: "fetchDetail", Temp: map[string]interface{}{"tag": "movie", "videoType": "dz"}})

			baseURL := "http://m.d9ys.com"
			types := map[string][]string{
				"movie": []string{"dz", "xj", "aq", "kh", "jq", "kb"},
				"tv":    []string{"gc", "rh", "om", "gt"},
				"other": []string{"man", "dh", "zy"},
			}
			for k, v := range types {
				for _, item := range v {
					ctx.AddQueue(&request.Request{Url: baseURL + "/" + item, Rule: "fetchList", Temp: map[string]interface{}{"tag": k, "videoType": item, "isStart": 1}})
				}
			}
		},

		Trunk: map[string]*spider.Rule{
			"fetchList": { //拉取列表页
				AidFunc: func(ctx *spider.Context, v map[string]interface{}) interface{} {
					maxPageNum := v["maxPageNum"].(int)
					videoType := v["videoType"].(string)
					tag := v["tag"].(string)
					baseURL := "http://m.d9ys.com"
					fmt.Println("fet...", tag, videoType, maxPageNum)

					for index := 2; index < maxPageNum; index++ {
						//"/dz/142.html"
						url := baseURL + fmt.Sprintf("/%s/%d.html", videoType, index)
						// fmt.Println(index, url)

						ctx.AddQueue(&request.Request{Url: url, Rule: "fetchList", Temp: map[string]interface{}{
							"videoType": videoType,
							"tag":       tag,
						}})
					}
					return nil
				},
				ParseFunc: func(ctx *spider.Context) {
					query := ctx.GetDom()
					items := query.Find(".stui-vodlist__item a")
					videoType := ctx.GetTemp("videoType", "")
					tag := ctx.GetTemp("tag", "")
					isStart := ctx.GetTemp("isStart", 0).(int)
					if isStart == 1 {
						lastpage := query.Find(".page a").Last()
						lastPageHref, _ := lastpage.Attr("href")
						reg := regexp.MustCompile(`\d+`)
						matches := reg.FindAllString(lastPageHref, -1)
						if len(matches) > 0 {
							maxPageNum, _ := strconv.Atoi(matches[0])
							ctx.Aid(map[string]interface{}{
								"maxPageNum": maxPageNum,
								"videoType":  videoType,
								"tag":        tag,
							}, "fetchList")
						}
					}

					items.Each(func(i int, s *goquery.Selection) {
						href, _ := s.Attr("href")
						url := "http://m.d9ys.com" + href
						// fmt.Println(url)
						ctx.AddQueue(&request.Request{Url: url, Rule: "fetchDetail", Temp: map[string]interface{}{
							"videoType": videoType,
							"tag":       tag,
						}})
					})
				},
			},
			"fetchDetail": { //拉取详情页信息
				ItemFields: []string{
					"title",
					"actor",
					"cover",
					"area",
					"year",
					"url",
					"playPageURL",
					"tag",
					"videoType",
					"short",
				},
				ParseFunc: func(ctx *spider.Context) {
					query := ctx.GetDom()
					title := query.Find(".stui-content__detail h3").Text()
					actor := query.Find(".stui-content__detail .data").Eq(2).Text()
					tempText := strings.Split(query.Find(".stui-content__detail .data").Eq(1).Text(), "   ")
					area := tempText[0]
					year := tempText[1]
					cover, _ := query.Find(".stui-vodlist__thumb").Attr("data-original")
					playPageURL, _ := query.Find("#m_html_p4 a").Attr("href")
					videoType := ctx.GetTemp("videoType", "")
					tag := ctx.GetTemp("tag", "")
					short := query.Find(".short").Text()
					detail := map[int]interface{}{
						0: title,
						1: actor,
						2: cover,
						3: area,
						4: year,
						5: ctx.Request.Url,
						6: playPageURL,
						7: tag,
						8: videoType,
						9: short,
					}
					// fmt.Println(detail)
					ctx.Output(detail)
				},
			},
		},
	},
}
