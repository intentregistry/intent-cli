package cmd

import (
	"bytes"
	"crypto/ed25519"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/intentregistry/intent-cli/internal/pack"
)

func TestVerifyCommand_SignedPackageWithPublicKey(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "intent-verify-*")
	if err != nil {
		t.Fatalf("tempdir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	projectDir := filepath.Join(tmpDir, "project")
	if err := os.MkdirAll(filepath.Join(projectDir, "intents"), 0o755); err != nil {
		t.Fatalf("mkdir intents: %v", err)
	}
	if err := os.MkdirAll(filepath.Join(projectDir, "policies"), 0o755); err != nil {
		t.Fatalf("mkdir policies: %v", err)
	}
	if err := os.WriteFile(filepath.Join(projectDir, "intents", "hello.itml"), []byte("intent \"hello\"\nworkflow:\n  → return(status=\"ok\")\n"), 0o644); err != nil {
		t.Fatalf("write hello: %v", err)
	}
	if err := os.WriteFile(filepath.Join(projectDir, "policies", "security.itml"), []byte("security: {}\n"), 0o644); err != nil {
		t.Fatalf("write policy: %v", err)
	}

	manifest := pack.ItpkgManifest{
		Name:         "@scope/test",
		Version:      "0.1.0",
		Description:  "test",
		ItmlVersion:  "0.1",
		Type:         "lib",
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
	manifestBytes, _ := json.MarshalIndent(manifest, "", "  ")
	if err := os.WriteFile(filepath.Join(projectDir, "itpkg.json"), manifestBytes, 0o644); err != nil {
		t.Fatalf("write manifest: %v", err)
	}

	pub, priv, err := ed25519.GenerateKey(nil)
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}

	pkgPath := filepath.Join(tmpDir, "test.itpkg")
	if _, err := pack.CreateItpkg(projectDir, pkgPath, priv, false); err != nil {
		t.Fatalf("create package: %v", err)
	}

	pubPath := filepath.Join(tmpDir, "public_key.hex")
	if err := os.WriteFile(pubPath, []byte(hex.EncodeToString(pub)), 0o644); err != nil {
		t.Fatalf("write public key: %v", err)
	}

	cmd := VerifyCmd()
	cmd.SetArgs([]string{pkgPath, "--public-key", pubPath})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("verify failed: %v", err)
	}
}

func TestVerifyCommand_FailsWithWrongPublicKey(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "intent-verify-badkey-*")
	if err != nil {
		t.Fatalf("tempdir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	projectDir := filepath.Join(tmpDir, "project")
	if err := os.MkdirAll(filepath.Join(projectDir, "intents"), 0o755); err != nil {
		t.Fatalf("mkdir intents: %v", err)
	}
	if err := os.MkdirAll(filepath.Join(projectDir, "policies"), 0o755); err != nil {
		t.Fatalf("mkdir policies: %v", err)
	}
	if err := os.WriteFile(filepath.Join(projectDir, "intents", "hello.itml"), []byte("intent \"hello\"\nworkflow:\n  → return(status=\"ok\")\n"), 0o644); err != nil {
		t.Fatalf("write hello: %v", err)
	}
	manifest := `{
  "name": "@scope/test",
  "version": "0.1.0",
  "description": "test",
  "type": "lib",
  "itmlVersion": "0.1",
  "capabilities": [],
  "policies": {
    "security": {
      "network": {
        "outbound": {
          "deny": ["*"]
        }
      }
    }
  }
}`
	if err := os.WriteFile(filepath.Join(projectDir, "itpkg.json"), []byte(manifest), 0o644); err != nil {
		t.Fatalf("write manifest: %v", err)
	}

	_, priv, err := ed25519.GenerateKey(nil)
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}
	pkgPath := filepath.Join(tmpDir, "test.itpkg")
	if _, err := pack.CreateItpkg(projectDir, pkgPath, priv, false); err != nil {
		t.Fatalf("create package: %v", err)
	}

	wrongPub, _, err := ed25519.GenerateKey(nil)
	if err != nil {
		t.Fatalf("generate wrong key: %v", err)
	}
	pubPath := filepath.Join(tmpDir, "wrong_public_key.hex")
	if err := os.WriteFile(pubPath, []byte(hex.EncodeToString(wrongPub)), 0o644); err != nil {
		t.Fatalf("write public key: %v", err)
	}

	cmd := VerifyCmd()
	cmd.SetArgs([]string{pkgPath, "--public-key", pubPath})
	err = cmd.Execute()
	if err == nil || !strings.Contains(err.Error(), "signature verification failed") {
		t.Fatalf("expected signature verification error, got: %v", err)
	}
}

func TestVerifyCommand_RequireSignatureFailsWithoutPublicKey(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "intent-verify-require-key-*")
	if err != nil {
		t.Fatalf("tempdir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	projectDir := filepath.Join(tmpDir, "project")
	if err := scaffoldVerifyProject(projectDir); err != nil {
		t.Fatalf("scaffold project: %v", err)
	}
	_, priv, err := ed25519.GenerateKey(nil)
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}
	pkgPath := filepath.Join(tmpDir, "signed.itpkg")
	if _, err := pack.CreateItpkg(projectDir, pkgPath, priv, false); err != nil {
		t.Fatalf("create package: %v", err)
	}

	cmd := VerifyCmd()
	cmd.SetArgs([]string{pkgPath, "--require-signature"})
	err = cmd.Execute()
	if err == nil || !strings.Contains(err.Error(), "--require-signature requires --public-key") {
		t.Fatalf("expected require-signature error, got: %v", err)
	}
}

func TestVerifyCommand_RequireSignatureFailsUnsigned(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "intent-verify-require-unsigned-*")
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
		t.Fatalf("create package: %v", err)
	}

	cmd := VerifyCmd()
	cmd.SetArgs([]string{pkgPath, "--allow-unsigned", "--require-signature"})
	err = cmd.Execute()
	if err == nil || !strings.Contains(err.Error(), "package is unsigned but --require-signature was set") {
		t.Fatalf("expected unsigned/signature error, got: %v", err)
	}
}

func TestVerifyCommand_JSONOutput(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "intent-verify-json-*")
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
		t.Fatalf("create package: %v", err)
	}

	var out bytes.Buffer
	cmd := VerifyCmd()
	cmd.SetOut(&out)
	cmd.SetErr(&out)
	cmd.SetArgs([]string{pkgPath, "--allow-unsigned", "--json"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("verify failed: %v", err)
	}

	var got verifyResult
	if err := json.Unmarshal(out.Bytes(), &got); err != nil {
		t.Fatalf("expected JSON output, got %q (err: %v)", out.String(), err)
	}
	if !got.Integrity || got.Signed || got.SignatureStatus != "unsigned" {
		t.Fatalf("unexpected verify JSON output: %+v", got)
	}
}
