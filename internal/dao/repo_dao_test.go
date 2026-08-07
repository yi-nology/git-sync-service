package dao

import (
	"os"
	"testing"

	"github.com/yi-nology/git-platform-sdk/pkg/credential"
)

func TestCryptoManager_EncryptDecrypt_Roundtrip(t *testing.T) {
	// Set up a test encryption key
	key := "0123456789abcdef0123456789abcdef" // 32 bytes for AES-256
	os.Setenv("ENCRYPTION_KEY", key)
	defer os.Unsetenv("ENCRYPTION_KEY")

	cm, err := credential.NewCryptoManager()
	if err != nil {
		t.Fatalf("failed to create CryptoManager: %v", err)
	}

	testCases := []struct {
		name      string
		plaintext string
	}{
		{"normal text", "my-secret-token-123"},
		{"with special chars", "token!@#$%^&*()_+-=[]{}|;':\",./<>?`~"},
		{"unicode", "token-with-unicode-\u4e2d\u6587"},
		{"long string", "a]very-long-token-that-is-quite-substantial-in-length-and-contains-many-characters-for-testing-purposes"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			encrypted, err := cm.Encrypt(tc.plaintext)
			if err != nil {
				t.Fatalf("Encrypt failed: %v", err)
			}

			if encrypted == tc.plaintext {
				t.Errorf("encrypted text should differ from plaintext")
			}

			decrypted, err := cm.Decrypt(encrypted)
			if err != nil {
				t.Fatalf("Decrypt failed: %v", err)
			}

			if decrypted != tc.plaintext {
				t.Errorf("decrypt roundtrip failed: got %q, want %q", decrypted, tc.plaintext)
			}
		})
	}
}

func TestCryptoManager_EmptyString(t *testing.T) {
	key := "0123456789abcdef0123456789abcdef"
	os.Setenv("ENCRYPTION_KEY", key)
	defer os.Unsetenv("ENCRYPTION_KEY")

	cm, err := credential.NewCryptoManager()
	if err != nil {
		t.Fatalf("failed to create CryptoManager: %v", err)
	}

	// Test encrypting empty string
	encrypted, err := cm.Encrypt("")
	if err != nil {
		t.Fatalf("Encrypt empty string failed: %v", err)
	}

	// Empty string should still produce some output (or handle gracefully)
	decrypted, err := cm.Decrypt(encrypted)
	if err != nil {
		t.Fatalf("Decrypt empty string failed: %v", err)
	}

	if decrypted != "" {
		t.Errorf("expected empty string, got %q", decrypted)
	}
}

func TestCryptoManager_WithoutKey(t *testing.T) {
	// Ensure ENCRYPTION_KEY is not set
	os.Unsetenv("ENCRYPTION_KEY")

	_, err := credential.NewCryptoManager()
	if err == nil {
		t.Error("expected error when ENCRYPTION_KEY is not set, got nil")
	}
}

func TestCryptoManager_NewCryptoManagerFromKey(t *testing.T) {
	key := "0123456789abcdef0123456789abcdef"
	cm := credential.NewCryptoManagerFromKey(key)

	plaintext := "test-token-value"
	encrypted, err := cm.Encrypt(plaintext)
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}

	decrypted, err := cm.Decrypt(encrypted)
	if err != nil {
		t.Fatalf("Decrypt failed: %v", err)
	}

	if decrypted != plaintext {
		t.Errorf("decrypt roundtrip failed: got %q, want %q", decrypted, plaintext)
	}
}

func TestDefaultPagination(t *testing.T) {
	tests := []struct {
		name           string
		offset, limit  int
		wantOff, wantLim int
	}{
		{"normal values", 10, 50, 10, 50},
		{"zero limit defaults to 50", 0, 0, 0, 50},
		{"negative limit defaults to 50", 5, -1, 5, 50},
		{"limit over 200 capped", 0, 300, 0, 200},
		{"negative offset becomes 0", -5, 50, 0, 50},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			p := DefaultPagination(tc.offset, tc.limit)
			if p.Offset != tc.wantOff {
				t.Errorf("Offset = %d, want %d", p.Offset, tc.wantOff)
			}
			if p.Limit != tc.wantLim {
				t.Errorf("Limit = %d, want %d", p.Limit, tc.wantLim)
			}
		})
	}
}
