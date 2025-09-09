//go:build !windows

package defaulter

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/rpanchyk/javaman/internal/models"
)

type PlatformDefaulter struct {
	Config *models.Config
}

func (d PlatformDefaulter) Default(version string) error {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("cannot get user home directory: %w", err)
	}

	filePath := strings.Replace(d.Config.UnixShellConfig, "~", userHomeDir, 1)
	cfgString := `[ -f "$HOME/.javaman/profile" ] && . "$HOME/.javaman/profile"`
	hasProfile, err := d.hasProfile(filePath, cfgString)
	if err != nil {
		return fmt.Errorf("could not check profile in %s: %w", filePath, err)
	}
	if !hasProfile {
		if err := d.addProfile(filePath, cfgString); err != nil {
			return fmt.Errorf("could not add profile in %s: %w", filePath, err)
		}
		fmt.Printf("profile is added to file: %s\n", filePath)
	}

	profileFile := filepath.Join(userHomeDir, ".javaman", "profile")
	if err := d.updateProfile(profileFile, version); err != nil {
		return fmt.Errorf("could not update profile in file %s: %w", profileFile, err)
	}

	fmt.Printf("User environment is set to %s version as default\n", version)
	return nil
}

func (d PlatformDefaulter) hasProfile(filePath, cfgString string) (bool, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return false, fmt.Errorf("could not open file %s: %w", filePath, err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	for {
		line, err := reader.ReadString('\n')
		if err != nil && err != io.EOF {
			return false, fmt.Errorf("could not read file %s: %w", filePath, err)
		}

		if strings.TrimSpace(line) == cfgString {
			fmt.Printf("profile is found in file %s\n", filePath)
			return true, nil
		}

		if err == io.EOF {
			break
		}
	}

	fmt.Printf("profile is not found in file %s\n", filePath)
	return false, nil
}

func (d PlatformDefaulter) addProfile(filePath, cfgString string) error {
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("could not open file %s: %w", filePath, err)
	}
	defer file.Close()

	if _, err := file.WriteString("\n" + cfgString + "\n"); err != nil {
		return fmt.Errorf("could not write file %s: %w", filePath, err)
	}
	return nil
}

func (d PlatformDefaulter) updateProfile(filePath, version string) error {
	dirPath := filepath.Dir(filePath)
	if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
		return fmt.Errorf("could not create directory %s: %w", dirPath, err)
	}

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("could not write file %s: %w", filePath, err)
	}
	defer file.Close()

	javaHome := filepath.Join(d.Config.InstallDir, version)
	javaHomeEnvVar := "export JAVA_HOME=\"" + javaHome + "\""
	pathEnvVar := "export PATH=\"$JAVA_HOME/bin:$PATH\""

	for _, line := range []string{javaHomeEnvVar, pathEnvVar} {
		if _, err := file.WriteString(line + "\n"); err != nil {
			return fmt.Errorf("could not write file %s: %w", filePath, err)
		}
	}

	fmt.Printf("profile is updated in file %s\n", filePath)
	return nil
}
