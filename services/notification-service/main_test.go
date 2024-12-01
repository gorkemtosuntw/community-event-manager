package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHealthCheck(t *testing.T) {
	app := NewNotificationService()
	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(app.HealthCheck)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Trim any whitespace and newlines
	got := strings.TrimSpace(rr.Body.String())
	expected := `{"status":"healthy"}`
	if got != expected {
		t.Errorf("handler returned unexpected body: got %q want %q",
			got, expected)
	}
}

func TestCreateNotification(t *testing.T) {
	app := NewNotificationService()
	notification := Notification{
		UserID:  "test-user",
		Message: "Test notification",
		Type:    "test",
	}

	payload, err := json.Marshal(notification)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", "/notifications", bytes.NewBuffer(payload))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(app.CreateNotification)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}

	var response Notification
	err = json.NewDecoder(rr.Body).Decode(&response)
	if err != nil {
		t.Fatal(err)
	}

	if response.ID == "" {
		t.Error("response should contain an id")
	}
	if response.UserID != notification.UserID {
		t.Errorf("expected UserID %v, got %v", notification.UserID, response.UserID)
	}
	if response.Message != notification.Message {
		t.Errorf("expected Message %v, got %v", notification.Message, response.Message)
	}
}
