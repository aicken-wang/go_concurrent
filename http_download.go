package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"strconv"
	"time"
)

const KB int = 1024

func createDir(dir string) error {
	fInfo, err := os.Stat(dir)
	if err != nil {
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			fmt.Printf("创建目录失败err:[%s], stat ret fileInfo:[%v]\n", err, fInfo)
			return errors.New(fmt.Sprintf("Create dir:[%s] failed, error:%s", dir, err))
		}
		fmt.Printf("创建目录:%s\n", dir)
	}
	return nil
}
func bytesToSize(bs string) string {

	length, err := strconv.Atoi(bs)
	if err != nil {
		fmt.Printf("err:%s\n", err)
		return fmt.Sprintf("%s Bytes", bs)
	}
	// by zero return
	if length == 0 {
		return "0 Bytes"
	}

	var sizes = []string{"Bytes", "KB", "MB", "GB", "TB"}
	// math.Log对数函数 在向下取整 得到sizes index
	i := math.Floor(math.Log(float64(length)) / math.Log(float64(KB)))
	r := float64(length) / math.Pow(float64(KB), i)
	return strconv.FormatFloat(r, 'f', 3, 64) + " " + sizes[int(i)]
}
func readFiles(filePath, saveDir string) (baseUrl string, files []string, errMsg error) {
	capacity := 1024
	fd, err := os.Open(filePath)
	if err != nil {
		panic(err)
		errMsg = errors.New(fmt.Sprintf("readFiles open fd failed, error:[%v]", err))
		return
	}
	defer fd.Close()
	bufio := bufio.NewReader(fd)
	index := 0
	files = make([]string, capacity)
	var input string
	for {
		buf, err := bufio.ReadBytes('\n')
		isLast := false
		if err != nil {
			if err == io.EOF {
				fmt.Println("EOF")
				isLast = true
			} else {
				errMsg = errors.New(fmt.Sprintf("readFiles open fd failed, error:[%v]", err))
				return
			}
		}
		if isLast && len(buf) == 0 {
			break
		}
		fmt.Println(string(buf[:len(buf)]))
		// 最后一个没有\n
		if isLast {
			input = string(buf[:len(buf)])
		} else {
			input = string(buf[:len(buf)-1])
		}
		strLen := len(input)
		if strLen == 0 {
			continue
		}
		index++
		if index == 1 {
			baseUrl = input
			continue
		}
		url := fmt.Sprintf("%s%s", baseUrl, input)
		fmt.Printf("第 %d 本, url:%s\n", index-1, url)
		// 当超过1024条记录才使用append,避免slice扩容会影响性能
		if index-2 >= capacity {
			files = append(files, input)
		} else {
			files[index-2] = input
		}
		// 最后一条跳出
		if isLast {
			break
		}
		//time.Sleep(time.Millisecond * 200)
		//downloadPDF(url, input, saveDir)
	}
	// 截掉后面的空格
	files = files[:index-1]
	return
}
func downloadPDF(index int, url, filename, saveDir string) {
	// 避免压力 crash goroutine Sleep 100
	time.Sleep(time.Millisecond * 100)
	res, err := http.Get(url)
	if err != nil {
		fmt.Println("下载失败", err)
		return
	}
	defer res.Body.Close()
	statusCode := res.StatusCode
	status := res.Status
	fmt.Printf("status:%s,statusCode:%d\n", status, statusCode)
	if statusCode != http.StatusOK {
		fmt.Printf("响应异常 Code:%d\n", statusCode)
		return
	}
	fmt.Printf("编号:[%d]正在下载中, url:%s\n", index, url)
	contentLen, ok := res.Header["Content-Length"]
	if ok && len(contentLen) >= 1 {
		fmt.Printf("content size:%s\n", bytesToSize(contentLen[0]))
	}
	fd, err := os.Create(saveDir + filename)
	if err != nil {
		fmt.Println("创建文件失败", err)
		return
	}
	defer fd.Close()
	io.Copy(fd, res.Body)
}
func main() {
	ThreadNum := 8
	dir, _ := os.Getwd()
	// 文档保存路径，不存在创建目录
	saveDir := dir + "/rfc/"
	baseUrl, files, _ := readFiles("./rfc.txt", saveDir)
	err := createDir(saveDir)
	if err != nil {
		fmt.Println("err:", err)
		return
	}
	size := len(files)
	if size <= ThreadNum || size <= 1 {
		ThreadNum = 1
	}

	ch := make(chan int, ThreadNum)
	// 可以用wg代替
	// var wg sync.WaitGroup
	// wg.Add(size)
	Done := make(chan int)
	for i := 0; i < len(files); i++ {
		fmt.Printf("url:%s file:%s\n", baseUrl, files[i])
		ch <- i
		go func(idx int, url, filename, dir string) {
			// defer wg.Done()
			downloadPDF(idx, url, filename, dir)
			<-ch
			if size == idx {
				Done <- 0
			}
		}(i+1, baseUrl+files[i], files[i], saveDir)

	}
	// 关闭 channel，应该由写数据方关闭channel,可以避免关闭的channel中写入元素导致程序panic
	// wg.Wait()
	<-Done

	close(ch)
	close(Done)
}
