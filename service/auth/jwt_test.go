package auth

import "testing"

func TestCreateJWT(t *testing.T) {
	token, err := CreateJWT([]byte("secret-jwt"), 1, 300)
	if err != nil {
		t.Errorf("error creating jwt token: %v", err)
	}

	if token == "" {
		t.Errorf("token is empty; not created")
	}
}
