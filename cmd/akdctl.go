package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"akdctl/compose"
	"akdctl/configs"
	"akdctl/internal"
)

func main() {
	cmdFlag := flag.String("cmd", "", "Команда для akd: start, stop, standup, standstop, standown")
	flag.Parse()

	if *cmdFlag != "start" && *cmdFlag != "stop" && *cmdFlag != "standup" && *cmdFlag != "standown" && *cmdFlag != "standstop" {
		fmt.Println("Ошибка: допустимые команды — start, stop, standup, standstop, standown")
		flag.Usage()
		os.Exit(1)
	}

	binaryPath := "/opt/MarketingPlatform/akd"

	switch *cmdFlag {
	case "start":
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
		}
		fmt.Println("Платформа успешно запущена.")

	case "stop":
		if err := internal.RunCommand(binaryPath, "stop"); err != nil {
			fmt.Println("Ошибка остановки:", err)
			os.Exit(1)
		}
		fmt.Println("Платформа успешно остановлена.")

	case "standup", "standstop", "standown":
		cm, err := compose.NewComposeManager("/opt/compose_file") // путь к директории с docker-compose файлами
		if err != nil {
			fmt.Println("Ошибка загрузки docker-compose файлов:", err)
			os.Exit(1)

		}

		configPath, err := cm.Select()
		if err != nil {
			fmt.Println("Ошибка выбора docker-compose файла:", err)
			os.Exit(1)
		}

		switch *cmdFlag {
		case "standup":
			if err := cm.Up(configPath); err != nil {
				fmt.Println("Ошибка запуска контейнеров:", err)
				os.Exit(1)
			}
			fmt.Println("Контейнеры успешно запущены.")

		case "standstop":
			if err := internal.RunCommand(binaryPath, "stop"); err != nil {
				fmt.Println("Ошибка остановки:", err)
				os.Exit(1)
			}
			fmt.Println("Платформа успешно остановлена.")

			if err := cm.Stop(configPath); err != nil {
				fmt.Println("Ошибка остановки контейнеров:", err)
				os.Exit(1)
			}
			fmt.Println("Контейнеры успешно остановлены.")

		case "standown":
			if err := internal.RunCommand(binaryPath, "stop"); err != nil {
				fmt.Println("Ошибка остановки:", err)
				os.Exit(1)
			}
			fmt.Println("Платформа успешно остановлена.")

			if err := cm.Down(configPath); err != nil {
				fmt.Println("Ошибка удаления контейнеров:", err)
				os.Exit(1)
			}
			fmt.Println("Контейнеры и тома успешно удалены.")
		}
	}
}
