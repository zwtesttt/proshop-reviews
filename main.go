package main

import (
	"flag"
	"fmt"
	"shop-reviews/pkg/amazon"
	"shop-reviews/utils"
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
	reviews, err := utils.LoadExcelFile(path, "review")
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	for _, review := range reviews {
		fmt.Println("review：", review)
	}

	reqs := amazon.GenerateReviews(reviews, pid, "5")
	for _, req := range reqs {
		if err := amazon.CreateReviews(host, token, *req); err != nil {
			fmt.Println("error:", err)
			return
		}
	}

}

//eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySWQiOiI2NjIwODg2NDkzNjA3NjU2MjE4M2VkMTgiLCJpYXQiOjE3MTM0MDgxNDgsImV4cCI6MTcxNjAwMDE0OH0.jcq-LhtdVIppidly0Xuo3RzBlQNMybNhmE64IF1k0i0
