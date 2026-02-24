package cmd

import (
	"crypto/ed25519"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"os"
	"path/filepath"
	"testing"
)

func TestLoadEd25519Key_HexAndPEM(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "intent-key-loader-*")
	if err != nil {
		t.Fatalf("tempdir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	_, priv, err := ed25519.GenerateKey(nil)
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}

	hexPath := filepath.Join(tmpDir, "private.hex")
	if err := os.WriteFile(hexPath, []byte(hex.EncodeToString(priv)), 0o600); err != nil {
		t.Fatalf("write hex key: %v", err)
	}
	gotHex, err := loadEd25519Key(hexPath)
	if err != nil {
		t.Fatalf("load hex private key: %v", err)
	}
	if string(gotHex) != string(priv) {
		t.Fatalf("loaded hex private key mismatch")
	}

	der, err := x509.MarshalPKCS8PrivateKey(priv)
	if err != nil {
		t.Fatalf("marshal pkcs8: %v", err)
	}
	pemPath := filepath.Join(tmpDir, "private.pem")
	pemBytes := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: der})
	if err := os.WriteFile(pemPath, pemBytes, 0o600); err != nil {
		t.Fatalf("write pem key: %v", err)
	}
	gotPEM, err := loadEd25519Key(pemPath)
	if err != nil {
		t.Fatalf("load pem private key: %v", err)
	}
	if string(gotPEM) != string(priv) {
		t.Fatalf("loaded pem private key mismatch")
	}
}

func TestLoadEd25519PublicKey_HexAndPEM(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "intent-pub-loader-*")
	if err != nil {
		t.Fatalf("tempdir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	pub, _, err := ed25519.GenerateKey(nil)
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}

	hexPath := filepath.Join(tmpDir, "public.hex")
	if err := os.WriteFile(hexPath, []byte(hex.EncodeToString(pub)), 0o644); err != nil {
		t.Fatalf("write hex public key: %v", err)
	}
	gotHex, err := loadEd25519PublicKey(hexPath)
	if err != nil {
		t.Fatalf("load hex public key: %v", err)
	}
	if string(gotHex) != string(pub) {
		t.Fatalf("loaded hex public key mismatch")
	}

	der, err := x509.MarshalPKIXPublicKey(pub)
	if err != nil {
		t.Fatalf("marshal pkix public key: %v", err)
	}
	pemPath := filepath.Join(tmpDir, "public.pem")
	pemBytes := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: der})
	if err := os.WriteFile(pemPath, pemBytes, 0o644); err != nil {
		t.Fatalf("write pem public key: %v", err)
	}
	gotPEM, err := loadEd25519PublicKey(pemPath)
	if err != nil {
		t.Fatalf("load pem public key: %v", err)
	}
	if string(gotPEM) != string(pub) {
		t.Fatalf("loaded pem public key mismatch")
	}
}
