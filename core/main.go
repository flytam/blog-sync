package core

import (
	"blog-sync/log"
	"blog-sync/util"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

type data struct {
	Id              string `json:"id"`
	Markdowncontent string `json:"markdowncontent"`
	Content         string `json:"content"`
	Tags            string `json:"tags"`
	Categories      string `json:"categories"`
	Title           string `json:"title"`
}

type apiRes struct {
	Status bool   `json:"status"`
	Data   data   `json:"data"`
	Error  string `json:"error"`
}

type articleItem struct {
	id              string
	date            string
	markdowncontent string
	content         string
	tags            []string
	categories      []string
	title           string
}

// 获取文章列表
func (core *Core) getArticle(articleList *[]*articleItem) {
	logger := log.GetLogger()

	client := http.Client{}

	for index := 0; true; index++ {
		resp, err := client.Get("https://blog.csdn.net/" + core.Csdn + "/article/list/" + strconv.Itoa(index))
		if err != nil {
			logger.Error(err.Error())
		}
		doc, err := goquery.NewDocumentFromReader(resp.Body)

		list := doc.Find(".article-item-box")

		if list.Length() > 0 {
			list.Each(func(i int, selection *goquery.Selection) {
				// id
				href, _ := selection.Find("a").Attr("href")

				// 时间
				// TODO 处理date的换行问题
				date := selection.Find(".date").Text()

				fmt.Println(href, date)
				// FindString
				reg, _ := regexp.Compile(`(\d+)$`)

				id := reg.FindString(href)

				*articleList = append(*articleList, &articleItem{date: date, id: id})
			})
		} else {
			logger.Info("加载文章id完成")
			break
		}
	}

	wg := sync.WaitGroup{}

	wg.Add(len(*articleList))
	// 加载文章详情
	for i := range *articleList {
		go func(i int) {
			// https://mp.csdn.net/mdeditor/getArticle?id
			url := "https://mp.csdn.net/mdeditor/getArticle?id=" + (*articleList)[i].id
			logger.Info("拉取url--%s", url)
			req, err := http.NewRequest("GET", url, nil)
			if err != nil {
				logger.Error(err.Error())
				return
			}
			req.Header.Add("cookie", core.Cookie)

			client := &http.Client{}
			resp, err := client.Do(req)

			var res apiRes
			json.NewDecoder(resp.Body).Decode(&res)

			if res.Status {
				d := &res.Data
				(*articleList)[i].title = d.Title
				(*articleList)[i].content = d.Content
				(*articleList)[i].markdowncontent = d.Markdowncontent
				(*articleList)[i].categories = strings.Split(d.Categories, ",")
				(*articleList)[i].tags = strings.Split(d.Tags, ",")
			} else {
				logger.Error(res.Error)
			}

			defer wg.Done()
		}(i)
	}

	wg.Wait()
}

type Core struct {
	Cookie string
	Csdn   string
	Output string
}

func (core *Core) Run() {

	logger := log.GetLogger()
	var articleList []*articleItem

	core.getArticle(&articleList)

	if !util.Exists(core.Output) {
		logger.Info("输出目录不存在，正在创建")
		if err := os.Mkdir(core.Output, 0777); err != nil {
			logger.Error(err.Error())
		}
	}

	// 生成文件
	logger.Info("开始输出文件")

	//
	wg := sync.WaitGroup{}
	wg.Add(len(articleList))
	for _, article := range articleList {
		go func(article *articleItem) {
			// 文件
			f, err := os.Create(path.Join(core.Output, article.title+".md"))
			if err != nil {
				// TODO 处理非法文件名
				logger.Error(err.Error())
			} else {
				f.WriteString(fmt.Sprintf("---\n"))
				f.WriteString(fmt.Sprintf("title: %s\n", article.title))
				f.WriteString(fmt.Sprintf("date: %s\n", article.date))
				f.WriteString(fmt.Sprintf("tags: %s\n", strings.Join(article.tags, " ")))
				f.WriteString(fmt.Sprintf("categories: %s\n", strings.Join(article.categories, " ")))
				f.WriteString(fmt.Sprintf("---\n\n"))
				f.WriteString(fmt.Sprintf("<!--more-->\n\n"))

				if len(article.markdowncontent) > 0 {
					f.WriteString(article.markdowncontent)
				} else if len(article.content) > 0 {
					f.WriteString(article.content)
				}
				defer f.Close()
				logger.Info("生成完成：%s", article.title)

			}

			wg.Done()
		}(article)
	}
	wg.Wait()
	logger.Info("生成完成")
}
