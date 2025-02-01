package main

import (
	"bufio"
	"database/sql"
	"flag"
	"fmt"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

const taskFileName string = "records.txt"

type tasks struct {
	id       int
	task     string
	complete bool
}

func main() {
	//Подключаемся к БД
	db, err := sql.Open("mysql", "root:@/taskPlaner")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	//Открываем файл для записи логов
	file, err := openFile(taskFileName)
	if err != nil {
		fmt.Println("Программа остановила свое выполнение из-за ошибки при открытии файла")
		return
	}
	defer file.Close()

	reader := bufio.NewReader(os.Stdin)
	command := ""
	help()
	for command != "exit" {
		fmt.Println("\nВведите команду (или 'exit' для выхода):")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		parts := strings.Split(input, " ")
		if len(parts) == 0 {
			continue
		}
		command = parts[0]

		addCmd := flag.NewFlagSet("add", flag.ExitOnError)
		changeCmd := flag.NewFlagSet("change", flag.ExitOnError)

		switch command {
		case "add":
			addCmd.Parse(parts[1:])
			if addCmd.NArg() < 1 {
				fmt.Println("Добавьте описание задачи.")
				continue
			}
			description := strings.Join(addCmd.Args(), " ")
			addTask(db, description)
			fmt.Println("Добавлена задача:", description)
			addRecord(file, description)
		case "list":
			showTasks(db)
		case "change":
			changeCmd.Parse(parts[1:])
			if changeCmd.NArg() < 1 {
				fmt.Println("Допишите id задачи")
				continue
			}
			changeCompleteStatus(db, changeCmd)
		case "clear":
			err := clearRecords(file)
			if err != nil {
				fmt.Println("Произошла ошибка при очистке файла - ", err)
			}
			clearTasks(db)
		case "help":
			help()
		case "exit":
			fmt.Println("Завершение работы.")
		default:
			help()
		}
	}
}

func help() {
	fmt.Println("Испольование программы: todo <команда> [аргументы]")
	fmt.Println("Команды:")
	fmt.Println("  add <описание задачи> - добавить задачу")
	fmt.Println("  list - показать задачи")
	fmt.Println("  change <номер id> - пометить задачу выполненной")
	fmt.Println("  clear - удалить все задачи")
	fmt.Println("  help - показать справочную информацию")
	fmt.Println("  exit - выйти из программы")
}

func addTask(db *sql.DB, description string) {
	_, err := db.Exec("insert into taskPlaner.tasks (task, complete) values (?, ?)",
		description, false)
	if err != nil {
		panic(err)
	}
}

func showTasks(db *sql.DB) {
	rows, err := db.Query("select id, task, complete from taskplaner.tasks")
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	tasksList := []tasks{}

	for rows.Next() {
		p := tasks{}
		err := rows.Scan(&p.id, &p.task, &p.complete)
		if err != nil {
			fmt.Println(err)
			continue
		}
		tasksList = append(tasksList, p)
	}
	fmt.Println("Список задач:")
	for _, p := range tasksList {
		fmt.Print(p.id, ") ", p.task, "\t")
		if p.complete {
			fmt.Println("->Сделано")
		} else {
			fmt.Println("->Не сделано")
		}
	}
}

func changeCompleteStatus(db *sql.DB, changeCmd *flag.FlagSet) {
	id := changeCmd.Arg(0)
	_, err := db.Exec("update tasks set complete = true where id = ?", id)
	if err != nil {
		panic(err)
	}
}

func clearTasks(db *sql.DB) {
	_, err := db.Exec("truncate table tasks")
	if err != nil {
		panic(err)
	}
	fmt.Println("Список задач очищен")
}
