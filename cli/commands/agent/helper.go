package agent

import (
	"archive/tar"
	"bufio"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/pkg/errors"
	"github.com/usestrix/cli/domain/consts"
	repository2 "github.com/usestrix/cli/domain/repository"
)

/*
	Created by aomerk at 6/14/21 for project cli
*/

/*
	helper functions for install procedure
*/

// global constants for file
const ()

// global variables (not cool) for this file
var ()

// extractTarGz function is not a generic function for extracting .tar.gz files.
//
// It is tailored for extracting strixeyed binary from downloaded .tar.
// gz file and placing it to predetermined download location.
//
// It uses constant output paths like /usr/bin/strixeyed on *NIX systems.
func extractTarGz(version repository2.Version) error {
	zipFile, err := os.Open(consts.DownloadZipName)
	if err != nil {
		return err
	}
	// extract gzip part
	gzf, err := gzip.NewReader(zipFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// start un-tar process
	tarReader := tar.NewReader(gzf)

	strixeyedPath := filepath.Join(consts.DaemonDir, consts.DaemonName)
	strixeyedTmpPath := fmt.Sprintf("%s-%d", strixeyedPath, time.Now().UnixNano())

	// os package doesn't work well when you simply try to create with root access, instead
	// it is easier to create in a temporary path and move to permissible location.
	outFile, err := os.Create(strixeyedTmpPath)
	if err != nil {
		return err
	}

	defer func(outFile *os.File) {
		err := outFile.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(outFile)

	// consume reader buffer
	for {
		header, err := tarReader.Next()

		if errors.Is(err, io.EOF) {
			break
		}

		if err != nil {
			return err
		}

		name := header.Name

		// incoming request must be a file, not a TypeDir or etc.
		if header.Typeflag != tar.TypeReg {
			continue
		}

		// handle strixeyed binary
		if strings.Compare(name, "strixeyed") == 0 {
			// Copy max n bytes to prevent decompression bomb.
			if written, err := io.CopyN(outFile, tarReader, version.Size); err != nil {
				_ = written
				return err
			}

			err = os.Rename(strixeyedTmpPath, strixeyedPath)
			if err != nil {
				return err
			}

			// make it executable only by owner
			// #nosec
			err = os.Chmod(strixeyedPath, 0700)
			if err != nil {
				return err
			}

			// remove zip file
			err = os.Remove(consts.DownloadZipName)
			if err != nil {
				return err
			}

			return nil
		}
	}

	return nil
}

// DownloadDaemonBinary downloads from install API and places it to predesignated location.
func DownloadDaemonBinary(
	userAPIToken, agentToken string, version repository2.Version,
	downloadDomain string,
) error {
	color.Blue("Installing StrixEye Daemon with version %s", version.Version)

	url := fmt.Sprintf(
		"https://%s/get/manager/%s/manager_%s_Linux_amd64.tar.gz",
		downloadDomain, version.Version, version.Version,
	)

	// Get tar.gz for strixeyed
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	req.Header.Add("X-USER-API-TOKEN", userAPIToken)
	req.Header.Add("X-AGENT-TOKEN", agentToken)

	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		return err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(resp.Body)
	// fail on non-ok responses
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("can not download binary from server, status code :%d ", resp.StatusCode)
	}

	// create zip file to write in it
	zipFile, err := os.Create(consts.DownloadZipName)
	if err != nil {
		return err
	}

	// write into zipfile
	buf := bufio.NewWriter(zipFile)
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	_, err = buf.Write(data)
	if err != nil {
		return err
	}

	// extract downloaded tarball to daemon location
	fmt.Printf("Extracting daemon to %s\n", filepath.Join(consts.DaemonDir, consts.DaemonName))
	err = extractTarGz(version)
	if err != nil {
		return err
	}

	return nil
}
