package user

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/nathaniel-alvin/tireappBE/types"
)

func TestUserService(t *testing.T) {
	// mock userStore
	userStore := &mockUserStore{}
	handler := NewHandler(userStore)

	t.Run("should fail when user request is invalid", func(t *testing.T) {
		request := types.RegisterUserRequest{
			Username: "testing",
			Email:    "",
			Password: "test123",
		}
		marshalled, _ := json.Marshal(request)

		req, err := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(marshalled))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()

		router.HandleFunc("/register", handler.handleRegister)
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected status code %d, got %d", http.StatusBadRequest, rr.Code)
		}
	})

	t.Run("should pass when user request is valid", func(t *testing.T) {
		request := types.RegisterUserRequest{
			Username: "testing",
			Email:    "validmail@mail.com",
			Password: "test123",
		}
		marshalled, _ := json.Marshal(request)

		req, err := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(marshalled))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()

		router.HandleFunc("/register", handler.handleRegister)
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusCreated {
			t.Errorf("expected status code %d, got %d", http.StatusCreated, rr.Code)
		}
	})
}

type mockUserStore struct{}

func (m *mockUserStore) GetUserByEmail(email string) (*types.User, error) {
	return &types.User{}, fmt.Errorf("user found")
}

func (m *mockUserStore) CreateUser(types.User) (int, error) {
	return 0, nil
}

func (m *mockUserStore) GetUserById(id int) (*types.User, error) {
	return nil, nil
}
