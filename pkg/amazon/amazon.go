package amazon

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/tealeg/xlsx"
	"net/http"
	"os"
)

type Review struct {
	ProductId string `json:"productId"`
	Rating    string `json:"rating"`
	Comment   string `json:"comment"`
	Name      string `json:"name"`
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
