package cmd

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"os/exec"
	"path/filepath"
	"shop-reviews/pkg/amazon"
	"shop-reviews/pkg/auth"
	"strings"
)

var App = &cli.App{
	Name:  "ProShop Review Scraper",
	Usage: "ProShop Review Scraper",
	Commands: []*cli.Command{
		{
			Name:    "insert",
			Aliases: []string{"i"},
			Usage:   "insert reviews to ProShop",
			Action:  Insert,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "host",
					Aliases: []string{"h"},
					Value:   "http://127.0.0.1:5000",
					Usage:   "Host of the API",
				},
				&cli.StringFlag{
					Name:  "path",
					Usage: "Path of the file",
				},
				&cli.StringFlag{
					Name:  "pid",
					Usage: "Product ID",
				},
				&cli.StringFlag{
					Name:  "token",
					Usage: "JWT Token",
				},
				&cli.StringFlag{
					Name:  "email",
					Usage: "Email of the user",
				},
				&cli.StringFlag{
					Name:  "password",
					Usage: "Password of the user",
				},
			},
		},
		{
			Name:    "update",
			Aliases: []string{"u"},
			Usage:   "update reviews to ProShop",
			Action:  Update,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "anis",
					Usage: "Anis ID",
				},
				&cli.StringFlag{
					Name:    "filename",
					Aliases: []string{"name"},
					Usage:   "FileName",
				},
				&cli.StringFlag{
					Name:  "path",
					Usage: "Path of the file",
				},
				&cli.StringFlag{
					Name:  "pid",
					Usage: "Product ID",
				},
				&cli.StringFlag{
					Name:  "email",
					Usage: "Email of the user",
				},
				&cli.StringFlag{
					Name:  "password",
					Usage: "Password of the user",
				},
				&cli.StringFlag{
					Name:  "host",
					Usage: "proshop host",
				},
			},
		},
	},
}

func Insert(c *cli.Context) error {
	if c.String("path") == "" {
		return fmt.Errorf("path is empty")
	}
	if c.String("pid") == "" {
		return fmt.Errorf("pid is empty")
	}
	if c.String("token") == "" {
		if c.String("email") == "" {
			return fmt.Errorf("email is empty")
		}
		if c.String("password") == "" {
			return fmt.Errorf("password is empty")
		}
	}
	if c.String("host") == "" {
		return fmt.Errorf("host is empty")
	}

	path := c.String("path")
	pid := c.String("pid")
	host := c.String("host")
	token := c.String("token")
	if c.String("token") == "" {
		token2, err := auth.AuthenticateUser(c.String("host"), c.String("email"), c.String("password"))
		if err != nil {
			return err
		}
		token = token2
	}
	fmt.Println("获取token", token)

	reviews2, err2 := amazon.GetLatestProductReviews(host, pid)
	if err2 != nil {
		fmt.Println("error:", err2)
		return err2
	}

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
		return err
	}

	if err != nil {
		fmt.Println("error:", err)
		return err2
	}

	//去重
	newReviews := amazon.SubtractReviews(reviews, reviews2)
	for _, review := range newReviews {
		fmt.Println("review：", review.Comment)
		fmt.Println("name：", review.Name)
		fmt.Println("rating：", review.Rating)
	}

	fmt.Println("新增评论：", len(newReviews))
	for _, req := range newReviews {
		if err := amazon.CreateReviews(host, token, *req); err != nil {
			fmt.Println("error:", err)
			return err
		}
	}
	return nil
}

func Update(c *cli.Context) error {
	if c.String("anis") == "" {
		return fmt.Errorf("anis is empty")
	}
	if c.String("filename") == "" {
		return fmt.Errorf("filename is empty")
	}
	if c.String("path") == "" {
		return fmt.Errorf("path is empty")
	}
	if c.String("pid") == "" {
		return fmt.Errorf("pid is empty")
	}
	if c.String("email") == "" {
		return fmt.Errorf("email is empty")
	}
	if c.String("password") == "" {
		return fmt.Errorf("password is empty")
	}

	cmd := exec.Command("amazon-buddy", "reviews", c.String("anis"), "--filename", c.String("filename"), "-n", "10")

	// 执行命令并捕获输出
	_, err := cmd.Output()
	if err != nil {
		fmt.Println("命令执行错误:", err)
		return err
	}

	// 创建一个新的cli.Context对象，设置参数
	insertContext := cli.NewContext(c.App, nil, c)

	err = insertContext.Set("path", c.String("path"))
	if err != nil {
		return err
	}
	err = insertContext.Set("pid", c.String("pid"))
	if err != nil {
		return err
	}
	err = insertContext.Set("email", c.String("email"))
	if err != nil {
		return err
	}
	err = insertContext.Set("password", c.String("password"))
	if err != nil {
		return err
	}
	host := c.String("host")
	if host == "" {
		host = "http://127.0.0.1:5000"
	}
	err = insertContext.Set("host", host)
	fmt.Println("host1:", host)
	err = Insert(insertContext)
	if err != nil {
		return err
	}
	return nil
}
