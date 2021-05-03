package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)

var (
	reQQEmail = `(\d+)@qq\.com`
	// \w代表大小写字母+数字+下划线
	reEmail = `\w+@\w+\.\w+`
	// +代表一次或多次；\s\S代表各种字符；+?代表贪婪模式：带第一个”停止匹配
	// s?代表有或者没有s
	reLink   = `href="(https?://[\s\S]+?)"`
	rePhone  = `1[3-9]\d\s?\d{4}\s?\d{4}`
	reIdCard = `[1-9]\d{5}((19\d{2})|(20[012]\d{2}))((0[1-9])|(1[012]))((0[1-9])|([12]\d)|(3[01]))\d{3}[\dXx]`
	reImg    = `https?://[^"]+?(\.((jpg)|(png)|(jpeg)|(gif)|(bmp)))`
)

// 爬 QQ 邮箱
func GetQQEmail(url string) {
	pageStr := GetPageStr(url)
	// 过滤数据
	compile, err := regexp.Compile(reQQEmail)
	HandleError(err, "regexp.Compile")
	results := compile.FindAllStringSubmatch(pageStr, -1)
	// 打印匹配内容
	PrintContent(results)
}

// 爬邮箱
func GetEmail(url string) {
	pageStr := GetPageStr(url)
	// 过滤数据
	compile, err := regexp.Compile(reEmail)
	HandleError(err, "regexp.Compile")
	results := compile.FindAllStringSubmatch(pageStr, -1)
	// 打印匹配内容
	for _, t := range results {
		fmt.Println(t)
	}
}

// 爬链接
func GetLink(url string) {
	pageStr := GetPageStr(url)
	// 过滤数据
	compile, err := regexp.Compile(reLink)
	HandleError(err, "regexp.Compile")
	results := compile.FindAllStringSubmatch(pageStr, -1)
	// 打印匹配内容
	for _, t := range results {
		fmt.Println(t)
	}
}

// 爬身份证
func GetIdCard(url string) {
	pageStr := GetPageStr(url)
	// 过滤数据
	compile, err := regexp.Compile(reIdCard)
	HandleError(err, "regexp.Compile")
	results := compile.FindAllStringSubmatch(pageStr, -1)
	// 打印匹配内容
	for _, t := range results {
		fmt.Println(t)
	}
}

// 爬手机号
func GetPhone(url string) {
	pageStr := GetPageStr(url)
	// 过滤数据
	compile, err := regexp.Compile(rePhone)
	HandleError(err, "regexp.Compile")
	results := compile.FindAllStringSubmatch(pageStr, -1)
	// 打印匹配内容
	for _, t := range results {
		fmt.Println(t)
	}
}

// 爬图片
func GetImg(url string) {
	pageStr := GetPageStr(url)

	// 过滤数据
	compile, err := regexp.Compile(reImg)
	HandleError(err, "regexp.Compile")
	results := compile.FindAllStringSubmatch(pageStr, -1)
	// 打印匹配内容
	for _, t := range results {
		fmt.Println(t[0])
	}
}

// 打印匹配内容
func PrintContent(results [][]string) {
	// 邮箱去重
	result := []string{}
	temp := map[string]struct{}{}
	for _, item := range results {
		for _, t := range item {
			if _, ok := temp[t]; !ok {
				temp[t] = struct{}{}
				result = append(result, item[0]+" "+item[1])
			}
			break
		}
	}
	for i := 0; i < len(result); i++ {
		fields := strings.Fields(result[i])
		fmt.Println("邮箱：", fields[0])
		fmt.Println("qq：", fields[1])
	}
}

// 处理异常
func HandleError(err error, reason string) {
	if err != nil {
		fmt.Println(reason)
		return
	}
}

// 获取页面内容
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

func main() {
	// 爬 QQ 邮箱
	//GetQQEmail("https://tieba.baidu.com/p/6051076813?red_tag=1573533731")
	// 爬邮箱
	//GetEmail("https://tieba.baidu.com/p/6051076813?pn=1")
	// 爬链接
	//GetLink("http://www.baidu.com/s?wd=%E8%B4%B4%E5%90%A7+%E7%95%99%E4%B8%8B%E9%82%AE%E7%AE%B1")
	// 爬手机号
	//GetPhone("http://www.zhaohaowang.com/")
	// 爬身份证
	//GetIdCard("https://henan.qq.com/a/20171107/069413.htm")
	// 爬图片
	GetImg("https://tieba.baidu.com/f?kw=%E5%9B%BE%E7%89%87&ie=utf-8&pn=0")
}
