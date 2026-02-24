package cmd

import (
	"crypto/ed25519"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/intentregistry/intent-cli/internal/pack"
)

func TestVerifyCommand_IntegrationUnsigned(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "intent-verify-int-unsigned-*")
	if err != nil {
		t.Fatalf("tempdir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	projectDir := filepath.Join(tmpDir, "project")
	if err := scaffoldVerifyProject(projectDir); err != nil {
		t.Fatalf("scaffold project: %v", err)
	}

	pkgPath := filepath.Join(tmpDir, "unsigned.itpkg")
	if _, err := pack.CreateItpkg(projectDir, pkgPath, nil, true); err != nil {
		t.Fatalf("create unsigned package: %v", err)
	}

	cmd := RootCmd()
	cmd.SetArgs([]string{"verify", pkgPath, "--allow-unsigned"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("verify unsigned failed: %v", err)
	}
}

func TestVerifyCommand_IntegrationSigned(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "intent-verify-int-signed-*")
	if err != nil {
		t.Fatalf("tempdir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	projectDir := filepath.Join(tmpDir, "project")
	if err := scaffoldVerifyProject(projectDir); err != nil {
		t.Fatalf("scaffold project: %v", err)
	}

	pub, priv, err := ed25519.GenerateKey(nil)
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}

	pkgPath := filepath.Join(tmpDir, "signed.itpkg")
	if _, err := pack.CreateItpkg(projectDir, pkgPath, priv, false); err != nil {
		t.Fatalf("create signed package: %v", err)
	}

	pubPath := filepath.Join(tmpDir, "public_key.hex")
	if err := os.WriteFile(pubPath, []byte(hex.EncodeToString(pub)), 0o644); err != nil {
		t.Fatalf("write public key: %v", err)
	}

	cmd := RootCmd()
	cmd.SetArgs([]string{"verify", pkgPath, "--public-key", pubPath})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("verify signed failed: %v", err)
	}
}

func scaffoldVerifyProject(projectDir string) error {
	if err := os.MkdirAll(filepath.Join(projectDir, "intents"), 0o755); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Join(projectDir, "policies"), 0o755); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(projectDir, "intents", "hello.itml"), []byte("intent \"hello\"\nworkflow:\n  → return(status=\"ok\")\n"), 0o644); err != nil {
		return err
	}

	manifest := pack.ItpkgManifest{
		Name:         "@scope/test",
		Version:      "0.1.0",
		Description:  "test",
		Type:         "lib",
		ItmlVersion:  "0.1",
		Capabilities: []string{},
		Policies: map[string]interface{}{
			"security": map[string]interface{}{
				"network": map[string]interface{}{
					"outbound": map[string]interface{}{
						"deny": []string{"*"},
					},
				},
			},
		},
	}
	manifestBytes, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(projectDir, "itpkg.json"), manifestBytes, 0o644)
}
