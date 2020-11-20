package gushiwen

import (
	"fmt"
	"strings"

	"github.com/henrylee2cn/pholcus/app/downloader/request"
	"github.com/henrylee2cn/pholcus/app/spider/common"
	"github.com/henrylee2cn/pholcus/common/goquery" //DOM解析

	. "github.com/henrylee2cn/pholcus/app/spider" //必需
)

// func init() {
// 	Gushiwen.Register()
// }

//Gushiwen 古诗文爬虫
var Gushiwen = &Spider{
	Name:         "古诗文爬虫",
	Description:  "古诗文爬虫 [Auto Page] [https://www.gushiwen.cn/]",
	EnableCookie: true,
	RuleTree: &RuleTree{
		Root: func(ctx *Context) {
			ctx.AddQueue(&request.Request{Url: "https://www.gushiwen.cn/", Rule: "分类"})
		},

		Trunk: map[string]*Rule{
			"分类": {
				ItemFields: []string{
					"大类",
					"地址",
				},
				ParseFunc: func(ctx *Context) {
					query := ctx.GetDom()
					//分类
					lis := query.Find(".right .sons .cont a")
					lis.Each(func(i int, s *goquery.Selection) {
						// if i > 0 {
						// 	return
						// }
						url, _ := s.Attr("href")
						if !strings.HasPrefix(url, "http") {
							url = "https://so.gushiwen.cn" + url
						}
						title := s.Text()
						types := map[int]interface{}{
							0: title,
							1: url,
						}
						fmt.Println(types)
						ctx.Output(types)
						ctx.AddQueue(&request.Request{Url: url, Rule: "古诗名称"})
					})
				},
			},

			"古诗名称": {
				ItemFields: []string{
					"子类",
					"古诗名称",
					"地址",
				},
				ParseFunc: func(ctx *Context) {
					query := ctx.GetDom()
					query.Find(".left .sons .typecont").Each(func(i int, s *goquery.Selection) {
						// if i > 0 {
						// 	return
						// }
						subType := s.Find(".bookMl").Text()
						s.Find("span a").Each(func(i int, s1 *goquery.Selection) {
							// if i > 0 {
							// 	return
							// }
							name := s1.Text()
							url, _ := s1.Attr("href")
							if !strings.HasPrefix(url, "http") {
								url = "https://so.gushiwen.cn" + url
							}
							data := map[int]interface{}{
								0: subType,
								1: name,
								2: url,
							}
							ctx.Output(data)
							fmt.Println(data)
							tempData := map[string]interface{}{
								"type": subType,
							}
							ctx.AddQueue(&request.Request{Url: url, Rule: "内容", Temp: tempData})
						})
					})
				},
			},

			"内容": {
				//注意：有无字段语义和是否输出数据必须保持一致
				ItemFields: []string{
					"名称",
					"朝代",
					"作者",
					"内容",
					"类别",
				},
				ParseFunc: func(ctx *Context) {
					query := ctx.GetDom()
					// 获取名称
					name := query.Find("#sonsyuanwen > div.cont > h1").Text()

					// 获取朝代作者
					auth := query.Find("#sonsyuanwen > div.cont > p").Text()

					// 获取内容
					cont := query.Find("#sonsyuanwen .contson").Text()

					authArr := strings.Split(auth, "：")
					length := len(authArr)
					var caodai, authName string
					if length == 1 {
						caodai = authArr[0]
					} else if length > 1 {
						caodai = authArr[0]
						authName = authArr[1]
					}

					data := map[int]interface{}{
						0: name,
						1: caodai,
						2: authName,
						3: common.DepriveBreak(cont),
						4: ctx.GetTemp("type", ""),
					}
					// 结果存入Response中转
					fmt.Println(data)
					ctx.Output(data)
				},
			},
		},
	},
}
