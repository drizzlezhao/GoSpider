package main

import (
	"regexp"
	"net/http"
	"io/ioutil"
	"fmt"
	"strings"
	"strconv"
	"time"
)

var(
	reImg string = `"(https?://[^"]+?(\.((jpg)|(jpeg)|(png)|(gif)|(ico))))"`
)

func GetPageString(url string) string{
	response, err := http.Get(url)
	HandleError(err,"http.Get url")
	defer response.Body.Close()
	pageBytes, err := ioutil.ReadAll(response.Body)
	HandleError(err,"ioutil.ReadAll")
	pageStr := string(pageBytes)
	return pageStr
}

func HandleError(err error, why string){
	if err != nil{
		fmt.Println(why,err)
	}
}


//	爬当前页所有图片链接
//	添加到管道
func getImageUrls(url string){
	urls := getImages(url)
	for _,url := range urls{
		chanImageUrls <- url
	}
	//	标志当前协程信号完成
	chanTask <- url
	waitGroup.Done()
}

//	拿图片链接
func getImages(url string)(urls []string){
	pageStr := GetPageString(url)
	//fmt.Println(pageStr)
	re := regexp.MustCompile(reImg)
	results := re.FindAllStringSubmatch(pageStr,-1)
	//fmt.Println(results)
	for _,result := range results{
		url := result[1]
		urls = append(urls,url)
	}
	return urls
}


func CheckOk(){
	var count int
	for {
		url := <- chanTask
		fmt.Printf("%s完成爬去任务",url)
		count ++
		if count == 65 {
			close(chanImageUrls)
			break
		}
	}
	waitGroup.Done()
}

func DownLoadFile(url string, filepath string) bool{
	resp,err := http.Get(url)
	if err != nil{
		return false
	}
	defer resp.Body.Close()
	fBytes, err := ioutil.ReadAll(resp.Body)
	err = ioutil.WriteFile(filepath,fBytes,644)
	//fmt.Println("err",err)
	if err == nil{
		return true
	}else {
		return false
	}
}


//	下载图片
func DownloadImg(){
	for url := range chanImageUrls{
		//	下载文件
		filepath := GetFileNameUrl(url, "F:/go_work/src/crawl/img/")
		ok := DownLoadFile(url,filepath)
		if ok{
			fmt.Printf("%s 下载成功",url)
		}else {
			fmt.Printf("%s 下载失败",url)
		}
	}
}

func GetFileNameUrl(url string, dirPath string) string{
	lastIndex := strings.LastIndex(url,"/")
	filename := url[lastIndex+1:]
	timePrefix := strconv.Itoa(int(time.Now().UnixNano()))
	filePath := dirPath + timePrefix + "_" + filename
	return filePath
}