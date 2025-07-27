package configs

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type ConfigManager struct {
	ConfigDir string
	Files     []string
}

// Сканирует директорию и сохраняет json-файлы
func NewConfigManager(configDir string) (*ConfigManager, error) {
	entries, err := os.ReadDir(configDir)
	if err != nil {
		return nil, fmt.Errorf("не удалось прочитать директорию: %w", err)
	}

	var files []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".json") {
			files = append(files, entry.Name())
		}
	}

	if len(files) == 0 {
		return nil, fmt.Errorf("в %s нет json-файлов", configDir)
	}

	return &ConfigManager{
		ConfigDir: configDir,
		Files:     files,
	}, nil
}

// Показывает список конфигов
func (cm *ConfigManager) List() {
	fmt.Println("Доступные конфиги:")
	for i, name := range cm.Files {
		fmt.Printf("[%d] %s\n", i+1, name)
	}
}

// Запрашивает выбор пользователя и возвращает --config путь
func (cm *ConfigManager) Select() (string, error) {
	cm.List()

	fmt.Print("\nВведите номер нужного конфига: ")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	var index int
	_, err := fmt.Sscanf(input, "%d", &index)
	if err != nil || index < 1 || index > len(cm.Files) {
		return "", fmt.Errorf("некорректный ввод: %s", input)
	}

	selected := filepath.Join(cm.ConfigDir, cm.Files[index-1])
	return fmt.Sprintf("--config %s", selected), nil
}
