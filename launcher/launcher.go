package launcher

import (
	"fmt"
	"marina/files"
	"marina/types"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func LaunchGame(version *marina.VersionDefinition, onClose func(error)) error {
	path := files.GetVersionInstallDirPath(version)

	executable := getGameExecutablePath(path)

	go runGame(executable, onClose)

	return nil
}

func runGame(exePath string, onClose func(error)) {
	args := os.Args

	cmd := exec.Command(exePath, args[1:]...)
	cmd.Env = append(os.Environ(), fmt.Sprintf("SHIP_HOME=%s", filepath.Dir(exePath)))

	err := cmd.Run()
	if err != nil {
		fmt.Printf("\n\n==Start==\n\nExe: %s\n\nArgs: %s\n\nEnvironment: %s\n\nError: %s\n\n==End==\n\n", exePath, strings.Join(cmd.Env, " "), strings.Join(args, " "), err)
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
