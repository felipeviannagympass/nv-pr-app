package config

import (
	"encoding/json"
	"os"

	"github.com/kelseyhightower/envconfig"
)

// ServiceConfig ...
type ServiceConfig struct {
	Debug    bool      `json:"debug"`
	Projects []Project `json:"projects"`
	Github   Github    `json:"github"`
	Database Database  `json:"database"`
}
type Squads struct {
	Name      string   `json:"name"`
	Engineers []string `json:"engineers"`
}
type Project struct {
	Owner  string   `json:"owner"`
	Repos  []string `json:"repos"`
	Squads []Squads `json:"squads"`
}

type Github struct {
	AccessToken string `json:"access_token"`
}

type Database struct {
	DbPath         string `json:"db_path"`
	MigrationsPath string `json:"migrations_path"`
}

// LoadServiceConfig ...
func LoadServiceConfig(configFile string) (*ServiceConfig, error) {
	var cfg ServiceConfig

	if err := loadServiceConfigFromFile(configFile, &cfg); err != nil {
		return nil, err
	}

	if err := envconfig.Process("", &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func loadServiceConfigFromFile(configFile string, cfg *ServiceConfig) error {
	_, err := os.Stat(configFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	jsonFile, err := os.ReadFile(configFile)
	if err != nil {
		return err
	}

	return json.Unmarshal(jsonFile, &cfg)
}
