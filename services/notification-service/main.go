package main

import (
	"context"
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

// HealthCheck handles health check requests
func (ns *NotificationService) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]string{"status": "healthy"}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// CreateNotification handles notification creation requests
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

// GetUserNotifications handles requests to get user notifications
func (ns *NotificationService) GetUserNotifications(w http.ResponseWriter, r *http.Request) {
	userID := mux.Vars(r)["user_id"]
	if userID == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	ns.mu.RLock()
	var userNotifications []Notification
	for _, n := range ns.notifications {
		if n.UserID == userID {
			userNotifications = append(userNotifications, n)
		}
	}
	ns.mu.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userNotifications)
}

// MarkNotificationAsRead handles requests to mark notifications as read
func (ns *NotificationService) MarkNotificationAsRead(w http.ResponseWriter, r *http.Request) {
	notificationID := mux.Vars(r)["id"]
	if notificationID == "" {
		http.Error(w, "Notification ID is required", http.StatusBadRequest)
		return
	}

	ns.mu.Lock()
	if notification, exists := ns.notifications[notificationID]; exists {
		notification.Read = true
		ns.notifications[notificationID] = notification
		ns.mu.Unlock()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(notification)
	} else {
		ns.mu.Unlock()
		http.Error(w, "Notification not found", http.StatusNotFound)
	}
}

func main() {
	app := NewNotificationService()
	router := mux.NewRouter()

	// Routes
	router.HandleFunc("/health", app.HealthCheck).Methods("GET")
	router.HandleFunc("/notifications", app.CreateNotification).Methods("POST")
	router.HandleFunc("/notifications/user/{user_id}", app.GetUserNotifications).Methods("GET")
	router.HandleFunc("/notifications/{id}/read", app.MarkNotificationAsRead).Methods("PUT")

	// Graceful shutdown setup
	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	// Channel to listen for errors coming from the listener.
	serverErrors := make(chan error, 1)
	// Channel to listen for an interrupt or terminate signal from the OS.
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// Start the service listening for requests.
	go func() {
		log.Printf("API listening on %s", srv.Addr)
		serverErrors <- srv.ListenAndServe()
	}()

	// Blocking main and waiting for shutdown.
	select {
	case err := <-serverErrors:
		log.Fatalf("Error starting server: %v", err)
	case <-shutdown:
		log.Println("Starting shutdown...")
		// Give outstanding requests a deadline for completion.
		const timeout = 5 * time.Second
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		// Asking listener to shut down and shed load.
		if err := srv.Shutdown(ctx); err != nil {
			log.Printf("Graceful shutdown did not complete in %v: %v", timeout, err)
			if err := srv.Close(); err != nil {
				log.Printf("Error killing server: %v", err)
			}
		}
	}
}
