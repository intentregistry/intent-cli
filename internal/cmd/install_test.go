package cmd

import (
	"archive/tar"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// helper to create a small tar.gz archive from in-memory file map and return its bytes and sha256
func createTarGz(files map[string]string) ([]byte, string, error) {
	var buf strings.Builder
	// we'll write to a temp file to reuse sha logic easily
	tmpFile, err := os.CreateTemp("", "install-test-*.tar.gz")
	if err != nil {
		return nil, "", err
	}
	defer os.Remove(tmpFile.Name())

	gz := gzip.NewWriter(tmpFile)
	tw := tar.NewWriter(gz)
	for name, content := range files {
		hdr := &tar.Header{ Name: name, Mode: 0644, Size: int64(len(content)) }
		if err := tw.WriteHeader(hdr); err != nil { return nil, "", err }
		if _, err := tw.Write([]byte(content)); err != nil { return nil, "", err }
	}
	if err := tw.Close(); err != nil { return nil, "", err }
	if err := gz.Close(); err != nil { return nil, "", err }
	if _, err := tmpFile.Seek(0, 0); err != nil { return nil, "", err }
	// compute sha
	h := sha256.New()
	if _, err := io.Copy(h, tmpFile); err != nil { return nil, "", err }
	sha := hex.EncodeToString(h.Sum(nil))
	// read bytes
	if _, err := tmpFile.Seek(0, 0); err != nil { return nil, "", err }
	data, err := io.ReadAll(tmpFile)
	if err != nil { return nil, "", err }
	_ = buf // silence unused for now
	return data, sha, nil
}

func TestInstallCommand_Success(t *testing.T) {
	// Create a tarball to serve
	tarBytes, sha, err := createTarGz(map[string]string{
		"manifest.yaml": "name: example\nversion: 1.0.0\n",
		"README.md": "# Example Intent\n",
	})
	if err != nil { t.Fatalf("failed to create tarball: %v", err) }

    // Setup test HTTP server
    mux := http.NewServeMux()
    server := httptest.NewServer(mux)
    base := server.URL
    mux.HandleFunc("/v1/packages/resolve", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        q := r.URL.Query().Get("spec")
        if q == "" { http.Error(w, "missing spec", http.StatusBadRequest); return }
        meta := map[string]string{
            "name":    "@scope/example",
            "version": "1.2.3",
            "tarball": base + "/artifacts/example-1.2.3.tar.gz",
            "sha256":  sha,
        }
        _ = json.NewEncoder(w).Encode(meta)
    })
    mux.HandleFunc("/artifacts/example-1.2.3.tar.gz", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/gzip")
        w.WriteHeader(http.StatusOK)
        _, _ = w.Write(tarBytes)
    })
	defer server.Close()

	// Prepare temp working dir
	tmpDir, err := os.MkdirTemp("", "intent-install-*")
	if err != nil { t.Fatalf("tempdir: %v", err) }
	defer os.RemoveAll(tmpDir)
	cwd, _ := os.Getwd()
	defer os.Chdir(cwd)
	_ = os.Chdir(tmpDir)

	// Point CLI to test server
	apiURLFlag = server.URL

	// Run command
	cmd := InstallCmd()
	cmd.SetArgs([]string{"@scope/example@1.2.3"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("install failed: %v", err)
	}

	// Verify extraction
    target := filepath.Join("intents", "@scope-example")
	if stat, err := os.Stat(target); err != nil || !stat.IsDir() {
		t.Fatalf("expected target dir %s: %v", target, err)
	}
	if _, err := os.Stat(filepath.Join(target, "manifest.yaml")); err != nil {
		t.Fatalf("expected extracted file: %v", err)
	}
	if _, err := os.Stat(filepath.Join(target, ".installed.json")); err != nil {
		t.Fatalf("expected install manifest: %v", err)
	}
}

func TestInstallCommand_ChecksumMismatch(t *testing.T) {
	// Create tarball
	tarBytes, _, err := createTarGz(map[string]string{"file.txt": "data"})
	if err != nil { t.Fatalf("failed to create tarball: %v", err) }
	badSha := strings.Repeat("0", 64)

    mux := http.NewServeMux()
    server := httptest.NewServer(mux)
    base := server.URL
    mux.HandleFunc("/v1/packages/resolve", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        meta := map[string]string{
            "name":    "example",
            "version": "0.0.1",
            "tarball": base + "/a.tgz",
            "sha256":  badSha,
        }
        _ = json.NewEncoder(w).Encode(meta)
    })
    mux.HandleFunc("/a.tgz", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        _, _ = w.Write(tarBytes)
    })
	defer server.Close()

	tmpDir, err := os.MkdirTemp("", "intent-install-*")
	if err != nil { t.Fatalf("tempdir: %v", err) }
	defer os.RemoveAll(tmpDir)
	cwd, _ := os.Getwd()
	defer os.Chdir(cwd)
	_ = os.Chdir(tmpDir)

	apiURLFlag = server.URL

	cmd := InstallCmd()
	cmd.SetArgs([]string{"example"})
	err = cmd.Execute()
	if err == nil || !strings.Contains(err.Error(), "checksum mismatch") {
		t.Fatalf("expected checksum mismatch error, got: %v", err)
	}
}
