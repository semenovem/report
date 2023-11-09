package spreadsheet

import (
	"errors"
)

const (
	csvSep = ';'

	ContentTypeXLSX    = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	ContentTypeTextCSV = "text/csv"

	FileExtensionCSV  = "csv"
	FileExtensionXLSX = "xlsx"
)

var (
	ErrUnknownContentType = errors.New("unknown content type")
	ErrNoData             = errors.New("table has is no data")
	ErrNoTableOfContent   = errors.New("no table of contents")
	ErrHeaderDuplicate    = errors.New("the table contains duplicated headers")
	ErrUnsuportedFormat   = errors.New("unsupported XLSX format")
	ErrCellCoord          = errors.New("getting cell coordinates")
	ErrCellVal            = errors.New("getting the value of a cell")
)
