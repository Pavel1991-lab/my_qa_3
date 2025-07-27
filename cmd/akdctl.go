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
	//Проверяем что из клавиатуры были введены только нужные команды
	cmdFlag := flag.String("cmd", "", "Команда для akd: start, stop, standup, standstop, standown")
	flag.Parse()

	if *cmdFlag != "start" && *cmdFlag != "stop" && *cmdFlag != "standup" && *cmdFlag != "standown" && *cmdFlag != "standstop" && *cmdFlag != "clenar" {
		fmt.Println("Ошибка: допустимые команды — start, stop, standup, standstop, standown")
		flag.Usage()
		os.Exit(1)
	}
	//Переменная бинарика платформы
	binaryPath := "/opt/MarketingPlatform/akd"

	//Дальнейшая логика в зависимости от того какой флаг был выбран
	switch *cmdFlag {
	//Если выбарали start
	case "start":
		// Создаем экземляр структуры ConfigManager и передаем путь где лежат конфиги
		cm, err := configs.NewConfigManager("/opt/configs")
		if err != nil {
			fmt.Println("Ошибка загрузки конфигов:", err)
			os.Exit(1)
		}
		// Вызываем метод Select чтобы выбрать нужный конфиг и получить полный путь /opt/config/config.json
		configArg, err := cm.Select()
		if err != nil {
			fmt.Println("Ошибка выбора конфига:", err)
			os.Exit(1)
		}
		//Переменаая args принимает в качестве аргумента выбраный конфиг
		//Далее args передается в функцию RunCommand которая запускает бинарь
		//akd start -- config /opt/configs/main.json
		args := append([]string{"start"}, strings.Split(configArg, " ")...)
		if err := internal.RunCommand(binaryPath, args...); err != nil {
			fmt.Println("Ошибка запуска:", err)
			os.Exit(1)
		}
		fmt.Println("Платформа успешно запущена.")

	case "stop":
		//Если нужно платформу остановить передаем в функцию
		//RunCommand аргумент stop
		if err := internal.RunCommand(binaryPath, "stop"); err != nil {
			fmt.Println("Ошибка остановки:", err)
			os.Exit(1)
		}
		fmt.Println("Платформа успешно остановлена.")

	//Три отельных кейса нужны для того чтобы работать с compose файлом
	case "standup", "standstop", "standown":
		// Создаем экземляр структуры ComposeManager и передаем путь где лежат композ файлы
		cm, err := compose.NewComposeManager("/opt/compose_file") // путь к директории с docker-compose файлами
		if err != nil {
			fmt.Println("Ошибка загрузки docker-compose файлов:", err)
			os.Exit(1)

		}
		//По аналогии с ConfigManager выбираем compose файл
		configPath, err := cm.Select()
		if err != nil {
			fmt.Println("Ошибка выбора docker-compose файла:", err)
			os.Exit(1)
		}
		//Тут три возможных варианта работы есил разрушаем, стопаем или запускаем контейнеры
		switch *cmdFlag {
		//Если запускаем контенры то методу Up передаем аргумент который получили
		//благодаря мтеоду Select это путь до нашего yaml файла .../.../compose.yaml
		case "standup":
			if err := cm.Up(configPath); err != nil {
				fmt.Println("Ошибка запуска контейнеров:", err)
				os.Exit(1)
			}
			fmt.Println("Контейнеры успешно запущены.")
		//Тут все по аналогии как и в standup
		//Единственно что прежде чем мы останавливаем или
		//Разрушаем контейнеры нужно остановить платфомру
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
	//Тут мы вызываем скрипт clenar который чистит плафторму
	case "clenar":
		binaryPath := "/opt/cleaner/cliner.sh"
		if err := internal.RunCommand(binaryPath); err != nil {
			fmt.Println("Ошибка запуска:", err)
			os.Exit(1)
		}
		fmt.Println("Платформа успешно запущена.")
	}
}
