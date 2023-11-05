package spreadsheet

import (
	"encoding/csv"
	"io"

	"errors"
	"github.com/xuri/excelize/v2"
)

const (
	parseUsersMaxRows    = 100000 // Максимальное кол-во загружаемых записей
	parseUsersMaxCells   = 500    // Максимальное кол-во обрабатываемых ячеек в строке
	parseExcelEmptyCells = 10     // Кол-во пустых ячеек подряд, что бы строка считалась законченной
	parseExcelEmptyRows  = 10     // Кол-во пустых строк подряд, что бы таблица считалась законченной
)

type XLSXParserConfig struct {
	MaxRows    int
	MaxCells   int
	EmptyCells int
	EmptyRows  int
}

func (p *XLSXParserConfig) check() {
	if p.MaxRows == 0 {
		p.MaxRows = parseUsersMaxRows
	}
	if p.MaxCells == 0 {
		p.MaxCells = parseUsersMaxCells
	}
	if p.EmptyCells == 0 {
		p.EmptyCells = parseExcelEmptyCells
	}
	if p.EmptyRows == 0 {
		p.EmptyRows = parseExcelEmptyRows
	}
}

func Parse(r io.Reader, contentType string) ([][]string, error) {
	switch contentType {
	case ContentTypeXLSX:
		return ParseXLSX(r, nil)
	case ContentTypeTextCSV:
		return ParseCSV(r)
	}

	return nil, ErrUnknownContentType
}

// ParseCSV Парсинг формата .csv
func ParseCSV(r io.Reader) ([][]string, error) {
	reader := csv.NewReader(r)
	reader.Comma = ';'

	records, err := reader.ReadAll()
	if err != nil {
		return nil, errors.Join(err, ErrUnsuportedFormat)
	}

	return records, nil
}

// ParseXLSX парсинг файла .xlsx
func ParseXLSX(r io.Reader, config *XLSXParserConfig) ([][]string, error) {
	var (
		ret = make([][]string, 0)
	)

	if config == nil {
		config = &XLSXParserConfig{}
	}

	config.check()

	f, err := excelize.OpenReader(r)
	if err != nil {
		return nil, errors.Join(err, ErrUnsuportedFormat)
	}

	defer f.Close()

	if f.SheetCount == 0 {
		return ret, nil
	}

	var (
		sheetName     = f.GetSheetName(f.GetActiveSheetIndex())
		axis, val     string
		emptyRowCount = 0
	)

	for rowNum := 1; rowNum < config.MaxRows; rowNum++ {
		var (
			row            = make([]string, 0)
			emptyCellCount = 0
		)

		for cellNum := 1; cellNum < config.MaxCells; cellNum++ {
			axis, err = excelize.CoordinatesToCellName(cellNum, rowNum)
			if err != nil {
				return nil, errors.Join(err, ErrCellCoord)
			}

			val, err = f.GetCellValue(sheetName, axis)
			if err != nil {
				return nil, errors.Join(err, ErrCellVal)
			}

			if val == "" {
				emptyCellCount++
			}

			if emptyCellCount > config.EmptyCells {
				break
			}

			row = append(row, val)
		}

		if IsRowEmpty(row) {
			emptyRowCount++

			if emptyRowCount > config.EmptyRows {
				break
			}
		}

		ret = append(ret, row)
	}

	return ret, nil
}
