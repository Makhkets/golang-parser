package main

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"sync"
)

var (
	f         *excelize.File
	sheetName = "Sheet1"

	mutex     sync.Mutex
	increment = 1
)

type Data struct {
	Category     string
	Title        string
	Price        string
	Producer     string
	Description  string
	Description2 string
	Url          string
	Images       string
}

func excelInit() {
	// Создаем новый файл
	f = excelize.NewFile()

	// Создаем новый лист

	index := f.NewSheet(sheetName)

	// Добавляем данные в ячейки
	f.SetCellValue(sheetName, "A1", "Категория")
	f.SetCellValue(sheetName, "B1", "Название")
	f.SetCellValue(sheetName, "C1", "Цена")
	f.SetCellValue(sheetName, "D1", "Производитель")
	f.SetCellValue(sheetName, "E1", "Описание")
	f.SetCellValue(sheetName, "F1", "Доп. Описание")
	f.SetCellValue(sheetName, "G1", "Ссылка")
	f.SetCellValue(sheetName, "H1", "Ссылки на изображения")

	// Устанавливаем активный лист
	f.SetActiveSheet(index)

	// Сохраняем файл
	if err := f.SaveAs("example.xlsx"); err != nil {
		fmt.Println(err)
		return
	}
}

func (d *Data) Z() {

}

func (d *Data) addDataToExcel() {
	mutex.Lock()
	increment++

	f.SetCellValue(sheetName, fmt.Sprintf("A%d", increment), d.Category)
	f.SetCellValue(sheetName, fmt.Sprintf("B%d", increment), d.Title)
	f.SetCellValue(sheetName, fmt.Sprintf("C%d", increment), d.Price)
	f.SetCellValue(sheetName, fmt.Sprintf("D%d", increment), d.Producer)
	f.SetCellValue(sheetName, fmt.Sprintf("E%d", increment), d.Description)
	f.SetCellValue(sheetName, fmt.Sprintf("F%d", increment), d.Description2)
	f.SetCellValue(sheetName, fmt.Sprintf("G%d", increment), d.Url)
	f.SetCellValue(sheetName, fmt.Sprintf("H%d", increment), d.Images)

	if err := f.SaveAs("example.xlsx"); err != nil {
		fmt.Println(err)
		return
	}

	mutex.Unlock()
}
