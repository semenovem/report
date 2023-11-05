package spreadsheet

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"github.com/xuri/excelize/v2"
)

// CreateXLSX excel file
func CreateXLSX(table [][]string, sheetName string) (*bytes.Buffer, error) {
	f := excelize.NewFile()

	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println("CreateXLSXFile err:", err.Error())
		}
	}()

	index, err := f.NewSheet(sheetName)
	if err != nil {
		return nil, err
	}

	for _, n := range f.GetSheetList() {
		if err = f.DeleteSheet(n); err != nil {
			return nil, err
		}
	}

	f.SetActiveSheet(index)

	for rowNum, row := range table {
		for k, v := range row {
			cellName, err := excelize.CoordinatesToCellName(k+1, rowNum+1)
			if err != nil {
				return nil, err
			}

			if err = f.SetCellStr(sheetName, cellName, v); err != nil {
				return nil, err
			}
		}
	}

	b, err := f.WriteToBuffer()
	if err != nil {
		return nil, err
	}

	return b, nil
}

// CreateCSV csv file
func CreateCSV(table [][]string) (*bytes.Buffer, error) {
	var (
		b      = new(bytes.Buffer)
		writer = csv.NewWriter(b)
	)

	writer.Comma = csvSep

	err := writer.WriteAll(table)

	return b, err
}
