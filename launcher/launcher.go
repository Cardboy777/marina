package launcher

import (
	"fmt"
	"marina/files"
	"marina/types"
	"os"
	"os/exec"
	"path/filepath"
)

func LaunchGame(version *marina.Version, onClose func(error)) error {
	path := files.GetVersionInstallDirPath(version)

	executable := getGameExecutablePath(path)

	fmt.Println(executable)

	go runGame(executable, onClose)

	return nil
}

func runGame(exePath string, onClose func(error)) {
	args := os.Args

	cmd := exec.Command(exePath, args[1:]...)
	cmd.Env = append(os.Environ(), fmt.Sprintf("SHIP_HOME=%s", filepath.Dir(exePath)))

	err := cmd.Run()
	if err != nil {
		onClose(err)
		return
	}

	fmt.Println("Successfully closed game.")
	onClose(nil)
}

func getGameExecutablePath(dirName string) string {
	installFiles, err := os.ReadDir(dirName)
	if err != nil {
		panic(fmt.Errorf("Error locating executable: %w", err))
	}

	for _, f := range installFiles {
		name := f.Name()
		if files.IsExecutable(name) {
			return filepath.Join(dirName, name)
		}
	}
	err = fmt.Errorf("No executable found in dir \"%s\"", dirName)

	panic(fmt.Errorf("Error locating executable: %w", err))
}
