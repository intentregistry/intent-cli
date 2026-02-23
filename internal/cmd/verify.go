package cmd

import (
	"archive/tar"
	"compress/gzip"
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

func VerifyCmd() *cobra.Command {
	var (
		publicKeyPath string
		allowUnsigned bool
	)

	c := &cobra.Command{
		Use:   "verify <package.itpkg>",
		Short: "Verify .itpkg package integrity and signature",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			pkgPath := args[0]
			files, err := readItpkgFiles(pkgPath)
			if err != nil {
				return err
			}

			manifestContent, ok := files["MANIFEST.sha256"]
			if !ok {
				return fmt.Errorf("package missing MANIFEST.sha256")
			}
			signature, ok := files["SIGNATURE"]
			if !ok {
				return fmt.Errorf("package missing SIGNATURE")
			}

			entries, err := parseManifestSHA256(string(manifestContent))
			if err != nil {
				return err
			}

			for _, entry := range entries {
				content, ok := files[entry.Path]
				if !ok {
					return fmt.Errorf("manifest entry missing from package: %s", entry.Path)
				}
				sum := sha256.Sum256(content)
				got := hex.EncodeToString(sum[:])
				if !strings.EqualFold(got, entry.Hash) {
					return fmt.Errorf("checksum mismatch for %s: got %s expected %s", entry.Path, got, entry.Hash)
				}
			}

			fmt.Println("✅ Integrity checks passed")

			if string(signature) == "UNSIGNED" {
				if !allowUnsigned {
					return fmt.Errorf("package is unsigned; re-run with --allow-unsigned to accept unsigned artifacts")
				}
				fmt.Println("⚠️  Package is unsigned (--allow-unsigned enabled)")
				return nil
			}

			if len(signature) != ed25519.SignatureSize {
				return fmt.Errorf("invalid signature length: got %d bytes, expected %d", len(signature), ed25519.SignatureSize)
			}

			if publicKeyPath == "" {
				fmt.Println("⚠️  Signature present but not cryptographically verified (use --public-key)")
				return nil
			}

			publicKey, err := loadEd25519PublicKey(publicKeyPath)
			if err != nil {
				return fmt.Errorf("failed to load public key: %w", err)
			}

			if !ed25519.Verify(publicKey, manifestContent, signature) {
				return fmt.Errorf("signature verification failed")
			}

			fmt.Println("✅ Signature verified (ed25519)")
			return nil
		},
	}

	c.Flags().StringVar(&publicKeyPath, "public-key", "", "Path to ed25519 public key (hex format)")
	c.Flags().BoolVar(&allowUnsigned, "allow-unsigned", false, "Allow unsigned packages (integrity only)")
	return c
}

type manifestEntry struct {
	Hash string
	Path string
}

func parseManifestSHA256(content string) ([]manifestEntry, error) {
	lines := strings.Split(content, "\n")
	entries := make([]manifestEntry, 0, len(lines))
	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, "  ", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid MANIFEST.sha256 line %d: %q", i+1, line)
		}
		hash := strings.TrimSpace(parts[0])
		path := strings.TrimSpace(parts[1])
		if hash == "" || path == "" {
			return nil, fmt.Errorf("invalid MANIFEST.sha256 line %d: %q", i+1, line)
		}
		if len(hash) != 64 {
			return nil, fmt.Errorf("invalid SHA256 hash length at line %d", i+1)
		}
		entries = append(entries, manifestEntry{Hash: hash, Path: path})
	}
	if len(entries) == 0 {
		return nil, fmt.Errorf("MANIFEST.sha256 has no entries")
	}
	return entries, nil
}

func readItpkgFiles(path string) (map[string][]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open package: %w", err)
	}
	defer f.Close()

	gz, err := gzip.NewReader(f)
	if err != nil {
		return nil, fmt.Errorf("failed to read gzip stream: %w", err)
	}
	defer gz.Close()

	tr := tar.NewReader(gz)
	files := map[string][]byte{}
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to read archive: %w", err)
		}
		if hdr.Typeflag != tar.TypeReg {
			continue
		}
		data, err := io.ReadAll(tr)
		if err != nil {
			return nil, fmt.Errorf("failed to read %s: %w", hdr.Name, err)
		}
		files[hdr.Name] = data
	}

	return files, nil
}
