package utils

import (
	"github.com/tealeg/xlsx"
)

func LoadExcelFile(fileName, colName string) ([]string, error) {
	var result []string
	// 打开 Excel 文件
	xlFile, err := xlsx.OpenFile(fileName)
	if err != nil {
		return nil, err
	}

	// 遍历每个工作表
	for _, sheet := range xlFile.Sheets {
		var colIndex int
		// 遍历每个表的第一行
		for cellIndex, cell := range sheet.Rows[0].Cells {
			text := cell.String()
			if text == colName {
				colIndex = cellIndex
				break
			}
		}
		for _, row := range sheet.Rows {
			value := row.Cells[colIndex].String()
			if value == "" || value == " " {
				continue
			}
			result = append(result, value)
		}
	}
	return result, nil
}
