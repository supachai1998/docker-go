package helpers

import (
	"bytes"
	"fmt"

	"github.com/xuri/excelize/v2"
)

func FindColumnInExcelBytes(file *[]byte, targets []string) (results []map[string]string, err error) {
	ioReader := bytes.NewReader(*file)
	excel, err := excelize.OpenReader(ioReader)
	count_rows := []int{}
	if err != nil {
		return nil, err
	}

	// get all sheet
	sheets := excel.GetSheetMap()
	// loop in all sheet
	for _, sheet := range sheets {
		// get column name
		columnss, err := excel.GetCols(sheet)
		if err != nil {
			return nil, err
		}

		// loop in all column
		for _, columns := range columnss {
			count_rows = append(count_rows, len(columns))
			for _, column := range columns {
				if column != "" && Contains(targets, column) {
					for _, valColumns := range columns {
						if !Contains(targets, valColumns) {
							results = append(results, map[string]string{column: valColumns})
						}
					}
				}
			}
		}

	}
	// len in all column should be same
	if !Every(count_rows, count_rows[0]) {
		return nil, fmt.Errorf("column length is not same %v", count_rows)
	}
	return results, nil
}
func FindColumnInExcelFile(file string, targets []string) (results []map[string]string, err error) {

	excel, err := excelize.OpenFile(file)
	count_rows := []int{}
	if err != nil {
		return nil, err
	}
	// get all sheet
	sheets := excel.GetSheetMap()
	// loop in all sheet
	for _, sheet := range sheets {
		// get column name
		columnss, err := excel.GetCols(sheet)
		if err != nil {
			return nil, err
		}
		// loop in all column
		for _, columns := range columnss {
			count_rows = append(count_rows, len(columns))
			for _, column := range columns {
				if column != "" && Contains(targets, column) {
					for _, valColumns := range columns {
						if !Contains(targets, valColumns) {
							results = append(results, map[string]string{column: valColumns})
						}
					}
				}
			}
		}

	}
	// len in all column should be same
	if !Every(count_rows, count_rows[0]) {
		return nil, fmt.Errorf("column length is not same %v", count_rows)
	}
	return results, nil
}
