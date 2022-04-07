package main

import (
	"fmt"
	"github.com/gocolly/colly"
	"log"
	"os"
)

func main() {
	file1, err := os.OpenFile("joke.txt", os.O_RDWR|os.O_APPEND|os.O_CREATE, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer file1.Close()

	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.75 Safari/537.36"),
		colly.AllowedDomains("xiaodiaodaya.cn"),
		colly.AllowURLRevisit(),
		colly.CacheDir("./joke_cache"),
	)
	c.DetectCharset = true

	c1 := c.Clone()
	c2 := c.Clone()
	c.OnRequest(func(r *colly.Request) {
		r.Headers.Add("cookie", "ASP.NET_SessionId=2cwokpmoimdhanqqmzvrvp2n; GUID=2674410612080532; Hm_lvt_2331286a697422bbd3257dce1b849877=1649218081; Hm_lpvt_2331286a697422bbd3257dce1b849877=1649218096")
		fmt.Println("Visiting", r.URL)
	})

	c.OnError(func(_ *colly.Response, err error) {
		log.Println("Something went wrong:", err)
	})

	c.OnResponse(func(r *colly.Response) {
		fmt.Println("Visited", r.Request.URL)

	})

	c.OnScraped(func(r *colly.Response) {
		fmt.Println("Finished", r.Request.URL)
	})
	//访问主页面，获取子链接
	c.OnHTML(`div[class=hotword_list]`, func(e *colly.HTMLElement) {
		e.ForEach("a[href]", func(_ int, element *colly.HTMLElement) {
			nextUrl := element.Attr("href")
			//拿到小分类的子链接，并访问呢
			c1.Visit(e.Request.AbsoluteURL(nextUrl))
		})
	})

	c1.OnHTML(`body`, func(e *colly.HTMLElement) {
		e.ForEach("a[href]", func(_ int, element *colly.HTMLElement) {
			nextUrl := element.Attr("href")
			//拿到具体文章的子链接并访问
			c2.Visit(e.Request.AbsoluteURL(nextUrl))
		})
	})

	c2.OnHTML(`div[class=content]`, func(e *colly.HTMLElement) {
		fmt.Println(e.Text)
		// 往创建的文件中写入笑话
		_, err = file1.WriteString(e.Text)
		if err != nil {
			panic(err)
		}
	})

	c.Visit("http://xiaodiaodaya.cn/")

}
