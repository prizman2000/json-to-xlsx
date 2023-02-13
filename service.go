package main

import (
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/xuri/excelize/v2"
)

type StringService interface {
	Uppercase(string) (string, error)
	Count(string) int
	JsonToXlsx([][]interface{}) (string, error)
}

type stringService struct{}

func (stringService) Uppercase(s string) (string, error) {
	if s == "" {
		return "", ErrEmpty
	}
	return strings.ToUpper(s), nil
}

func (stringService) Count(s string) int {
	return len(s)
}

func (stringService) JsonToXlsx(jsonData [][]interface{}) (string, error) {
	f := excelize.NewFile()
	sheetName := "Report"

	if err := f.Close(); err != nil {
		return "", ErrFileCreation
	}

	if err := f.SetSheetName("Sheet1", sheetName); err != nil {
		return "", ErrSheetRename
	}

	for i, row := range jsonData {
		if cell, err := excelize.CoordinatesToCellName(1, i + 1); err != nil {
			return "", ErrCellCoordinates
		} else {
			f.SetSheetRow(sheetName, cell, &row)
		}
	}

	f.SaveAs("TestReport.xlsx")

	buf, _ := f.WriteToBuffer()

	if resp, err := http.Post("https://genom-report.free.beeceptor.com/file", "application/octet-stream", buf); err != nil {
		return "", ErrFileStorage
	} else {
		b, _ := io.ReadAll(resp.Body)
		return string(b), nil
	}
}

// ErrEmpty is returned when input string is empty
var ErrEmpty = errors.New("empty string")
// ErrFileCreation is returned when error while file creation
var ErrFileCreation = errors.New("cannot create file")
// ErrFileCreation is returned when base sheet cannot be renamed
var ErrSheetRename = errors.New("cannot rename sheet")
// ErrCellCordinates is returned when cannot getting cell by coordinates
var ErrCellCoordinates = errors.New("cannot get cell coordinates")
// ErrFileStorage is returned when cannot post file to file storage
var ErrFileStorage = errors.New("cannot post to file storage")