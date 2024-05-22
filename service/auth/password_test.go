package auth

import "testing"

func TestHashPassword(t *testing.T) {
	hash, err := HashPassword("password")
	if err != nil {
		t.Errorf("error hashing password: %v", err)
	}

	if hash == "" {
		t.Errorf("hash empty")
	}

	if hash == "password" {
		t.Errorf("input is not hashed")
	}
}

func TestComparePassword(t *testing.T) {
	pass := "password"
	hash, err := HashPassword(pass)
	if err != nil {
		t.Errorf("error hashing password: %v", err)
	}

	if !ComparePassword(hash, []byte(pass)) {
		t.Errorf("expected password to match hash")
	}

	diffPass := "different password"
	if ComparePassword(hash, []byte(diffPass)) {
		t.Errorf("expected password to not match hash")
	}
}
