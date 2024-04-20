package main

import (
	"flag"
	"fmt"
	"path/filepath"
	"shop-reviews/pkg/amazon"
	"strings"
)

var (
	path  = ""
	host  = ""
	pid   = ""
	token = ""
)

func init() {
	flag.StringVar(&path, "path", "", "file path")
	flag.StringVar(&host, "host", "", "host")
	flag.StringVar(&pid, "pid", "", "pid")
	flag.StringVar(&token, "token", "", "token")
	flag.Parse()
}

func main() {
	fmt.Println("path:", path, "host:", host, "pid:", pid)

	// 判断文件后缀名
	ext := filepath.Ext(path)
	var reviews []*amazon.Review
	var err error
	switch strings.ToLower(ext) {
	case ".xlsx":
		reviews, err = amazon.LoadExcelFile(path, pid)
	case ".csv":
		reviews, err = amazon.LoadCSVFile(path, pid)
	default:
		fmt.Println("Unsupported file format")
		return
	}

	if err != nil {
		fmt.Println("error:", err)
		return
	}

	for _, review := range reviews {
		fmt.Println("review：", review.Comment)
		fmt.Println("name：", review.Name)
		fmt.Println("rating：", review.Rating)
	}

	for _, req := range reviews {
		if err := amazon.CreateReviews(host, token, *req); err != nil {
			fmt.Println("error:", err)
			return
		}
	}
}

//eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySWQiOiI2NjIwODg2NDkzNjA3NjU2MjE4M2VkMTgiLCJpYXQiOjE3MTM0MDgxNDgsImV4cCI6MTcxNjAwMDE0OH0.jcq-LhtdVIppidly0Xuo3RzBlQNMybNhmE64IF1k0i0
