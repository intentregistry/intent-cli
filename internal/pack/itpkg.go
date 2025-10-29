package pack

import (
	"archive/tar"
	"compress/gzip"
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// ItpkgManifest represents the itpkg.json manifest file
type ItpkgManifest struct {
	Name         string                 `json:"name"`
	Version      string                 `json:"version"`
	Description  string                 `json:"description"`
	Entry        string                 `json:"entry,omitempty"`
	Type         string                 `json:"type,omitempty"` // "app" or "lib", default "app"
	ItmlVersion  string                 `json:"itmlVersion"`
	Capabilities []string               `json:"capabilities"`
	Policies     map[string]interface{} `json:"policies"`
	Meta         *ItpkgMeta             `json:"meta,omitempty"`
}

// ItpkgMeta contains optional metadata
type ItpkgMeta struct {
	Signature *SignatureMeta `json:"signature,omitempty"`
}

// SignatureMeta contains signature key information
type SignatureMeta struct {
	Algorithm string `json:"algorithm"` // "ed25519"
	KeyID     string `json:"keyId,omitempty"`
}

// ManifestEntry represents a file entry in MANIFEST.sha256
type ManifestEntry struct {
	Hash string
	Path string
}

// CreateItpkg creates a signed .itpkg package with flat structure
// Root contains: itpkg.json, MANIFEST.sha256, SIGNATURE (if signed), and project files
func CreateItpkg(srcDir, outputPath string, signKey ed25519.PrivateKey, unsignedAllowed bool) (string, error) {
	// Read and validate itpkg.json
	manifestPath := filepath.Join(srcDir, "itpkg.json")
	manifest, err := ReadItpkgManifest(manifestPath)
	if err != nil {
		return "", fmt.Errorf("failed to read itpkg.json: %w", err)
	}

	// Validate manifest
	if err := ValidateManifest(manifest, srcDir); err != nil {
		return "", fmt.Errorf("manifest validation failed: %w", err)
	}

	// Validate directory structure
	if err := ValidateStructure(srcDir, manifest); err != nil {
		return "", fmt.Errorf("structure validation failed: %w", err)
	}

	// Create output file
	outFile, err := os.Create(outputPath)
	if err != nil {
		return "", err
	}
	defer func() {
		if closeErr := outFile.Close(); err == nil {
			err = closeErr
		}
	}()

	gz := gzip.NewWriter(outFile)
	defer gz.Close()
	tw := tar.NewWriter(gz)
	defer tw.Close()

	// Build MANIFEST.sha256 while adding files
	var manifestEntries []ManifestEntry
	
	// Add itpkg.json first
	manifestJSON, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal manifest: %w", err)
	}
	manifestHash := sha256.Sum256(manifestJSON)
	manifestEntries = append(manifestEntries, ManifestEntry{
		Hash: hex.EncodeToString(manifestHash[:]),
		Path: "itpkg.json",
	})
	if err := addBytesToTar(tw, manifestJSON, 0644, "itpkg.json"); err != nil {
		return "", fmt.Errorf("failed to add itpkg.json: %w", err)
	}

	// Add all project files and build manifest
	err = filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		// Skip itpkg.json (already added) and output file
		rel, _ := filepath.Rel(srcDir, path)
		if rel == "itpkg.json" {
			return nil
		}
		absPath, _ := filepath.Abs(path)
		absOut, _ := filepath.Abs(outputPath)
		if absPath == absOut {
			return nil
		}

		// Read file content
		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		// Compute hash
		hash := sha256.Sum256(content)
		manifestEntries = append(manifestEntries, ManifestEntry{
			Hash: hex.EncodeToString(hash[:]),
			Path: rel,
		})

		// Add to archive
		return addBytesToTar(tw, content, info.Mode(), rel)
	})
	if err != nil {
		return "", fmt.Errorf("failed to add project files: %w", err)
	}

	// Sort manifest entries by path for deterministic output
	sort.Slice(manifestEntries, func(i, j int) bool {
		return manifestEntries[i].Path < manifestEntries[j].Path
	})

	// Generate MANIFEST.sha256 (does not include itself)
	var manifestLines []string
	for _, entry := range manifestEntries {
		manifestLines = append(manifestLines, fmt.Sprintf("%s  %s", entry.Hash, entry.Path))
	}
	manifestContent := strings.Join(manifestLines, "\n") + "\n"

	if err := addBytesToTar(tw, []byte(manifestContent), 0644, "MANIFEST.sha256"); err != nil {
		return "", fmt.Errorf("failed to add MANIFEST.sha256: %w", err)
	}

	// Create signature over MANIFEST.sha256
	var signature []byte
	if signKey != nil {
		signature = ed25519.Sign(signKey, []byte(manifestContent))
		if manifest.Meta == nil {
			manifest.Meta = &ItpkgMeta{}
		}
		if manifest.Meta.Signature == nil {
			manifest.Meta.Signature = &SignatureMeta{
				Algorithm: "ed25519",
			}
		}
	} else if unsignedAllowed {
		signature = []byte("UNSIGNED")
	} else {
		return "", errors.New("signing key not provided; use --unsigned to allow unsigned package")
	}

	if err := addBytesToTar(tw, signature, 0644, "SIGNATURE"); err != nil {
		return "", fmt.Errorf("failed to add SIGNATURE: %w", err)
	}

	if err := tw.Flush(); err != nil {
		return "", err
	}
	if err := gz.Flush(); err != nil {
		return "", err
	}

	return outputPath, nil
}

