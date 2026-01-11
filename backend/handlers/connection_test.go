package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cobyabrahams/hungr/models"
	"github.com/cobyabrahams/hungr/storage"
)

const testEmail2 = "test2@example.com"

func ensureTestUser2(t *testing.T) {
	_, err := storage.GetUserByEmail(testEmail2)
	if err != nil {
		_, err = storage.CreateUser(testEmail2, "Test User 2")
		if err != nil {
			t.Fatalf("Failed to create test user 2: %v", err)
		}
	}
}

func TestCreateConnection_Success(t *testing.T) {
	ensureTestUser(t)
	ensureTestUser2(t)

	user1, _ := storage.GetUserByEmail(testEmail)
	user2, _ := storage.GetUserByEmail(testEmail2)

	// Clean up any existing connection
	storage.DeleteConnection(user1.UUID, user2.UUID)

	body := `{"target_user_uuid": "` + user2.UUID.String() + `"}`
	req := httptest.NewRequest("POST", "/api/connections?email="+testEmail, bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	CreateConnection(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	// Cleanup
	storage.DeleteConnection(user1.UUID, user2.UUID)
}

func TestCreateConnection_MissingEmail(t *testing.T) {
	body := `{"target_user_uuid": "00000000-0000-0000-0000-000000000001"}`
	req := httptest.NewRequest("POST", "/api/connections", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	CreateConnection(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", resp.StatusCode)
	}
}

func TestCreateConnection_InvalidEmail(t *testing.T) {
	body := `{"target_user_uuid": "00000000-0000-0000-0000-000000000001"}`
	req := httptest.NewRequest("POST", "/api/connections?email=nonexistent@example.com", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	CreateConnection(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", resp.StatusCode)
	}
}

func TestCreateConnection_MissingTargetUser(t *testing.T) {
	ensureTestUser(t)

	body := `{}`
	req := httptest.NewRequest("POST", "/api/connections?email="+testEmail, bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	CreateConnection(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", resp.StatusCode)
	}
}

func TestCreateConnection_SelfConnection(t *testing.T) {
	ensureTestUser(t)

	user1, _ := storage.GetUserByEmail(testEmail)

	body := `{"target_user_uuid": "` + user1.UUID.String() + `"}`
	req := httptest.NewRequest("POST", "/api/connections?email="+testEmail, bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	CreateConnection(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", resp.StatusCode)
	}
}

func TestCreateConnection_TargetNotFound(t *testing.T) {
	ensureTestUser(t)

	body := `{"target_user_uuid": "00000000-0000-0000-0000-000000000001"}`
	req := httptest.NewRequest("POST", "/api/connections?email="+testEmail, bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	CreateConnection(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", resp.StatusCode)
	}
}

func TestCreateConnection_Duplicate(t *testing.T) {
	ensureTestUser(t)
	ensureTestUser2(t)

	user1, _ := storage.GetUserByEmail(testEmail)
	user2, _ := storage.GetUserByEmail(testEmail2)

	// Clean up and create initial connection
	storage.DeleteConnection(user1.UUID, user2.UUID)
	storage.CreateConnection(user1.UUID, user2.UUID)
	defer storage.DeleteConnection(user1.UUID, user2.UUID)

	body := `{"target_user_uuid": "` + user2.UUID.String() + `"}`
	req := httptest.NewRequest("POST", "/api/connections?email="+testEmail, bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	CreateConnection(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusConflict {
		t.Errorf("Expected status 409, got %d", resp.StatusCode)
	}
}

func TestGetConnections_Success(t *testing.T) {
	ensureTestUser(t)
	ensureTestUser2(t)

	user1, _ := storage.GetUserByEmail(testEmail)
	user2, _ := storage.GetUserByEmail(testEmail2)

	// Clean up and create connection
	storage.DeleteConnection(user1.UUID, user2.UUID)
	storage.CreateConnection(user1.UUID, user2.UUID)
	defer storage.DeleteConnection(user1.UUID, user2.UUID)

	req := httptest.NewRequest("GET", "/api/connections?user_uuid="+user1.UUID.String(), nil)
	w := httptest.NewRecorder()

	GetConnections(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	var response models.ConnectionsResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if !response.Success {
		t.Error("Expected success to be true")
	}

	found := false
	for _, u := range response.Connections {
		if u.UUID == user2.UUID {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected to find user2 in connections")
	}
}

func TestGetConnections_Incoming(t *testing.T) {
	ensureTestUser(t)
	ensureTestUser2(t)

	user1, _ := storage.GetUserByEmail(testEmail)
	user2, _ := storage.GetUserByEmail(testEmail2)

	// Clean up and create connection (user1 -> user2)
	storage.DeleteConnection(user1.UUID, user2.UUID)
	storage.CreateConnection(user1.UUID, user2.UUID)
	defer storage.DeleteConnection(user1.UUID, user2.UUID)

	// Get incoming connections for user2
	req := httptest.NewRequest("GET", "/api/connections?user_uuid="+user2.UUID.String()+"&direction=incoming", nil)
	w := httptest.NewRecorder()

	GetConnections(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	var response models.ConnectionsResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	found := false
	for _, u := range response.Connections {
		if u.UUID == user1.UUID {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected to find user1 in incoming connections")
	}
}

func TestGetConnections_MissingUserUUID(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/connections", nil)
	w := httptest.NewRecorder()

	GetConnections(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", resp.StatusCode)
	}
}

func TestGetConnections_InvalidUserUUID(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/connections?user_uuid=invalid", nil)
	w := httptest.NewRecorder()

	GetConnections(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", resp.StatusCode)
	}
}

func TestDeleteConnection_Success(t *testing.T) {
	ensureTestUser(t)
	ensureTestUser2(t)

	user1, _ := storage.GetUserByEmail(testEmail)
	user2, _ := storage.GetUserByEmail(testEmail2)

	// Clean up and create connection
	storage.DeleteConnection(user1.UUID, user2.UUID)
	storage.CreateConnection(user1.UUID, user2.UUID)

	req := httptest.NewRequest("DELETE", "/api/connections?email="+testEmail+"&target_user_uuid="+user2.UUID.String(), nil)
	w := httptest.NewRecorder()

	DeleteConnection(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	// Verify connection was deleted
	exists, _ := storage.ConnectionExists(user1.UUID, user2.UUID)
	if exists {
		t.Error("Expected connection to be deleted")
	}
}

func TestDeleteConnection_MissingEmail(t *testing.T) {
	req := httptest.NewRequest("DELETE", "/api/connections?target_user_uuid=00000000-0000-0000-0000-000000000001", nil)
	w := httptest.NewRecorder()

	DeleteConnection(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", resp.StatusCode)
	}
}

func TestDeleteConnection_InvalidEmail(t *testing.T) {
	req := httptest.NewRequest("DELETE", "/api/connections?email=nonexistent@example.com&target_user_uuid=00000000-0000-0000-0000-000000000001", nil)
	w := httptest.NewRecorder()

	DeleteConnection(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", resp.StatusCode)
	}
}

func TestDeleteConnection_MissingTargetUUID(t *testing.T) {
	ensureTestUser(t)

	req := httptest.NewRequest("DELETE", "/api/connections?email="+testEmail, nil)
	w := httptest.NewRecorder()

	DeleteConnection(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", resp.StatusCode)
	}
}

func TestDeleteConnection_InvalidTargetUUID(t *testing.T) {
	ensureTestUser(t)

	req := httptest.NewRequest("DELETE", "/api/connections?email="+testEmail+"&target_user_uuid=invalid", nil)
	w := httptest.NewRecorder()

	DeleteConnection(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", resp.StatusCode)
	}
}
