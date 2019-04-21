package main

import (
	"sync"
	"strconv"
)

var (
	//	存图片链接数据通道
	chanImageUrls chan string
	//	监控通道
	chanTask chan string
	waitGroup sync.WaitGroup
)



func main(){
	chanImageUrls = make(chan string,1000)
	chanTask = make(chan string,65)
	//	爬虫协程
	for i:=1;i<66;i++{
		waitGroup.Add(1)
		url := "http://www.umei.cc/bizhitupian/meinvbizhi/" + strconv.Itoa(i) + ".htm"
		go getImageUrls(url)
	}
	//	任务统计协程
	waitGroup.Add(1)
	go CheckOk()
	//	下载协程
	for i:=0;i<5;i++{
		waitGroup.Add(1)
		go DownloadImg()
	}
	waitGroup.Wait()
}