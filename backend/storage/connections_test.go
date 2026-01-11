package storage

import (
	"testing"

	"github.com/gofrs/uuid"
)

const testEmail2 = "test2@example.com"

func ensureTestUser2(t *testing.T) {
	_, err := GetUserByEmail(testEmail2)
	if err != nil {
		_, err = CreateUser(testEmail2, "Test User 2")
		if err != nil {
			t.Fatalf("Failed to create test user 2: %v", err)
		}
	}
}

func TestCreateConnection(t *testing.T) {
	ensureTestUser(t)
	ensureTestUser2(t)

	user1, _ := GetUserByEmail(testEmail)
	user2, _ := GetUserByEmail(testEmail2)

	// Clean up any existing connection
	DeleteConnection(user1.UUID, user2.UUID)

	err := CreateConnection(user1.UUID, user2.UUID)
	if err != nil {
		t.Fatalf("CreateConnection failed: %v", err)
	}

	// Cleanup
	defer DeleteConnection(user1.UUID, user2.UUID)

	// Verify connection exists
	exists, err := ConnectionExists(user1.UUID, user2.UUID)
	if err != nil {
		t.Fatalf("ConnectionExists failed: %v", err)
	}
	if !exists {
		t.Error("Expected connection to exist")
	}
}

func TestCreateConnection_Duplicate(t *testing.T) {
	ensureTestUser(t)
	ensureTestUser2(t)

	user1, _ := GetUserByEmail(testEmail)
	user2, _ := GetUserByEmail(testEmail2)

	// Clean up any existing connection
	DeleteConnection(user1.UUID, user2.UUID)

	err := CreateConnection(user1.UUID, user2.UUID)
	if err != nil {
		t.Fatalf("First CreateConnection failed: %v", err)
	}
	defer DeleteConnection(user1.UUID, user2.UUID)

	// Try to create duplicate
	err = CreateConnection(user1.UUID, user2.UUID)
	if err == nil {
		t.Error("Expected error creating duplicate connection")
	}
}

func TestConnectionExists_NotFound(t *testing.T) {

	// Use random UUIDs that don't exist
	fakeUUID1 := uuid.Must(uuid.NewV4())
	fakeUUID2 := uuid.Must(uuid.NewV4())

	exists, err := ConnectionExists(fakeUUID1, fakeUUID2)
	if err != nil {
		t.Fatalf("ConnectionExists failed: %v", err)
	}
	if exists {
		t.Error("Expected connection to not exist")
	}
}

func TestGetConnectionsBySourceUser(t *testing.T) {
	ensureTestUser(t)
	ensureTestUser2(t)

	user1, _ := GetUserByEmail(testEmail)
	user2, _ := GetUserByEmail(testEmail2)

	// Clean up and create connection
	DeleteConnection(user1.UUID, user2.UUID)
	err := CreateConnection(user1.UUID, user2.UUID)
	if err != nil {
		t.Fatalf("CreateConnection failed: %v", err)
	}
	defer DeleteConnection(user1.UUID, user2.UUID)

	// Get connections for user1 (outgoing)
	connections, err := GetConnectionsBySourceUser(user1.UUID)
	if err != nil {
		t.Fatalf("GetConnectionsBySourceUser failed: %v", err)
	}

	found := false
	for _, u := range connections {
		if u.UUID == user2.UUID {
			found = true
			if u.Email != testEmail2 {
				t.Errorf("Expected email %q, got %q", testEmail2, u.Email)
			}
			break
		}
	}
	if !found {
		t.Error("Expected to find user2 in connections")
	}
}

func TestGetConnectionsByTargetUser(t *testing.T) {
	ensureTestUser(t)
	ensureTestUser2(t)

	user1, _ := GetUserByEmail(testEmail)
	user2, _ := GetUserByEmail(testEmail2)

	// Clean up and create connection
	DeleteConnection(user1.UUID, user2.UUID)
	err := CreateConnection(user1.UUID, user2.UUID)
	if err != nil {
		t.Fatalf("CreateConnection failed: %v", err)
	}
	defer DeleteConnection(user1.UUID, user2.UUID)

	// Get connections for user2 (incoming)
	connections, err := GetConnectionsByTargetUser(user2.UUID)
	if err != nil {
		t.Fatalf("GetConnectionsByTargetUser failed: %v", err)
	}

	found := false
	for _, u := range connections {
		if u.UUID == user1.UUID {
			found = true
			if u.Email != testEmail {
				t.Errorf("Expected email %q, got %q", testEmail, u.Email)
			}
			break
		}
	}
	if !found {
		t.Error("Expected to find user1 in incoming connections")
	}
}

func TestDeleteConnection(t *testing.T) {
	ensureTestUser(t)
	ensureTestUser2(t)

	user1, _ := GetUserByEmail(testEmail)
	user2, _ := GetUserByEmail(testEmail2)

	// Clean up and create connection
	DeleteConnection(user1.UUID, user2.UUID)
	err := CreateConnection(user1.UUID, user2.UUID)
	if err != nil {
		t.Fatalf("CreateConnection failed: %v", err)
	}

	// Delete connection
	err = DeleteConnection(user1.UUID, user2.UUID)
	if err != nil {
		t.Fatalf("DeleteConnection failed: %v", err)
	}

	// Verify connection no longer exists
	exists, err := ConnectionExists(user1.UUID, user2.UUID)
	if err != nil {
		t.Fatalf("ConnectionExists failed: %v", err)
	}
	if exists {
		t.Error("Expected connection to not exist after deletion")
	}
}

func TestGetConnectionsBySourceUser_Empty(t *testing.T) {

	// Use a random UUID that won't have any connections
	fakeUUID := uuid.Must(uuid.NewV4())

	connections, err := GetConnectionsBySourceUser(fakeUUID)
	if err != nil {
		t.Fatalf("GetConnectionsBySourceUser failed: %v", err)
	}

	// Should return empty slice or nil, not error
	if len(connections) != 0 {
		t.Errorf("Expected 0 connections for non-existent user, got %d", len(connections))
	}
}
