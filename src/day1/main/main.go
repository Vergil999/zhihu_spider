package main

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"net/http"
	"strings"
	"time"
)


func main() {
	//http.HandleFunc("/",spiderDouban)
	//log.Println("启动了")
	//err := http.ListenAndServe(":9000",nil)
	//if err != nil{
	//	log.Fatal("List 9000")
	//}
	need :=GetIfNeedCapture()
	if need{
		capture := getCapture()
		postCapture(capture)

	}


}


func spider(w http.ResponseWriter, r *http.Request)  {
	url := "https://d.weibo.com/231650_ctg1_-_all#"
	timeout := time.Duration(5 * time.Second) //超时时间5s
	client := &http.Client{
		Timeout: timeout,
	}
	var Body io.Reader
	request, err := http.NewRequest("GET", url, Body)
	if err != nil {
		return
	}
	request.Header.Add("User-Agent", `Mozilla/5.0 (Linux; Android 6.0; Nexus 5 Build/MRA58N) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/75.0.3770.100 Mobile Safari/537.36`)
	request.Header.Add("Upgrade-Insecure-Requests", `1`)
	request.Header.Add("Referer", `https://bbs.hupu.com/`)
	request.Header.Add("Host", `bbs.hupu.com`)
	res, err := client.Do(request)

	if err != nil {
		return
	}
	defer res.Body.Close()

	var allData []map[string]interface{}
	document, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return
	}

	document.Find(".bbsHotPit li").Each(func(i int, selection *goquery.Selection) {
		s := selection.Find(".textSpan a")
		url, boolUrl := s.Attr("href")
		text := s.Text()
		if boolUrl {
			allData = append(allData, map[string]interface{}{"title": strings.TrimSpace(text), "url": "https://bbs.hupu.com/" + url})
		}
	})
	fmt.Println(allData)
	paymentDataBuf, _ := json.Marshal(&allData)
	if err != nil{

	}
	w.Write(paymentDataBuf)
}

func spiderDouban(w http.ResponseWriter, r *http.Request)  {
	client := http.Client{}
	var reader io.Reader
	req,err := http.NewRequest("GET","https://movie.douban.com/top250?start=",reader)
	if err != nil{

	}
	req.Header.Add("User-Agent", `Mozilla/5.0 (Linux; Android 6.0; Nexus 5 Build/MRA58N) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/75.0.3770.100 Mobile Safari/537.36`)
	req.Header.Add("Upgrade-Insecure-Requests", `1`)
	resp,err := client.Do(req)
	if err != nil{

	}
	defer resp.Body.Close()
	doc,err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil{

	}
	doc.Find("div[class=item]").Each(func(i int, selection *goquery.Selection) {
		img,bool := selection.Find("div[class=pic]").Find("a").Attr("href");
		if bool{
			fmt.Println(img)
		}
		detail,bool := selection.Find("div[class=hd]").Find("a").Attr("href");
		if bool{
			fmt.Println(detail)
		}
		dir_node := selection.Find("div[class=bd]")

		fmt.Println(dir_node.Eq(1).Text())
		star_node := dir_node.Next()
		score := star_node.Find("span[class=rating_num]").Text()
		fmt.Println(score)
		people_num := star_node.Last().Text()
		fmt.Println(people_num)

	})
}