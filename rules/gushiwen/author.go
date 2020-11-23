package gushiwen

import (
	"fmt"

	"github.com/henrylee2cn/pholcus/app/downloader/request" //DOM解析
	"github.com/henrylee2cn/pholcus/common/simplejson"

	. "github.com/henrylee2cn/pholcus/app/spider" //必需
)

// func init() {
// 	Author.Register()
// }

var (
	//作者列表页
	urlAuthList string = "https://app.gushiwen.cn/api/author/Default10.aspx?c=&page=%d&token=gswapi"

	//作者详情页面
	urlAuthInfo string = "https://app.gushiwen.cn/api/author/author10.aspx?id=%s&token=gswapi"

	//作者作品列表页
	urlAuthWorks string = "https://app.gushiwen.cn/api/shiwen/Default11.aspx?page=%d&astr=%s&token=gswapi"
)

//Author 古诗文爬虫-作者
var Author = &Spider{
	Name:            "古诗文爬虫",
	Description:     "古诗文爬虫 [作者数据] [https://www.gushiwen.cn/]",
	EnableCookie:    true,
	NotDefaultField: true,
	RuleTree: &RuleTree{
		Root: func(ctx *Context) {
			for i := 1; i <= 100; i++ {
				ctx.AddQueue(&request.Request{Url: fmt.Sprintf(urlAuthList, i), Rule: "作者列表"})
			}
		},

		Trunk: map[string]*Rule{
			"作者列表": {
				ParseFunc: func(ctx *Context) {
					json, _ := simplejson.NewJson([]byte(ctx.GetText()))
					authors, _ := json.Get("authors").Array()
					for _, v := range authors {
						if item, ok := v.(map[string]interface{}); ok {
							fmt.Println(item["nameStr"])
							ctx.Output(item)
							ctx.AddQueue(&request.Request{Url: fmt.Sprintf(urlAuthInfo, item["idnew"]), Rule: "作者信息"})
							sumPage, _ := json.Get("sumPage").Int()
							for i := 1; i <= sumPage; i++ {
								ctx.AddQueue(&request.Request{Url: fmt.Sprintf(urlAuthWorks, i, item["nameStr"]), Rule: "作者作品"})
							}
						}
					}
				},
			},
			"作者信息": {
				ParseFunc: func(ctx *Context) {
					json, _ := simplejson.NewJson([]byte(ctx.GetText()))
					author, _ := json.Get("tb_author").Map()
					ziliaos, _ := json.Get("tb_ziliaos").Get("ziliaos").Array()
					ctx.Output(author)
					for _, ziliao := range ziliaos {
						if item, ok := ziliao.(map[string]interface{}); ok {
							ctx.Output(item, "作者资料")
						}
					}
				},
			},
			"作者资料": {},
			"作者作品": {
				ParseFunc: func(ctx *Context) {
					json, _ := simplejson.NewJson([]byte(ctx.GetText()))
					gushiwens, _ := json.Get("gushiwens").Array()
					for _, v := range gushiwens {
						if item, ok := v.(map[string]interface{}); ok {
							fmt.Println(item["nameStr"])
							ctx.Output(item)
						}
					}
				},
			},
		},
	},
}
