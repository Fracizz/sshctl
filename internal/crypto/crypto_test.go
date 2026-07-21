package crypto_test

import (
	"strings"
	"testing"

	"github.com/Fracizz/invossh/internal/crypto"
)

func TestEncryptDecryptV1Machine(t *testing.T) {
	crypto.SetMasterPassword("")
	enc, err := crypto.Encrypt("hello")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(enc, "enc:v1:") {
		t.Fatalf("expected v1, got %s", enc)
	}
	plain, err := crypto.Decrypt(enc)
	if err != nil || plain != "hello" {
		t.Fatalf("got %q err=%v", plain, err)
	}
}

func TestEncryptDecryptV2Master(t *testing.T) {
	crypto.SetMasterPassword("test-master-pw")
	crypto.SetBindMachine(false)
	defer crypto.SetMasterPassword("")

	enc, err := crypto.Encrypt("secret")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(enc, "enc:v2:") {
		t.Fatalf("expected v2, got %s", enc)
	}
	plain, err := crypto.Decrypt(enc)
	if err != nil || plain != "secret" {
		t.Fatalf("got %q err=%v", plain, err)
	}

	crypto.SetMasterPassword("wrong")
	if _, err := crypto.Decrypt(enc); err == nil {
		t.Fatal("expected decrypt failure with wrong password")
	}
}

func TestV2BindMachineMismatch(t *testing.T) {
	crypto.SetMasterPassword("pw")
	crypto.SetBindMachine(true)
	enc, err := crypto.Encrypt("x")
	if err != nil {
		t.Fatal(err)
	}
	crypto.SetBindMachine(false)
	if _, err := crypto.Decrypt(enc); err == nil {
		t.Fatal("expected failure when bind flag differs")
	}
	crypto.SetBindMachine(true)
	plain, err := crypto.Decrypt(enc)
	if err != nil || plain != "x" {
		t.Fatalf("got %q %v", plain, err)
	}
	crypto.SetMasterPassword("")
}

func TestPassthroughPlain(t *testing.T) {
	plain, err := crypto.Decrypt("not-encrypted")
	if err != nil || plain != "not-encrypted" {
		t.Fatalf("%q %v", plain, err)
	}
}
