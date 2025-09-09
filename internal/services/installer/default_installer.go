package installer

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/rpanchyk/javaman/internal/models"
	"github.com/rpanchyk/javaman/internal/services/downloader"
	"github.com/rpanchyk/javaman/internal/services/lister"
)

type DefaultInstaller struct {
	config      *models.Config
	listFetcher lister.ListFetcher
	downloader  downloader.Downloader
}

func NewDefaultInstaller(
	config *models.Config,
	listFetcher lister.ListFetcher,
	downloader downloader.Downloader) *DefaultInstaller {

	return &DefaultInstaller{
		config:      config,
		listFetcher: listFetcher,
		downloader:  downloader,
	}
}

func (i DefaultInstaller) Install(version string) error {
	//sdk, err := i.downloader.Download(version)
	//if err != nil {
	//	return fmt.Errorf("cannot download specified SDK: %w", err)
	//}
	//fmt.Printf("Found SDK: %+v\n", *sdk)

	//sdks, err := i.listFetcher.Fetch()
	//if err != nil {
	//	return fmt.Errorf("cannot get list of SDKs: %w", err)
	//}
	//
	//sdk, err := utils.FindByVersion(version, sdks)
	//if err != nil {
	//	return fmt.Errorf("cannot find specified SDK: %w", err)
	//}
	//fmt.Printf("Found SDK: %+v\n", *sdk)

	installDir := filepath.Join(i.config.InstallDir, version)
	if _, err := os.Stat(installDir); os.IsNotExist(err) {
		sdk, err := i.downloader.Download(version)
		if err != nil {
			return fmt.Errorf("cannot download specified SDK: %w", err)
		}

		if strings.HasSuffix(sdk.FilePath, ".zip") {
			unpackedDir, err := i.unpackZip(sdk.FilePath, i.config.InstallDir)
			if err != nil {
				return fmt.Errorf("cannot unpack specified SDK: %w", err)
			}
			if err := os.Rename(unpackedDir, installDir); err != nil {
				return fmt.Errorf("cannot rename directory for specified SDK: %w", err)
			}
		} else if strings.HasSuffix(sdk.FilePath, ".tar.gz") {
			unpackedDir, err := i.unpackTarGz(sdk.FilePath, i.config.InstallDir)
			if err != nil {
				return fmt.Errorf("cannot unpack specified SDK: %w", err)
			}
			if err := os.Rename(unpackedDir, installDir); err != nil {
				return fmt.Errorf("cannot rename directory for specified SDK: %w", err)
			}
		}

		fmt.Printf("SDK has been installed: %s\n", installDir)
	} else {
		fmt.Printf("SDK already installed: %s\n", installDir)
	}

	return nil
}

func (i DefaultInstaller) unpackZip(src, dst string) (string, error) {
	archive, err := zip.OpenReader(src)
	if err != nil {
		return "", fmt.Errorf("cannot open file: %s error: %w", src, err)
	}
	defer archive.Close()

	dstDir := ""
	for _, f := range archive.File {
		filePath := filepath.Join(dst, f.Name)
		fmt.Println("unzipping", filePath)

		if !strings.HasPrefix(filePath, filepath.Clean(dst)+string(os.PathSeparator)) {
			return "", fmt.Errorf("invalid file path: %s error: %w", filePath, err)
		}
		if f.FileInfo().IsDir() {
			fmt.Println("creating directory...")
			if err = os.MkdirAll(filePath, os.ModePerm); err != nil {
				return "", fmt.Errorf("cannot create dir: %s error: %w", filePath, err)
			}
			if dstDir == "" {
				dstDir = filePath
			}
			continue
		}

		dirPath := filepath.Dir(filePath)
		if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
			return "", fmt.Errorf("could not create directory %s: %w", dirPath, err)
		}

		dstFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return "", fmt.Errorf("cannot open file: %s error: %w", filePath, err)
		}

		fileInArchive, err := f.Open()
		if err != nil {
			return "", fmt.Errorf("cannot open zip file: %s error: %w", filePath, err)
		}

		if _, err := io.Copy(dstFile, fileInArchive); err != nil {
			return "", fmt.Errorf("cannot copy file: %s error: %w", filePath, err)
		}

		dstFile.Close()
		fileInArchive.Close()
	}

	return dstDir, nil
}

func (i DefaultInstaller) unpackTarGz(src, dst string) (string, error) {
	archive, err := os.Open(src)
	if err != nil {
		return "", fmt.Errorf("cannot open file: %s error: %w", src, err)
	}
	defer archive.Close()

	uncompressedStream, err := gzip.NewReader(archive)
	if err != nil {
		return "", fmt.Errorf("cannot create reader for file: %s error: %w", src, err)
	}

	dstDir := ""
	tarReader := tar.NewReader(uncompressedStream)
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", fmt.Errorf("cannot read entry in file %s: %w", src, err)
		}

		path := filepath.Join(dst, header.Name)
		fileInfo := header.FileInfo()

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(path, os.ModePerm); err != nil {
				return "", fmt.Errorf("could not create directory %s: %w", path, err)
			}
			if dstDir == "" {
				dstDir = path
			}
		case tar.TypeReg:
			outFile, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, fileInfo.Mode().Perm())
			if err != nil {
				return "", fmt.Errorf("could not create file %s: %w", path, err)
			}
			defer outFile.Close()
			if _, err := io.Copy(outFile, tarReader); err != nil {
				return "", fmt.Errorf("could not copy file %s: %w", path, err)
			}
		case tar.TypeSymlink:
			os.Symlink(header.Linkname, path)
		default:
			return "", fmt.Errorf("unknown type: %s in %s", string(header.Typeflag), header.Name)
		}
	}

	return dstDir, nil
}