// ReadItpkgManifest reads and parses itpkg.json
func ReadItpkgManifest(path string) (*ItpkgManifest, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var manifest ItpkgManifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		return nil, err
	}

	return &manifest, nil
}

// ValidateManifest validates the manifest structure
func ValidateManifest(manifest *ItpkgManifest, srcDir string) error {
	if manifest.Name == "" {
		return errors.New("name is required")
	}
	if manifest.Version == "" {
		return errors.New("version is required")
	}
	if manifest.ItmlVersion == "" {
		return errors.New("itmlVersion is required")
	}

	pkgType := manifest.Type
	if pkgType == "" {
		pkgType = "app"
	}

	if pkgType == "app" && manifest.Entry == "" {
		return errors.New("entry is required for app packages")
	}

	if pkgType == "app" && manifest.Entry != "" {
		entryPath := filepath.Join(srcDir, manifest.Entry)
		if _, err := os.Stat(entryPath); os.IsNotExist(err) {
			return fmt.Errorf("entry file not found: %s", manifest.Entry)
		}
	}

	// Validate policies
	if manifest.Policies == nil {
		return errors.New("policies are required")
	}
	if pkgType == "app" {
		security, ok := manifest.Policies["security"].(map[string]interface{})
		if !ok || security == nil {
			return errors.New("policies.security is required for app packages")
		}
		network, ok := security["network"].(map[string]interface{})
		if !ok || network == nil {
			return errors.New("policies.security.network is required for app packages")
		}
	}

	return nil
}

// ValidateStructure validates the directory structure
func ValidateStructure(srcDir string, manifest *ItpkgManifest) error {
	var errs []string
	var warnings []string

	pkgType := manifest.Type
	if pkgType == "" {
		pkgType = "app"
	}

	// Required directories
	requiredDirs := []string{"intents", "policies"}
	for _, dir := range requiredDirs {
		dirPath := filepath.Join(srcDir, dir)
		if _, err := os.Stat(dirPath); os.IsNotExist(err) {
			errs = append(errs, fmt.Sprintf("required directory missing: %s", dir))
		}
	}

	// Validate entry file exists (for app packages)
	if pkgType == "app" && manifest.Entry != "" {
		entryPath := filepath.Join(srcDir, manifest.Entry)
		if _, err := os.Stat(entryPath); os.IsNotExist(err) {
			errs = append(errs, fmt.Sprintf("entry file not found: %s", manifest.Entry))
		}
	}

	// Recommended directories (warnings)
	recommendedDirs := []string{"schemas", "tests", ".ci", "assets"}
	for _, dir := range recommendedDirs {
		dirPath := filepath.Join(srcDir, dir)
		if _, err := os.Stat(dirPath); os.IsNotExist(err) {
			warnings = append(warnings, fmt.Sprintf("recommended directory missing: %s", dir))
		}
	}

	// Check for tests in app packages
	if pkgType == "app" {
		testsPath := filepath.Join(srcDir, "tests")
		if _, err := os.Stat(testsPath); os.IsNotExist(err) {
			warnings = append(warnings, "no tests directory found for app package")
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("validation errors: %s", strings.Join(errs, "; "))
	}

	// Print warnings but don't fail
	if len(warnings) > 0 {
		fmt.Printf("⚠️  Warnings: %s\n", strings.Join(warnings, "; "))
	}

	return nil
}

func addBytesToTar(tw *tar.Writer, b []byte, mode os.FileMode, name string) error {
	hdr := &tar.Header{
		Name: name,
		Mode: int64(mode),
		Size: int64(len(b)),
	}
	if err := tw.WriteHeader(hdr); err != nil {
		return err
	}
	_, err := tw.Write(b)
	return err
}

