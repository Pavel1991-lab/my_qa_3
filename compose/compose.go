package compose

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Структура для выбора compose файла
type ComposeManager struct {
	composeDir string
	Files      []string
}

// Сканирует директорию и сохраняет yaml файл
func NewComposeManager(composeDir string) (*ComposeManager, error) {
	entries, err := os.ReadDir(composeDir)
	if err != nil {
		return nil, fmt.Errorf("не удалось прочитать директорию: %w", err)
	}
	//Создаем срез files в него будем передавать yaml файлы
	var files []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".yaml") {
			files = append(files, entry.Name())
		}
	}
	//Если в папке нету yaml файлов, то выводим ошибку
	if len(files) == 0 {
		return nil, fmt.Errorf("в %s нет yaml-файлов", composeDir)
	}
	//Возвращаем стуктуру ComposeManger с путем до json файлов и с срезом в котором будут json файлы
	return &ComposeManager{
		composeDir: composeDir,
		Files:      files,
	}, nil
}

// Показывает список композ файлов
func (cm *ComposeManager) List() {
	fmt.Println("Доступные композ файлы:")
	for i, name := range cm.Files {
		fmt.Printf("[%d] %s\n", i+1, name)
	}
}

// Запрашивает выбор пользователя и возвращает --compose путь
func (cm *ComposeManager) Select() (string, error) {
	// Вызываем метод List который покажет доступные compose файлы
	cm.List()

	fmt.Print("\nВведите номер нужного композ файла: ")
	//Создаем буферизованый ввод из клавиатуры
	reader := bufio.NewReader(os.Stdin)
	//Читаем строку
	input, _ := reader.ReadString('\n')
	//Удаляем пробелы
	input = strings.TrimSpace(input)
	//Из ведденой строки делаем из нее целове число и записываем в int
	var index int
	_, err := fmt.Sscanf(input, "%d", &index)
	// Если не число а что-то другое выводим ошибку
	if err != nil || index < 1 || index > len(cm.Files) {
		return "", fmt.Errorf("некорректный ввод: %s", input)
	}
	//Получаем полный путь до композ файла
	selected := filepath.Join(cm.composeDir, cm.Files[index-1])
	return selected, nil
}

// Метод запуска docker-compose up -d
func (cm *ComposeManager) Up(configPath string) error {
	cmd := exec.Command("docker-compose", "-f", configPath, "up", "-d")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Метод запуска docker-compose stop
func (cm *ComposeManager) Stop(configPath string) error {
	cmd := exec.Command("docker-compose", "-f", configPath, "stop")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Метод запуска docker-compose down -v
func (cm *ComposeManager) Down(configPath string) error {
	cmd := exec.Command("docker-compose", "-f", configPath, "down", "-v")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
