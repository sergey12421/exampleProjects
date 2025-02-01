package main

import (
	"fmt"
	"io"
	"os"
)

func openFile(fileName string) (*os.File, error) {
	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0)
	if err != nil {
		return nil, fmt.Errorf("ошибка работы с файлом: %w", err)
	}
	return file, nil
}

/* func addTask1(file *os.File, task string) error {
	_, err := file.WriteString(task + "\n")
	if err != nil {
		return fmt.Errorf("при записи задачи в файл произошла ошибка: %w", err)
	}
	return nil
}

func readTasks(file *os.File) []string {
	var tasks []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		task := scanner.Text()
		if strings.TrimSpace(task) != "" {
			tasks = append(tasks, task)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Ошибка чтения из файла: ", err)
	}
	return tasks
} */

func clearRecords(file *os.File) error {
	_, err := file.Seek(0, io.SeekStart)
	if err != nil {
		return fmt.Errorf("ошибка переноса указателя: %w", err)
	}
	err = file.Truncate(0)
	if err != nil {
		return fmt.Errorf("ошибка очищения файла: %w", err)
	}
	return nil

}

func addRecord(file *os.File, record string) error {
	_, err := file.WriteString("Добавлена задача: " + record + "\n")
	if err != nil {
		return fmt.Errorf("при записи задачи в файл произошла ошибка: %w", err)
	}
	return nil
}
