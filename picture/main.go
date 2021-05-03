package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	// 存放图片链接的数据管道
	chanImgUrls chan string
	// 监控协程
	chanTask  chan string
	waitGroup sync.WaitGroup
	reImg     = `https?://[^"]+?(\.((jpg)|(png)|(jpeg)|(gif)|(bmp)))`
)

func main() {
	// 初始化管道,存取图片的 url
	chanImgUrls = make(chan string, 1000000)
	// 用于检测开启任务的协程是否完成
	chanTask = make(chan string, 2)
	// 爬取几个网页开启多少协程
	for i := 1; i <= 2; i++ {
		waitGroup.Add(1)
		url := "https://tieba.baidu.com/f?kw=%E5%9B%BE%E7%89%87&ie=utf-8&pn=" +
			strconv.Itoa((i-1)*50)
		go getImgUrls(url)
	}
	// 统计协程是否全部完成
	waitGroup.Add(1)
	go TaskStatistic()
	// 下载协程：从管道中读取链接并下载
	for i := 0; i <= 1; i++ {
		waitGroup.Add(1)
		go DownLoadImg()
	}
	waitGroup.Wait()
}

// 任务统计协程
func TaskStatistic() {
	var count int
	for {
		url := <-chanTask
		fmt.Println(url + "   网页爬取完成")
		count++
		if count == 2 {
			close(chanImgUrls)
			break
		}
	}
	waitGroup.Done()
}

// 爬图片链接到管道
func getImgUrls(url string) {
	urls := getImg(url)
	// 遍历切片里面所有图片链接，存入数据管道
	for _, u := range urls {
		chanImgUrls <- u
	}
	// 用于监控协程，看有几个协程完成
	chanTask <- url
	waitGroup.Done()
}

// 获取图片url
func getImg(url string) (urls []string) {
	pageStr := GetPageStr(url)
	compile := regexp.MustCompile(reImg)
	results := compile.FindAllStringSubmatch(pageStr, -1)
	for _, result := range results {
		urls = append(urls, result[0])
	}
	return
}

// 返回页面内容
func GetPageStr(url string) string {
	// 1.去网站拿数据
	resp, err := http.Get(url)
	HandleError(err, "http.get ")
	defer resp.Body.Close()

	// 2.读取页面内容
	PageBytes, err := ioutil.ReadAll(resp.Body)
	HandleError(err, "ioutil.ReadAll ")

	// 字节转字符串
	pageStr := string(PageBytes)
	return pageStr
}

// 下载图片
func DownLoadImg() {
	for pictureUrl := range chanImgUrls {
		filename := GetPictureName(pictureUrl)
		// 图片名字过长会报错，将名字过长的的图片剔除
		if len(filename) > 70 {
			continue
		}
		DownLoadFile(pictureUrl, filename)
	}
	waitGroup.Done()
}

func DownLoadFile(url string, filename string) bool {
	resp, err := http.Get(url)
	HandleError(err, "http.get err")
	defer resp.Body.Close()
	bytes, _ := ioutil.ReadAll(resp.Body)
	name := "/Users/yangtongchun/Desktop/Img/" + filename
	err1 := ioutil.WriteFile(name, bytes, 0666)
	if err1 != nil {
		panic(err1)
		return false
	}
	return true
}

// 截取图片名字
func GetPictureName(pictureUrl string) (filename string) {
	// 返回最后一个/的位置
	lastIndex := strings.LastIndex(pictureUrl, "/")
	name := pictureUrl[lastIndex+1:]
	// 时间戳解决文件重名
	timePrefix := strconv.Itoa(int(time.Now().UnixNano()))
	filename = timePrefix + "_" + name
	return
}

// 处理异常
func HandleError(err error, reason string) {
	if err != nil {
		fmt.Println(reason)
		return
	}
}
