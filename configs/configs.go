package configs

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Стурктура для выбора конфига
type ConfigManager struct {
	ConfigDir string
	Files     []string
}

// Сканирует директорию и сохраняет json-файлы, которые передаются в структуру(функция для обновления структуры ConfigManager)
func NewConfigManager(configDir string) (*ConfigManager, error) {
	entries, err := os.ReadDir(configDir)
	if err != nil {
		return nil, fmt.Errorf("не удалось прочитать директорию: %w", err)
	}

	//Создаем срез files в него будем передавать json файлы
	var files []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".json") {
			files = append(files, entry.Name())
		}
	}

	//Если в папке нету json файлов, то выводим ошибку
	if len(files) == 0 {
		return nil, fmt.Errorf("в %s нет json-файлов", configDir)
	}

	//Возвращаем стуктуру ConfigManger с путем до json файлов и с срезом в котором будут json файлы
	return &ConfigManager{
		ConfigDir: configDir,
		Files:     files,
	}, nil
}

// Метод который выводит на экран список конфигов
func (cm *ConfigManager) List() {
	fmt.Println("Доступные конфиги:")
	for i, name := range cm.Files {
		fmt.Printf("[%d] %s\n", i+1, name)
	}
}

// Запрашивает выбор пользователя и возвращает --config путь
func (cm *ConfigManager) Select() (string, error) {
	// Вызываем метод List который покажет доступные конфиги
	cm.List()

	fmt.Print("\nВведите номер нужного конфига: ")
	//Создаем буферизованый ввод из клавиатуры
	reader := bufio.NewReader(os.Stdin)
	//Читаем строку
	input, _ := reader.ReadString('\n')
	//Удаляем пробелы
	input = strings.TrimSpace(input)
	//Из веденой строки делаем из нее целове число и записываем в int
	var index int
	_, err := fmt.Sscanf(input, "%d", &index)
	// Если не число а что-то другое выводим ошибку
	if err != nil || index < 1 || index > len(cm.Files) {
		return "", fmt.Errorf("некорректный ввод: %s", input)
	}
	//Получаем полный путь до выбраного конфиг файла "/opt/config/config.json"
	selected := filepath.Join(cm.ConfigDir, cm.Files[index-1])
	return fmt.Sprintf("--config %s", selected), nil
}
