package amazon

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/tealeg/xlsx"
	"io/ioutil"
	"net/http"
	"os"
)

type Review struct {
	ProductId string `json:"productId"`
	Rating    string `json:"rating"`
	Comment   string `json:"comment"`
	Name      string `json:"name"`
}

type ReviewResponse struct {
	Id        string `json:"_id"`
	Name      string `json:"name"`
	Comment   string `json:"comment"`
	User      string `json:"user"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
	Rating    int    `json:"rating"`
}

func GenerateReviews(reviews []string, productId string, rating string) []*Review {
	var reviewsList []*Review
	for _, review := range reviews {
		reviewsList = append(reviewsList, &Review{
			ProductId: productId,
			Rating:    rating,
			Comment:   review,
		})
	}
	return reviewsList
}

func CreateReviews(host string, jwt string, req2 Review) error {
	url := host + "/api/products/" + req2.ProductId + "/reviews"

	jsonData, err := json.Marshal(req2)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Cookie", fmt.Sprintf("jwt=%s", jwt))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	fmt.Println("response Status:", resp.Status)
	return nil
}

func LoadExcelFile(fileName, productId string) ([]*Review, error) {
	var result []*Review
	// 打开 Excel 文件
	xlFile, err := xlsx.OpenFile(fileName)
	if err != nil {
		return nil, err
	}
	reviewIndex := -1
	nameIndex := -1
	ratingIndex := -1
	// 遍历每个工作表
	for _, sheet := range xlFile.Sheets {
		// 遍历每个表的第一行
		for cellIndex, cell := range sheet.Rows[0].Cells {
			text := cell.String()
			switch text {
			case "review":
				reviewIndex = cellIndex
				continue
			case "name":
				nameIndex = cellIndex
				continue
			case "rating":
				ratingIndex = cellIndex
				continue
			}
		}

		for i, row := range sheet.Rows {
			if i == 0 {
				continue
			}
			if reviewIndex == -1 {
				break
			}

			review := row.Cells[reviewIndex].String()
			if review == "" || review == " " || review == "Translate review to English" {
				continue
			}

			var name string
			if nameIndex != -1 {
				name = row.Cells[nameIndex].String()
			}
			rating := "5"
			if ratingIndex != -1 {
				rating = row.Cells[ratingIndex].String()
			}
			value := &Review{
				ProductId: productId,
				Rating:    rating,
				Comment:   review,
				Name:      name,
			}
			result = append(result, value)
		}
	}
	return result, nil
}

func LoadCSVFile(fileName, productId string) ([]*Review, error) {
	var result []*Review

	// 打开 CSV 文件
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// 创建 CSV 文件的读取器
	reader := csv.NewReader(file)

	// 读取 CSV 文件的内容
	lines, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	// 假设 CSV 文件的第一行是标题行，包含字段名
	reviewIndex := -1
	nameIndex := -1
	ratingIndex := -1
	for cellIndex, cell := range lines[0] {
		switch cell {
		case "review":
			reviewIndex = cellIndex
		case "name":
			nameIndex = cellIndex
		case "rating":
			ratingIndex = cellIndex
		}
	}

	// 遍历 CSV 文件的每一行，跳过标题行
	for i, line := range lines {
		if i == 0 {
			continue
		}

		review := line[reviewIndex]
		if review == "" || review == " " {
			continue
		}

		var name string
		if nameIndex != -1 {
			name = line[nameIndex]
		}

		rating := "5"
		if ratingIndex != -1 {
			rating = line[ratingIndex]
		}

		value := &Review{
			ProductId: productId,
			Rating:    rating,
			Comment:   review,
			Name:      name,
		}
		result = append(result, value)
	}

	return result, nil
}

func GetLatestProductReviews(host, id string) ([]*Review, error) {
	url := fmt.Sprintf("%s/api/products/%s/latest-reviews", host, id)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var reviews []*ReviewResponse
	err = json.Unmarshal(body, &reviews)
	if err != nil {
		return nil, err
	}

	var result []*Review
	for _, review := range reviews {
		result = append(result, &Review{
			ProductId: id,
			Rating:    fmt.Sprintf("%d", review.Rating),
			Comment:   review.Comment,
			Name:      review.Name,
		})
	}
	return result, nil
}

// 切片1减去切片2的操作
func SubtractReviews(reviews1, reviews2 []*Review) []*Review {
	// 创建一个 map 用于存储切片2中的评论
	reviews2Map := make(map[string]struct{})
	for _, review := range reviews2 {
		reviews2Map[review.Comment] = struct{}{}
	}

	// 创建一个新的切片用于存储切片1减去切片2后的评论
	var subtractedReviews []*Review

	// 遍历切片1，并检查其评论是否在切片2中存在
	for _, review := range reviews1 {
		if _, ok := reviews2Map[review.Comment]; !ok {
			// 如果评论不存在于切片2中，则将其保留在结果中
			subtractedReviews = append(subtractedReviews, review)
		}
	}

	return subtractedReviews
}
