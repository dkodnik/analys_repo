package main

import (
	"fmt"
	//"sort"
	"os"
	"strings"
)

// Свойства типа
type ExtProperties struct {
	Amount  int     // общее количество файлов данного типа
	Size    int64   // суммарный размер файлов данного типа
	Percent float64 // процент данного типа от общего размера репозитория
}

type listFls struct {
	DirAmount  int                      // количество директорий в репозитории
	FileAmount int                      // количество файлов в репозитории
	FilesSize  int64                    // общий размер размер репозитория
	ExtFiles   map[string]ExtProperties // мапа свойств типа файлов в репозитории
}

var glbLF listFls

func arrayTree(path string) error {

	stream, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("Ошибка открытия")
	}
	defer stream.Close()

	thisDir, err := stream.Readdir(0)
	if err != nil {
		return fmt.Errorf("Ошибка чтения")
	}
	for _, objInf := range thisDir {
		if objInf.IsDir() && strings.HasPrefix(objInf.Name(), ".") {
			// пропускаем директорию начинающуюся с '.' (точки) - "скрытая"
			continue
		}
		if objInf.IsDir() {
			glbLF.DirAmount++
			arrayTree(path + string(os.PathSeparator) + objInf.Name())
		} else {
			glbLF.FileAmount++
			glbLF.FilesSize += objInf.Size()

			arrFlN := strings.Split(objInf.Name(), ".")
			flsExt := arrFlN[len(arrFlN)-1]

			if _, ok := glbLF.ExtFiles[flsExt]; ok {
				tmp := glbLF.ExtFiles[flsExt]
				tmp.Amount++
				tmp.Size += objInf.Size()
				glbLF.ExtFiles[flsExt] = tmp
			} else {
				glbLF.ExtFiles[flsExt] = ExtProperties{Amount: 1, Size: objInf.Size()}
			}
		}
	}
	return nil
}

func main() {
	glbLF.ExtFiles = make(map[string]ExtProperties)
	path := "." // текущий каталог
	if len(os.Args) == 2 {
		path = os.Args[1]
	}
	err := arrayTree(path)
	if err != nil {
		panic(err.Error())
	}

	/*for amp, val := range glbLF.ExtFiles {
		val.Percent = float64((val.Size * 100) / glbLF.FilesSize)
		glbLF.ExtFiles[amp] = val
	}
	// TODO: Сортировка по Percent% у MAP'ы
	//fmt.Printf("%#v\n", glbLF)*/

	fmt.Println("Анализ репозитория") //Analysis of the repository
	fmt.Println("Каталог:", path)
	fmt.Println("Количество директорий:", glbLF.DirAmount)
	fmt.Println("Количество файлов:", glbLF.FileAmount)
	fmt.Println("Содержимое:")
	for amp, val := range glbLF.ExtFiles {
		val.Percent = float64((val.Size * 100) / glbLF.FilesSize)
		fmt.Printf("	%s=%f%s\n", amp, val.Percent, "%")
		//fmt.Println("	", amp, "=", val.Percent, "%")
	}
}
