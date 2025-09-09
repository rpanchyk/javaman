package utils

import (
	"runtime"

	"github.com/rpanchyk/javaman/internal/models"
)

/*
Used values from https://github.com/golang/go/blob/master/src/internal/syslist/syslist.go
*/

func CurrentOs() models.Os {
	switch runtime.GOOS {
	case "linux":
		return models.Linux
	case "darwin":
		return models.Macos
	case "windows":
		return models.Windows
	default:
		return models.Unknown
	}
}

func CurrentArch() models.Arch {
	switch runtime.GOARCH {
	case "amd64":
		return models.X64
	case "arm":
		return models.ARM
	default:
		return models.UNKNOWN
	}
}
