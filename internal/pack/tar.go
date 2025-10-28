package pack

import (
	"archive/tar"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"errors"
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

// UntarGz extracts a .tar.gz archive into destDir, preserving file modes.
func UntarGz(archivePath, destDir string) error {
    f, err := os.Open(archivePath)
    if err != nil { return err }
    defer f.Close()

    gz, err := gzip.NewReader(f)
    if err != nil { return err }
    defer gz.Close()

    tr := tar.NewReader(gz)
    for {
        hdr, err := tr.Next()
        if errors.Is(err, io.EOF) { break }
        if err != nil { return err }
        target := filepath.Join(destDir, hdr.Name)
        switch hdr.Typeflag {
        case tar.TypeDir:
            if err := os.MkdirAll(target, os.FileMode(hdr.Mode)); err != nil { return err }
        case tar.TypeReg:
            if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil { return err }
            out, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.FileMode(hdr.Mode))
            if err != nil { return err }
            if _, err := io.Copy(out, tr); err != nil { out.Close(); return err }
            if err := out.Close(); err != nil { return err }
        default:
            // skip other types for now
        }
    }
    return nil
}