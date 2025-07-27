package internal

import (
	"fmt"
	"os"
	"os/exec"
)

// RunCommand запускает бинарник с произвольным количеством аргументов
func RunCommand(binaryPath string, args ...string) error {
	fileInfo, err := os.Stat(binaryPath)
	if err != nil {
		return fmt.Errorf("не удалось найти бинарь %s: %w", binaryPath, err)
	}
	if fileInfo.IsDir() {
		return fmt.Errorf("%s - это директория, а не исполняемый файл", binaryPath)
	}

	fmt.Printf("Запускаем команду: %s %s\n\n", binaryPath, args)

	cmd := exec.Command(binaryPath, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("ошибка выполнения команды: %w", err)
	}
	return nil
}
