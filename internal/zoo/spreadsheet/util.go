package spreadsheet

import (
	"strconv"
	"strings"
)

func IsRowEmpty(a []string, ignoreIndex ...int) bool {
common:
	for i, v := range a {
		if v == "" {
			continue
		}

		if len(ignoreIndex) != 0 {
			for _, ign := range ignoreIndex {
				if ign == i {
					continue common
				}
			}
		}

		return false
	}

	return true
}

// CompressRows Удаляет пустые строки
func CompressRows(pointerTorRows *[][]string, ignoreIndex ...int) {
	var a = -1
	rows := *pointerTorRows

	for b := 0; b < len(rows); b++ {
		if IsRowEmpty(rows[b], ignoreIndex...) {
			if a == -1 {
				a = b
			}
			continue
		}

		if a != -1 {
			rows[a] = rows[b]
			a++
		}
	}

	if a != -1 {
		*pointerTorRows = rows[:a]
	}
}

// CompressCols Удаляет пустые столбцы
func CompressCols(rows [][]string) {
	var (
		maxLen      = 0
		excludeCols = make(map[int]struct{}) // Столбцы для удаления
	)

	for y := 0; y < len(rows); y++ {
		if len(rows[y]) > maxLen {
			maxLen = len(rows[y])
		}
	}

	for i := 0; i < maxLen; i++ {
		isEmpty := true
		for y := 0; y < len(rows); y++ {
			if len(rows[y]) > i && rows[y][i] != "" {
				isEmpty = false
			}
		}

		if isEmpty {
			excludeCols[i] = struct{}{}
		}
	}

	if len(excludeCols) == 0 {
		return
	}

	for y := 0; y < len(rows); y++ {
		rows[y] = DelCols(rows[y], excludeCols)
	}
}

// DelCols Удаляет указанные индексы в массиве
func DelCols(row []string, excludeCols map[int]struct{}) []string {
	var a = -1

	for b := 0; b < len(row); b++ {
		if _, ok := excludeCols[b]; ok {
			if a == -1 {
				a = b
			}
			continue
		}

		if a != -1 {
			row[a] = row[b]
			a++
		}
	}

	if a != -1 {
		row = row[:a]
	}

	return row
}

// EnumerateTable Добавление первого столбика с номером строки
func EnumerateTable(a [][]string) {
	for y := 0; y < len(a); y++ {
		a[y] = append([]string{strconv.Itoa(y + 1)}, a[y]...)
	}
}

func TrimSpaceInTable(a [][]string) {
	for _, row := range a {
		TrimSpaceInRow(row)
	}
}

func TrimSpaceInRow(a []string) {
	for i := 0; i < len(a); i++ {
		a[i] = strings.TrimSpace(a[i])
	}
}

func toMap[T comparable](a []T) map[T]struct{} {
	b := make(map[T]struct{})

	for _, v := range a {
		b[v] = struct{}{}
	}

	return b
}
