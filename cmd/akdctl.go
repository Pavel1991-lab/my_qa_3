package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"akdctl/configs" // импортируем пакет с менеджером конфигов
	"akdctl/internal"
)

func main() {
	cmdFlag := flag.String("cmd", "", "Команда для akd: start или stop")
	flag.Parse()

	if *cmdFlag != "start" && *cmdFlag != "stop" {
		fmt.Println("Ошибка: допустимые команды — start или stop")
		flag.Usage()
		os.Exit(1)
	}

	binaryPath := "/opt/MarketingPlatform/akd"

	if *cmdFlag == "start" {
		// Используем configs.NewConfigManager с префиксом configs
		cm, err := configs.NewConfigManager("/opt/configs")
		if err != nil {
			fmt.Println("Ошибка загрузки конфигов:", err)
			os.Exit(1)
		}

		configArg, err := cm.Select()
		if err != nil {
			fmt.Println("Ошибка выбора конфига:", err)
			os.Exit(1)
		}

		args := append([]string{"start"}, strings.Split(configArg, " ")...)

		if err := internal.RunCommand(binaryPath, args...); err != nil {
			fmt.Println("Ошибка запуска:", err)
			os.Exit(1)
		} else {
			fmt.Println("Платформа успешно запущена.")
		}
		return
	}

	if err := internal.RunCommand(binaryPath, "stop"); err != nil {
		fmt.Println("Ошибка остановки:", err)
		os.Exit(1)
	} else {
		fmt.Println("Платформа успешно остановлена.")
	}
}
