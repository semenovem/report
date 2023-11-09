package spreadsheet

import (
	"github.com/pkg/errors"
	"strings"
)

type (
	Table struct {
		rows   []*Row
		cursor int
		colMap map[ColTitle]int
	}

	Row struct {
		t    *Table
		cols []string
		num  int // Номер строки в исходном файле (начинается с 1)
	}
)

// NewTable heuristic - функция для определения положения столбцов по содержимому
func NewTable(table [][]string, colTitles []ColTitle, heuristic func([]string) map[ColTitle]int) (*Table, error) {
	CompressCols(table)
	CompressRows(&table)

	switch len(table) {
	case 0:
		return nil, ErrNoTableOfContent
	case 1:
		return nil, ErrNoData
	}

	var (
		t    = &Table{}
		rows = make([]*Row, 0, len(table))
	)

	for i, r := range table[1:] {
		rows = append(rows, &Row{
			cols: r,
			num:  i + 1,
			t:    t,
		})
	}

	t.rows = rows

	var (
		colMap = make(map[ColTitle]int)
		dupls  = make([]string, 0)
	)

	// Проверить наличие оглавления
	for i, col := range GetTableOfContent(colTitles, table[0]) {
		if col == colUnknown {
			continue
		}

		if n, ok := colMap[col]; ok {
			dupls = append(dupls, table[0][n])
		} else {
			colMap[col] = i
		}
	}

	if len(dupls) != 0 {
		err := errors.WithMessagef(
			ErrHeaderDuplicate,
			"duplicated:[%s]",
			strings.Join(dupls, ", "),
		)

		return nil, err
	}

	if len(colMap) == 0 && heuristic != nil {
		colMap = heuristic(table[0])
	}

	if len(colMap) == 0 {
		return nil, ErrNoTableOfContent
	}

	t.colMap = colMap

	return t, nil
}

// Next передвигает курсор на следующую строку, и если она есть - true
func (t *Table) Next() bool {
	if t.cursor < len(t.rows) {
		t.cursor++
		return true
	}

	return false
}

// Row возвращает строку под курсором
func (t *Table) Row() *Row {
	if len(t.rows) != 0 && t.cursor-1 < len(t.rows) {
		return t.rows[t.cursor-1]
	}

	return nil
}

// Len кол-во строк в таблице
func (t *Table) Len() int {
	return len(t.rows)
}

// Title порядок столбцов оглавления
func (t *Table) Title() map[ColTitle]int {
	return t.colMap
}

// Num номер строки в исходном файле, отсчет от 1
func (r *Row) Num() int {
	return r.num
}

// Cols получить ячейки в строке.
// Значение слайса передается по ссылке, не модифицируйте без необходимости
func (r *Row) Cols() []string {
	return r.cols
}

// Val возвращает значение поля по указанному оглавлению
func (r *Row) Val(col ColTitle) string {
	i, ok := r.t.colMap[col]
	if ok && i < len(r.cols) {
		return r.cols[i]
	}

	return ""
}
