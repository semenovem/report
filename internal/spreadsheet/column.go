package spreadsheet

import "strings"

type ColTitle int

const (
	colUnknown ColTitle = iota
	ColEmail
	ColUserName
	ColUserPosition
	ColUserAvatar
	ColUserRole
	ColUserStatus
	ColUsersGroup
	ColUsersGroupRole
)

var (
	titles = map[ColTitle][]string{
		ColEmail:          {"mail", "email", "емайл", "электронная почта", "почта"},
		ColUserName:       {"имя", "fio", "name"},
		ColUserPosition:   {"позиция", "должность", "title", "position"},
		ColUserStatus:     {"статус", "status"},
		ColUserRole:       {"рол", "role"}, // рол[ь|и] используется для поиска по префиксу
		ColUsersGroup:     {"группа", "группа пользователей", "group", "users group"},
		ColUsersGroupRole: {"роль", "роль в группе пользователей", "role"},
	}
)

// GetTableOfContent Получить порядок столбцов в таблице.
// Если столбца в таблице нет - в результате его не будет
func GetTableOfContent(cols []ColTitle, row []string) []ColTitle {
	// TODO добавить ограничение на одновременный поиск полей ColUserRole и ColUsersGroupRole

	var (
		ret = make([]ColTitle, 0)
		m   = toMap(cols)
	)

	for _, v := range row {
		ret = append(ret, equalCols(m, v))
	}

	return ret
}

func equalCols(m map[ColTitle]struct{}, name string) ColTitle {
	for col := range m {
		if equalCol(col, name) {
			return col
		}
	}
	return colUnknown
}

func equalCol(col ColTitle, name string) bool {
	n := strings.ToLower(strings.TrimSpace(name))
	for _, v := range titles[col] {
		if strings.HasPrefix(n, v) {
			return true
		}
	}

	return false
}
