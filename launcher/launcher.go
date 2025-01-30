package launcher

import (
	"fmt"
	"marina/filemanager"
	"marina/types"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func LaunchGame(version *marina.VersionDefinition) error {
	path := filemanager.GetVersionInstallDirPath(version)

	executable := getGameExecutablePath(path)

	go runGame(executable)

	return nil
}

func runGame(exePath string) {
	args := os.Args

	cmd := exec.Command(exePath, args[1:]...)
	cmd.Env = append(os.Environ(), fmt.Sprintf("SHIP_HOME=%s", filepath.Dir(exePath)))

	err := cmd.Run()
	if err != nil {
		fmt.Printf("\n\n==Start==\n\nExe: %s\n\nArgs: %s\n\nEnvironment: %s\n\nError: %s\n\n==End==\n\n", exePath, strings.Join(cmd.Env, " "), strings.Join(args, " "), err)
		return
	}

	fmt.Println("Successfully closed game.")
}

func getGameExecutablePath(dirName string) string {
	files, err := os.ReadDir(dirName)
	if err != nil {
		panic(fmt.Errorf("Error locating executable: %w", err))
	}

	for _, f := range files {
		name := f.Name()
		if filemanager.IsExecutable(name) {
			return filepath.Join(dirName, name)
		}
	}
	err = fmt.Errorf("No executable found in dir \"%s\"", dirName)

	panic(fmt.Errorf("Error locating executable: %w", err))
}
