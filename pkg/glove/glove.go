package glove

import (
	"fmt"
	"os"
	"os/exec"
)

// Run запускает GloVe с использованием скрипта glove.sh
func Run() error {
	fmt.Println("Запуск GloVe...")
	cmd := exec.Command("sh", "./scripts/glove.sh")
	cmd.Dir = "." // Указываем корневую директорию проекта
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("ошибка при запуске GloVe: %v", err)
	}
	return nil
}
