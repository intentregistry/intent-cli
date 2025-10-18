package pack

import (
	"archive/tar"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
	"path/filepath"
)

func TarGz(srcDir string) (tarPath, sha string, err error) {
	tmp := filepath.Join(os.TempDir(), "intent-package.tgz")
	f, err := os.Create(tmp)
	if err != nil { return "", "", err }
	defer func() {
		if closeErr := f.Close(); closeErr != nil && err == nil {
			err = closeErr
		}
	}()

	gz := gzip.NewWriter(f)
	defer func() {
		if closeErr := gz.Close(); closeErr != nil && err == nil {
			err = closeErr
		}
	}()

	tw := tar.NewWriter(gz)
	defer func() {
		if closeErr := tw.Close(); closeErr != nil && err == nil {
			err = closeErr
		}
	}()

	err = filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil { return err }
		if info.IsDir() { return nil }
		rel, _ := filepath.Rel(srcDir, path)
		hdr, err := tar.FileInfoHeader(info, "")
		if err != nil { return err }
		hdr.Name = rel
		if err := tw.WriteHeader(hdr); err != nil { return err }
		in, err := os.Open(path)
		if err != nil { return err }
		defer func() {
			if closeErr := in.Close(); closeErr != nil && err == nil {
				err = closeErr
			}
		}()
		_, err = io.Copy(tw, in)
		return err
	})
	if err != nil { return "", "", err }

	// calcular sha256
	if err := gz.Flush(); err != nil { return "", "", err }
	if err := tw.Flush(); err != nil { return "", "", err }

	if _, err := f.Seek(0, 0); err != nil { return "", "", err }
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil { return "", "", err }
	sum := hex.EncodeToString(h.Sum(nil))
	return tmp, sum, nil
}