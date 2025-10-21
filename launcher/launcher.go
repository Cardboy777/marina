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
	err := files.CopyRomsToVersionInstall(version)
	if err != nil {
		return err
	}

	path := files.GetVersionInstallDirPath(version)
	return launch(path, onClose)
}

func LaunchUnstableGame(version *marina.UnstableVersion, onClose func(error)) error {
	err := files.CopyRomsToUnstableVersionInstall(version)
	if err != nil {
		return err
	}

	path := files.GetUnstableVersionInstallDirPath(version)
	return launch(path, onClose)
}

func launch(installPath string, onClose func(error)) error {
	executable := getGameExecutablePath(installPath)

	fmt.Println(executable)

	go runGame(executable, onClose)

	return nil
}

func runGame(exePath string, onClose func(error)) {
	args := os.Args

	workingDirectory := filepath.Dir(exePath)

	var cmd *exec.Cmd

	// Check if running inside Flatpak
	if isRunningInFlatpak() {
		// Use flatpak-spawn to run on host system (escapes sandbox)
		spawnArgs := []string{"--host", exePath}
		spawnArgs = append(spawnArgs, args[1:]...)
		cmd = exec.Command("flatpak-spawn", spawnArgs...)
	} else {
		// Normal execution outside Flatpak
		cmd = exec.Command(exePath, args[1:]...)
	}

	cmd.Dir = workingDirectory
	cmd.Env = append(os.Environ(), fmt.Sprintf("SHIP_HOME=%s", workingDirectory))

	err := cmd.Run()
	if err != nil {
		onClose(err)
		return
	}

	fmt.Println("Successfully closed game.")
	onClose(nil)
}

func isRunningInFlatpak() bool {
	// Flatpak always provides this file inside the sandbox
	_, err := os.Stat("/.flatpak-info")
	return err == nil
}

func getGameExecutablePath(dirName string) string {
	installFiles, err := os.ReadDir(dirName)
	if err != nil {
		panic(fmt.Errorf("Error locating executable: %w", err))
	}

	for _, f := range installFiles {
		name := f.Name()
		fullPath := filepath.Join(dirName, name)
		info, err := f.Info()
		if err != nil {
			panic(fmt.Errorf("Error locating executable: %w", err))
		}
		if files.IsExecutable(info) {
			return fullPath
		}
	}
	err = fmt.Errorf("No executable found in dir \"%s\"", dirName)

	panic(fmt.Errorf("Error locating executable: %w", err))
}
