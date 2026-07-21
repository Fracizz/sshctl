package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"os"
	"os/user"
	"runtime"
	"strings"
	"sync"

	"golang.org/x/crypto/argon2"
)

const (
	prefixV1 = "enc:v1:"
	prefixV2 = "enc:v2:"
	saltSize = 16
)

var (
	mu               sync.RWMutex
	masterPassword   string
	bindMachineFlag  bool
	bindMachineSet   bool
)

// IsEncrypted reports whether value uses sshfrac ciphertext (v1 or v2).
func IsEncrypted(value string) bool {
	return strings.HasPrefix(value, prefixV1) || strings.HasPrefix(value, prefixV2)
}

// SetMasterPassword sets the process-wide master password (empty clears).
func SetMasterPassword(pw string) {
	mu.Lock()
	defer mu.Unlock()
	masterPassword = pw
}

// SetBindMachine enables/disables mixing machine material into the v2 KDF.
func SetBindMachine(v bool) {
	mu.Lock()
	defer mu.Unlock()
	bindMachineFlag = v
	bindMachineSet = true
}

// Encrypt encrypts plaintext.
// Prefer enc:v2 with master password (flag/env); otherwise enc:v1 machine-derived key.
func Encrypt(plaintext string) (string, error) {
	if pw, ok := resolveMasterPassword(); ok && pw != "" {
		return encryptV2(plaintext, pw, resolveBindMachine())
	}
	return encryptV1(plaintext)
}

// Decrypt decrypts enc:v1 / enc:v2. Plaintext is returned unchanged.
func Decrypt(value string) (string, error) {
	switch {
	case strings.HasPrefix(value, prefixV2):
		pw, ok := resolveMasterPassword()
		if !ok || pw == "" {
			return "", errors.New("enc:v2 requires master password (--master-password or SSHFRAC_MASTER_PASSWORD)")
		}
		return decryptV2(strings.TrimPrefix(value, prefixV2), pw, resolveBindMachine())
	case strings.HasPrefix(value, prefixV1):
		return decryptV1(strings.TrimPrefix(value, prefixV1))
	default:
		return value, nil
	}
}

func resolveMasterPassword() (string, bool) {
	mu.RLock()
	pw := masterPassword
	mu.RUnlock()
	if pw != "" {
		return pw, true
	}
	if env := os.Getenv("SSHFRAC_MASTER_PASSWORD"); env != "" {
		return env, true
	}
	if env := os.Getenv("INVOSSH_MASTER_PASSWORD"); env != "" {
		return env, true
	}
	if env := os.Getenv("SSHCTL_MASTER_PASSWORD"); env != "" {
		return env, true
	}
	return "", false
}

func resolveBindMachine() bool {
	mu.RLock()
	set, v := bindMachineSet, bindMachineFlag
	mu.RUnlock()
	if set {
		return v
	}
	for _, key := range []string{"SSHFRAC_BIND_MACHINE", "INVOSSH_BIND_MACHINE", "SSHCTL_BIND_MACHINE"} {
		switch strings.ToLower(os.Getenv(key)) {
		case "1", "true", "yes", "on":
			return true
		}
	}
	return false
}

func encryptV1(plaintext string) (string, error) {
	key, err := machineKey()
	if err != nil {
		return "", err
	}
	out, err := seal(key, []byte(plaintext))
	if err != nil {
		return "", err
	}
	return prefixV1 + base64.StdEncoding.EncodeToString(out), nil
}

func decryptV1(b64 string) (string, error) {
	raw, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return "", fmt.Errorf("decode ciphertext: %w", err)
	}
	key, err := machineKey()
	if err != nil {
		return "", err
	}
	plain, err := open(key, raw)
	if err != nil {
		return "", fmt.Errorf("decrypt failed (wrong machine or corrupted data): %w", err)
	}
	return string(plain), nil
}

func encryptV2(plaintext, password string, bind bool) (string, error) {
	salt := make([]byte, saltSize)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return "", err
	}
	key := deriveV2(password, salt, bind)
	sealed, err := seal(key, []byte(plaintext))
	if err != nil {
		return "", err
	}
	out := append(salt, sealed...)
	return prefixV2 + base64.StdEncoding.EncodeToString(out), nil
}

func decryptV2(b64, password string, bind bool) (string, error) {
	raw, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return "", fmt.Errorf("decode ciphertext: %w", err)
	}
	if len(raw) < saltSize+12 {
		return "", errors.New("ciphertext too short")
	}
	salt, rest := raw[:saltSize], raw[saltSize:]
	key := deriveV2(password, salt, bind)
	plain, err := open(key, rest)
	if err != nil {
		return "", fmt.Errorf("decrypt failed (wrong master password, bind setting, or corrupted data): %w", err)
	}
	return string(plain), nil
}

func deriveV2(password string, salt []byte, bind bool) []byte {
	pw := password
	if bind {
		if mk, err := machineMaterial(); err == nil {
			pw = password + "|" + mk
		}
	}
	// time=1, memory=64MiB, threads=4, keyLen=32
	return argon2.IDKey([]byte(pw), salt, 1, 64*1024, 4, 32)
}

func seal(key, plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	return gcm.Seal(nonce, nonce, plaintext, nil), nil
}

func open(key, raw []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	if len(raw) < gcm.NonceSize() {
		return nil, errors.New("ciphertext too short")
	}
	nonce, ciphertext := raw[:gcm.NonceSize()], raw[gcm.NonceSize():]
	return gcm.Open(nil, nonce, ciphertext, nil)
}

func machineKey() ([]byte, error) {
	material, err := machineMaterial()
	if err != nil {
		return nil, err
	}
	sum := sha256.Sum256([]byte(material))
	return sum[:], nil
}

func machineMaterial() (string, error) {
	host, _ := os.Hostname()
	u, err := user.Current()
	username := "unknown"
	if err == nil {
		username = u.Username
	}
	return strings.Join([]string{
		"sshctl-v1", // stable KDF material; do not rename (breaks enc:v1)
		runtime.GOOS,
		runtime.GOARCH,
		host,
		username,
		machineID(),
	}, "|"), nil
}
