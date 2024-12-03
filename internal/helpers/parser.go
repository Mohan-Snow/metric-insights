package helpers

import (
	"io/ioutil"
	"strings"
)

func ParseFile(filePath string) ([]string, error) {
	// Читаем содержимое файла
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	// Преобразуем данные в строку
	content := string(data)

	// Разделяем текст по точке
	sentences := strings.Split(content, ".")

	return sentences, nil
}
