package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gorilla/mux"
)

type Notification struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Message   string    `json:"message"`
	Type      string    `json:"type"`
	CreatedAt time.Time `json:"created_at"`
	Read      bool      `json:"read"`
}

type NotificationService struct {
	notifications map[string]Notification
	mu           sync.RWMutex
}

func NewNotificationService() *NotificationService {
	return &NotificationService{
		notifications: make(map[string]Notification),
	}
}

func (ns *NotificationService) healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
}

func (ns *NotificationService) CreateNotification(w http.ResponseWriter, r *http.Request) {
	var notification Notification
	if err := json.NewDecoder(r.Body).Decode(&notification); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ns.mu.Lock()
	notification.ID = fmt.Sprintf("notif_%d", len(ns.notifications)+1)
	notification.CreatedAt = time.Now()
	notification.Read = false
	ns.notifications[notification.ID] = notification
	ns.mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(notification)
}

func (ns *NotificationService) GetUserNotifications(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["userId"]

	ns.mu.RLock()
	defer ns.mu.RUnlock()

	userNotifications := []Notification{}
	for _, notif := range ns.notifications {
		if notif.UserID == userID {
			userNotifications = append(userNotifications, notif)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userNotifications)
}

func (ns *NotificationService) MarkNotificationAsRead(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	notificationID := vars["notificationId"]

	ns.mu.Lock()
	defer ns.mu.Unlock()

	notification, exists := ns.notifications[notificationID]
	if !exists {
		http.Error(w, "Notification not found", http.StatusNotFound)
		return
	}

	notification.Read = true
	ns.notifications[notificationID] = notification

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(notification)
}

func main() {
	// Initialize the service
	notificationService := NewNotificationService()
	
	// Create router
	r := mux.NewRouter()

	// Register routes
	r.HandleFunc("/health", notificationService.healthCheck).Methods("GET")
	r.HandleFunc("/notifications", notificationService.CreateNotification).Methods("POST")
	r.HandleFunc("/notifications/user/{userId}", notificationService.GetUserNotifications).Methods("GET")
	r.HandleFunc("/notifications/{notificationId}/read", notificationService.MarkNotificationAsRead).Methods("PUT")

	// Create server
	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	// Channel for graceful shutdown
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error starting server: %v\n", err)
		}
	}()

	log.Printf("Server started on port 8080")

	<-done
	log.Print("Server stopped")
}
