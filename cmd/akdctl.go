package main

import (
	"flag"
	"fmt"
	"os"

	"akd-control/internal"  
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

	if err := internal.RunCommand(binaryPath, *cmdFlag); err != nil {
		fmt.Println(err)
		os.Exit(1)
	} else {
		fmt.Println("\nКоманда выполнена успешно.")
	}
}