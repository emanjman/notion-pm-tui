package config

import (
	"bufio"
	"os"
	"strings"
)

// Config holds application configuration
type Config struct {
	PriorityProjectIDs map[string]bool
}

// LoadConfig reads the .notion-tui-config file and returns the configuration
// Returns an empty config if file doesn't exist or on error
// Accepts UUIDs with or without hyphens and normalizes them
func LoadConfig(configPath string) *Config {
	cfg := &Config{
		PriorityProjectIDs: make(map[string]bool),
	}

	file, err := os.Open(configPath)
	if err != nil {
		// Config file is optional - return empty config
		return cfg
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Normalize UUID format (add hyphens if missing)
		normalizedID := normalizeUUID(line)
		if normalizedID != "" {
			cfg.PriorityProjectIDs[normalizedID] = true
		}
	}

	return cfg
}

// GetPriorityProjectIDs returns all configured priority project IDs (for debugging)
func (c *Config) GetPriorityProjectIDs() []string {
	ids := make([]string, 0, len(c.PriorityProjectIDs))
	for id := range c.PriorityProjectIDs {
		ids = append(ids, id)
	}
	return ids
}

// normalizeUUID converts a UUID with or without hyphens to the standard format
// Input: "1dbddd31ce45489089940b99f4b6bd45" or "1dbddd31-ce45-4890-8b94-0b99f4b6bd45"
// Output: "1dbddd31-ce45-4890-8b94-0b99f4b6bd45"
func normalizeUUID(id string) string {
	// Remove all hyphens
	clean := strings.ReplaceAll(id, "-", "")

	// Verify it's 32 hex characters
	if len(clean) != 32 {
		return id // Return as-is if not valid length
	}

	// Insert hyphens at correct positions: 8-4-4-4-12
	return clean[0:8] + "-" + clean[8:12] + "-" + clean[12:16] + "-" + clean[16:20] + "-" + clean[20:32]
}

// IsPriorityProject checks if a project ID is in the priority list
// Also normalizes the input ID before checking
func (c *Config) IsPriorityProject(projectID string) bool {
	normalizedID := normalizeUUID(projectID)
	return c.PriorityProjectIDs[normalizedID]
}
