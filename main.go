package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
)

func main() {
	cmdFlag := flag.String("cmd", "", "Команда для akd: start или stop")
	flag.Parse()

	if *cmdFlag != "start" && *cmdFlag != "stop" {
		fmt.Println("Ошибка: допустимые команды — start или stop")
		flag.Usage()
		os.Exit(1)
	}

	// Проверяем, существует ли бинарь и доступен ли он для выполнения
	binaryPath := "/opt/MarketingPlatform/akd"
	fileInfo, err := os.Stat(binaryPath)
	if err != nil {
		fmt.Printf("Ошибка: не удалось найти бинарь %s: %v\n", binaryPath, err)
		os.Exit(1)
	}
	if fileInfo.IsDir() {
		fmt.Printf("Ошибка: %s - это директория, а не исполняемый файл\n", binaryPath)
		os.Exit(1)
	}
	
	fmt.Printf("Запускаем команду: %s %s\n\n", binaryPath, *cmdFlag)

	cmd := exec.Command(binaryPath, *cmdFlag)

	// Перенаправляем вывод команды сразу в консоль (stdout и stderr)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Запускаем команду
	err = cmd.Run()
	if err != nil {
		// Если команда завершилась с ошибкой — выводим её
		fmt.Printf("\nОшибка выполнения команды: %v\n", err)
		os.Exit(1)
	} else {
		fmt.Println("\nКоманда выполнена успешно.")
	}
}
